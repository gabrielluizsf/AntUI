//go:build windows

package antui

import (
	"fmt"
	"math"
	"runtime"
	"syscall"
	"time"
	"unsafe"
)

// WASAPI, driven through syscall like the rest of the Windows backend: no
// cgo, no C toolchain, nothing to link.
//
// Shared mode and event driven, which is what a game wants — the device wakes
// us when it needs the next block instead of being polled, and the block is
// about ten milliseconds. The older waveOut is three times that and would be
// heard on anything that has to line up with a hit.
//
// COM without a compiler is a vtable and an index: every interface pointer is
// a pointer to a table of functions, `this` is the first argument, and the
// index is the method's place in the table. The tables are in the order the
// interfaces are declared in mmdeviceapi.h and audioclient.h, which is what
// the index comments below are.
//
// NOT YET RUN ON WINDOWS. It is written from the interface definitions and
// checked by a test that only runs there.

var (
	ole32    = syscall.NewLazyDLL("ole32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	procCoInitializeEx      = ole32.NewProc("CoInitializeEx")
	procCoUninitialize      = ole32.NewProc("CoUninitialize")
	procCoCreateInstance    = ole32.NewProc("CoCreateInstance")
	procCoTaskMemFree       = ole32.NewProc("CoTaskMemFree")
	procCreateEventW        = kernel32.NewProc("CreateEventW")
	procSetEvent            = kernel32.NewProc("SetEvent")
	procCloseHandle         = kernel32.NewProc("CloseHandle")
	procWaitForSingleObject = kernel32.NewProc("WaitForSingleObject")
)

// A guid is what COM calls everything by.
type guid struct {
	a uint32
	b uint16
	c uint16
	d [8]byte
}

var (
	clsidMMDeviceEnumerator = guid{0xBCDE0395, 0xE52F, 0x467C,
		[8]byte{0x8E, 0x3D, 0xC4, 0x57, 0x92, 0x91, 0x69, 0x2E}}
	iidIMMDeviceEnumerator = guid{0xA95664D2, 0x9614, 0x4F35,
		[8]byte{0xA7, 0x46, 0xDE, 0x8D, 0xB6, 0x36, 0x17, 0xE6}}
	iidIAudioClient = guid{0x1CB9AD4C, 0xDBFA, 0x4C32,
		[8]byte{0xB1, 0x78, 0xC2, 0xF5, 0x68, 0xA7, 0x03, 0xB2}}
	iidIAudioRenderClient = guid{0xF294ACFC, 0x3146, 0x4483,
		[8]byte{0xA7, 0xBF, 0xAD, 0xDC, 0xA7, 0xC2, 0x60, 0xE2}}

	// The two sample layouts a mix format is ever in.
	subtypeFloat = guid{0x00000003, 0x0000, 0x0010,
		[8]byte{0x80, 0x00, 0x00, 0xAA, 0x00, 0x38, 0x9B, 0x71}}
	subtypePCM = guid{0x00000001, 0x0000, 0x0010,
		[8]byte{0x80, 0x00, 0x00, 0xAA, 0x00, 0x38, 0x9B, 0x71}}
)

const (
	coinitMultithreaded = 0
	clsctxAll           = 23

	eRender  = 0 // the direction: out of the machine
	eConsole = 0 // and what for: games and everything else that is not a call

	shareModeShared     = 0
	streamEventCallback = 0x00040000

	waitObject0 = 0
	waitTimeout = 0x00000102

	formatFloat      = 3      // WAVE_FORMAT_IEEE_FLOAT
	formatPCM        = 1      // WAVE_FORMAT_PCM
	formatExtensible = 0xFFFE // and the tag that says look at the subformat

	// The one error worth telling apart: it is what a machine with no sound
	// card at all answers with.
	errNotFound          = 0x80070490 // ERROR_NOT_FOUND, as an HRESULT
	errDeviceInvalidated = 0x88890004
)

// A waveFormatEx is the shape of a stream. The extensible form adds three
// fields after it, and a mix format is always the extensible one.
type waveFormatEx struct {
	tag            uint16
	channels       uint16
	samplesPerSec  uint32
	avgBytesPerSec uint32
	blockAlign     uint16
	bitsPerSample  uint16
	size           uint16
}

type waveFormatExtensible struct {
	waveFormatEx
	samples     uint16 // a union; the valid bits per sample
	channelMask uint32
	subFormat   guid
}

// call invokes the method at an index of an interface's vtable. `this` goes
// in front of the arguments, which is what makes a COM call a COM call.
func comCall(this uintptr, index int, args ...uintptr) uintptr {
	vtable := *(**[64]uintptr)(unsafe.Pointer(this))
	all := make([]uintptr, 0, len(args)+1)
	all = append(all, this)
	all = append(all, args...)
	ret, _, _ := syscall.SyscallN(vtable[index], all...)
	return ret
}

// release is IUnknown::Release, which every interface has third.
func comRelease(this uintptr) {
	if this != 0 {
		comCall(this, 2)
	}
}

// failed turns an HRESULT into an error. Anything with the top bit set is a
// failure; everything else, S_FALSE included, is not.
func failed(hr uintptr, what string) error {
	if int32(hr) >= 0 {
		return nil
	}
	if uint32(hr) == errNotFound {
		return fmt.Errorf("%w: %s found no device", ErrNoAudio, what)
	}
	return fmt.Errorf("antui: %s failed (0x%08X)", what, uint32(hr))
}

// A wasapi is one render stream and the thread feeding it.
type wasapi struct {
	audio *Audio

	ready chan error
	quit  chan struct{}
	done  chan struct{}

	event   uintptr
	delay   time.Duration
	stopped bool
}

func (w *wasapi) latency() time.Duration { return w.delay }

func (w *wasapi) close() error {
	select {
	case <-w.quit:
		return nil // already asked
	default:
	}
	close(w.quit)

	// Wake the thread so it notices, rather than leaving it asleep until the
	// device happens to want another block.
	if w.event != 0 {
		procSetEvent.Call(w.event)
	}
	<-w.done
	return nil
}

// openAudio is the Windows half of OpenAudio.
//
// Everything happens on the stream's own thread, because COM is threaded: an
// interface pointer belongs to the apartment it was made in, and using it
// from another goroutine is undefined in a way that works right up until it
// does not.
func openAudio(name string, a *Audio) (driver, error) {
	w := &wasapi{
		audio: a,
		ready: make(chan error, 1),
		quit:  make(chan struct{}),
		done:  make(chan struct{}),
	}
	go w.run()
	if err := <-w.ready; err != nil {
		return nil, err
	}
	return w, nil
}

func (w *wasapi) run() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	defer close(w.done)

	hr, _, _ := procCoInitializeEx.Call(0, coinitMultithreaded)
	if err := failed(hr, "CoInitializeEx"); err != nil {
		w.ready <- err
		return
	}
	defer procCoUninitialize.Call()

	client, render, frames, format, err := w.start()
	if err != nil {
		w.ready <- err
		return
	}
	defer comRelease(render)
	defer comRelease(client)
	defer func() {
		if w.event != 0 {
			procCloseHandle.Call(w.event)
		}
	}()
	defer comCall(client, 11) // IAudioClient::Stop

	w.ready <- nil
	w.feed(client, render, frames, format)
}

