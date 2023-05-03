package cwflags

import (
	"github.com/ncw/ncwtester/cwgenerator"
	"github.com/spf13/pflag"
)

const (
	channelNum      = 2
	bitDepthInBytes = 2
	maxSampleValue  = 32767
)

var (
	sampleRate int
	wpm        float64
	frequency  float64
)

// Add the CW flags to the flagset passed in
func Add(flags *pflag.FlagSet) {
	flags.IntVarP(&sampleRate, "samplerate", "s", 44100, "sample rate")
	flags.Float64VarP(&wpm, "wpm", "", 25.0, "WPM to send at")
	flags.Float64VarP(&frequency, "frequency", "", 600.0, "HZ of morse")
}

// NewOpt creates a new set of cwgenerator.Options from the command
// line flags
func NewOpt() cwgenerator.Options {
	return cwgenerator.Options{
		WPM:             wpm,
		Frequency:       frequency,
		SampleRate:      sampleRate,
		ChannelNum:      channelNum,
		BitDepthInBytes: bitDepthInBytes,
		MaxSampleValue:  maxSampleValue,
	}
}
