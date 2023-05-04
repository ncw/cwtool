// Package cwflags configures a CW player from the flags
package cwflags

import (
	"github.com/ncw/cwtool/cw"
	"github.com/ncw/cwtool/cwfile"
	"github.com/ncw/cwtool/cwplayer"
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
	outputFile string
)

// Add the CW flags to the flagset passed in
func Add(flags *pflag.FlagSet) {
	flags.IntVarP(&sampleRate, "samplerate", "s", 44100, "sample rate")
	flags.Float64VarP(&wpm, "wpm", "", 25.0, "WPM to send at")
	flags.Float64VarP(&frequency, "frequency", "", 600.0, "HZ of morse")
	flags.StringVarP(&outputFile, "out", "", "", "WAV file for output instead of speaker")
}

// NewOpt creates a new set of cw.Options from the command line flags
func NewOpt() *cw.Options {
	return &cw.Options{
		WPM:             wpm,
		Frequency:       frequency,
		SampleRate:      sampleRate,
		ChannelNum:      channelNum,
		BitDepthInBytes: bitDepthInBytes,
		MaxSampleValue:  maxSampleValue,
		OutputFile:      outputFile,
	}
}

// NewPlayer creates a new player from the options
func NewPlayer(opt *cw.Options) (cw.CW, error) {
	if opt.OutputFile == "" {
		return cwplayer.New(opt)
	}
	return cwfile.New(opt)
}