// start opens the default endpoint and gets as far as a running stream. It
// answers the client, the render interface, how many frames the device's
// buffer holds, and the format it wants them in.
func (w *wasapi) start() (client, render uintptr, frames uint32, format *waveFormatExtensible, err error) {
	var enumerator uintptr
	hr, _, _ := procCoCreateInstance.Call(
		uintptr(unsafe.Pointer(&clsidMMDeviceEnumerator)), 0, clsctxAll,
		uintptr(unsafe.Pointer(&iidIMMDeviceEnumerator)),
		uintptr(unsafe.Pointer(&enumerator)))
	if err := failed(hr, "CoCreateInstance"); err != nil {
		return 0, 0, 0, nil, err
	}
	defer comRelease(enumerator)

	// IMMDeviceEnumerator::GetDefaultAudioEndpoint is the fourth method.
	var device uintptr
	if err := failed(comCall(enumerator, 4, eRender, eConsole,
		uintptr(unsafe.Pointer(&device))), "GetDefaultAudioEndpoint"); err != nil {
		return 0, 0, 0, nil, err
	}
	defer comRelease(device)

	// IMMDevice::Activate is the fourth.
	if err := failed(comCall(device, 3,
		uintptr(unsafe.Pointer(&iidIAudioClient)), clsctxAll, 0,
		uintptr(unsafe.Pointer(&client))), "Activate"); err != nil {
		return 0, 0, 0, nil, err
	}

	// IAudioClient::GetMixFormat, the ninth. What comes back is what the
	// shared mixer is running at, and taking it as it stands is what makes
	// this work on a device at any rate at all — the mixer above follows,
	// rather than a resampler being written here.
	var mix *waveFormatExtensible
	if err := failed(comCall(client, 8, uintptr(unsafe.Pointer(&mix))),
		"GetMixFormat"); err != nil {
		comRelease(client)
		return 0, 0, 0, nil, err
	}
	defer procCoTaskMemFree.Call(uintptr(unsafe.Pointer(mix)))

	// A copy, because the original is COM's memory and is about to go.
	format = new(waveFormatExtensible)
	*format = *mix
	if !w.understands(format) {
		comRelease(client)
		return 0, 0, 0, nil, fmt.Errorf(
			"antui: the device wants %d-bit samples in format %d, which this does not write",
			format.bitsPerSample, format.tag)
	}

	// IAudioClient::Initialize, the fourth. A duration of zero in shared mode
	// asks for the engine's own period, which is the shortest it will give.
	if err := failed(comCall(client, 3, shareModeShared, streamEventCallback,
		0, 0, uintptr(unsafe.Pointer(&format.waveFormatEx)), 0),
		"Initialize"); err != nil {
		comRelease(client)
		return 0, 0, 0, nil, err
	}

	event, _, _ := procCreateEventW.Call(0, 0, 0, 0)
	if event == 0 {
		comRelease(client)
		return 0, 0, 0, nil, fmt.Errorf("antui: no event for the audio device")
	}
	w.event = event

	// SetEventHandle is the fourteenth.
	if err := failed(comCall(client, 13, event), "SetEventHandle"); err != nil {
		comRelease(client)
		return 0, 0, 0, nil, err
	}
	// GetBufferSize, the fifth.
	if err := failed(comCall(client, 4, uintptr(unsafe.Pointer(&frames))),
		"GetBufferSize"); err != nil {
		comRelease(client)
		return 0, 0, 0, nil, err
	}
	// GetService, the fifteenth.
	if err := failed(comCall(client, 14,
		uintptr(unsafe.Pointer(&iidIAudioRenderClient)),
		uintptr(unsafe.Pointer(&render))), "GetService"); err != nil {
		comRelease(client)
		return 0, 0, 0, nil, err
	}
	// Start, the eleventh.
	if err := failed(comCall(client, 10), "Start"); err != nil {
		comRelease(render)
		comRelease(client)
		return 0, 0, 0, nil, err
	}

	w.audio.rate = int(format.samplesPerSec)
	w.delay = time.Duration(frames) * time.Second / time.Duration(format.samplesPerSec)
	return client, render, frames, format, nil
}

