package audio

import (
	"testing"

	"github.com/gabrielluizsf/antui/template/event"
)

func TestBanksLoadedFully(t *testing.T) {
	for _, bank := range []*Template{Tech, Simple} {
		if bank == nil {
			t.Fatal("a shipped bank is nil")
		}
		for _, sound := range []*Sound{bank.Button, bank.Checkbox, bank.Radio, bank.TextInput} {
			if sound == nil {
				t.Errorf("%s has a missing sound", bank.Name)
				continue
			}
			if len(sound.sample) == 0 {
				t.Errorf("%s/%v has no samples", bank.Name, sound)
			}
			if sound.rate != 44100 {
				t.Errorf("%s/%v at rate %d, want 44100", bank.Name, sound, sound.rate)
			}
		}
	}
}

func TestBanksDiffer(t *testing.T) {
	// The two banks come from different packs; they should not be the same
	// recording or every look would sound identical.
	if len(Tech.Button.sample) == len(Simple.Button.sample) {
		same := true
		for i := range Tech.Button.sample {
			if Tech.Button.sample[i] != Simple.Button.sample[i] {
				same = false
				break
			}
		}
		if same {
			t.Error("Tech.Button and Simple.Button are the same samples")
		}
	}
}

func TestOnPlaysWithoutADevice(t *testing.T) {
	Use(nil)
	// Every event a template can report must be answerable in silence: the
	// bank should never crash no matter what the visual half says happened.
	events := []event.Event{
		{Component: event.Button, Kind: event.Click},
		{Component: event.Checkbox, Kind: event.Toggle},
		{Component: event.Radio, Kind: event.Select},
		{Component: event.Slider, Kind: event.Change},
		{Component: event.TextInput, Kind: event.Focus},
		{Component: event.TextInput, Kind: event.Type, Step: 5},
	}
	for _, e := range events {
		Tech.On(e)
		Simple.On(e)
	}
}

func TestUseRoutesPlayback(t *testing.T) {
	p := NewPlayer(44100)
	Use(p) // Tech now plays through p
	Tech.Button.Play()

	block := make([]float32, 512)
	p.Fill(block)
	var heard bool
	for _, v := range block {
		if v != 0 {
			heard = true
			break
		}
	}
	if !heard {
		t.Error("Tech.Button.Play should be heard once a Player is in use")
	}
	Use(nil)
}
