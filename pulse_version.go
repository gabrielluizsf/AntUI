//go:build linux && !android

package antui

const (
	pulseVersion = 13

	// The tags a value carries. The protocol is big-endian throughout, and
	// the samples inside a data packet are not — they are whatever the sample
	// spec said, which here is little-endian floats.
	tagString     = 't'
	tagStringNil  = 'N'
	tagU32        = 'L'
	tagU8         = 'B'
	tagArbitrary  = 'x'
	tagTrue       = '1'
	tagFalse      = '0'
	tagSampleSpec = 'a'
	tagChannelMap = 'm'
	tagCvolume    = 'v'
	tagProplist   = 'P'

	// The commands this needs. A server says a great deal more than this and
	// all of it is ignored.
	cmdError                = 0
	cmdReply                = 2
	cmdCreatePlaybackStream = 3
	cmdDeletePlaybackStream = 4
	cmdAuth                 = 8
	cmdSetClientName        = 9
	cmdRequest              = 61

	sampleFloat32LE = 5 // PA_SAMPLE_FLOAT32LE
	frontLeft       = 1 // PA_CHANNEL_POSITION_FRONT_LEFT
	frontRight      = 2
	volumeNorm      = 0x10000 // PA_VOLUME_NORM, which is unattenuated

	noIndex    = 0xFFFFFFFF // PA_INVALID_INDEX
	noPacket   = 0xFFFFFFFF // a descriptor's channel, for a command
	maxPayload = 64 << 10   // the biggest packet this will read
)