// understands reports whether the samples are ones this can write.
func (w *wasapi) understands(f *waveFormatExtensible) bool {
	if f.channels == 0 || f.samplesPerSec == 0 {
		return false
	}
	switch f.tag {
	case formatFloat:
		return f.bitsPerSample == 32
	case formatPCM:
		return f.bitsPerSample == 16
	case formatExtensible:
		if f.subFormat == subtypeFloat {
			return f.bitsPerSample == 32
		}
		return f.subFormat == subtypePCM && f.bitsPerSample == 16
	}
	return false
}

// floats reports whether the device takes samples as they already are.
func (w *wasapi) floats(f *waveFormatExtensible) bool {
	return f.tag == formatFloat || (f.tag == formatExtensible && f.subFormat == subtypeFloat)
}

// feed is the stream's loop: wait to be asked, hand over a block, repeat.
func (w *wasapi) feed(client, render uintptr, frames uint32, format *waveFormatExtensible) {
	block := make([]float32, int(frames)*2)
	channels := int(format.channels)
	asFloats := w.floats(format)

	for {
		state, _, _ := procWaitForSingleObject.Call(w.event, 2000)
		select {
		case <-w.quit:
			return
		default:
		}
		if state == waitTimeout {
			// Two seconds without the device asking for anything is a device
			// that has gone away — pulled out, or switched to another one.
			w.audio.fail(fmt.Errorf("antui: the audio device stopped asking for sound"))
			return
		}
		if state != waitObject0 {
			w.audio.fail(fmt.Errorf("antui: waiting on the audio device gave %d", state))
			return
		}

		// GetCurrentPadding, the seventh: how much of the buffer the device
		// has not played yet. The rest is ours to fill.
		var padding uint32
		if err := failed(comCall(client, 6, uintptr(unsafe.Pointer(&padding))),
			"GetCurrentPadding"); err != nil {
			w.audio.fail(err)
			return
		}
		want := frames - padding
		if want == 0 {
			continue
		}

		// IAudioRenderClient::GetBuffer, the fourth.
		var buffer uintptr
		if err := failed(comCall(render, 3, uintptr(want),
			uintptr(unsafe.Pointer(&buffer))), "GetBuffer"); err != nil {
			w.audio.fail(err)
			return
		}

		w.audio.pull(block[:int(want)*2])
		w.spread(buffer, block[:int(want)*2], int(want), channels, asFloats)

		// ReleaseBuffer, the fifth.
		if err := failed(comCall(render, 4, uintptr(want), 0), "ReleaseBuffer"); err != nil {
			w.audio.fail(err)
			return
		}
	}
}

// spread writes a stereo block into the device's buffer, in the device's own
// layout.
//
// The mixer above is stereo and a device may be anything: a headset is two, a
// living room is six. The two go to the front pair and the rest are left
// silent, which is what every game does with a stereo mix — putting the same
// sound in the surrounds as well makes a room sound like a corridor.
func (w *wasapi) spread(buffer uintptr, block []float32, frames, channels int, asFloats bool) {
	if asFloats {
		out := unsafe.Slice((*float32)(unsafe.Pointer(buffer)), frames*channels)
		for i := range out {
			out[i] = 0
		}
		for frame := range frames {
			out[frame*channels] = block[frame*2]
			if channels > 1 {
				out[frame*channels+1] = block[frame*2+1]
			}
		}
		return
	}

	out := unsafe.Slice((*int16)(unsafe.Pointer(buffer)), frames*channels)
	for i := range out {
		out[i] = 0
	}
	for frame := range frames {
		out[frame*channels] = pcm16(block[frame*2])
		if channels > 1 {
			out[frame*channels+1] = pcm16(block[frame*2+1])
		}
	}
}

// pcm16 is a sample as the older devices want it, clipped rather than wrapped
// — a sample that overflows into the opposite sign is a click, and a click is
// the loudest thing a speaker can make.
func pcm16(sample float32) int16 {
	v := float64(sample) * 32767
	if v > 32767 {
		v = 32767
	}
	if v < -32768 {
		v = -32768
	}
	return int16(math.Round(v))
}