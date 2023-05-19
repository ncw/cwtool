// Package cwflags configures a CW player from the flags
package cwflags

import (
	"github.com/ncw/cwtool/cmd"
	"github.com/ncw/cwtool/cw"
	"github.com/ncw/cwtool/cwfile"
	"github.com/ncw/cwtool/cwplayer"
	"github.com/spf13/pflag"
)

const (
	bitDepthInBytes = 2
	maxSampleValue  = 32767
)

var (
	sampleRate int
	channels   int
	wpm        float64
	farnsworth float64
	frequency  float64
	outputFile string
)

// Add the CW flags to the flagset passed in
func Add(flags *pflag.FlagSet) {
	flags.IntVarP(&sampleRate, "samplerate", "s", 8000, "sample rate in samples/s")
	flags.IntVarP(&channels, "channels", "c", 1, "channels to generate")
	flags.Float64VarP(&wpm, "wpm", "", 25.0, "WPM to send at")
	flags.Float64VarP(&farnsworth, "farnsworth", "", 0.0, "Increase character spacing to match this WPM")
	flags.Float64VarP(&frequency, "frequency", "", 600.0, "HZ of Morse")
	flags.StringVarP(&outputFile, "out", "", "", "WAV file for output instead of speaker")
}

// NewOpt creates a new set of cw.Options from the command line flags
func NewOpt() *cw.Options {
	return &cw.Options{
		WPM:             wpm,
		Farnsworth:      farnsworth,
		Frequency:       frequency,
		SampleRate:      sampleRate,
		Channels:        channels,
		BitDepthInBytes: bitDepthInBytes,
		MaxSampleValue:  maxSampleValue,
		OutputFile:      outputFile,
		Debug:           cmd.Debug,
	}
}

// NewPlayer creates a new player from the options
func NewPlayer(opt *cw.Options) (cw.CW, error) {
	if opt.OutputFile == "" {
		return cwplayer.New(opt)
	}
	return cwfile.New(opt)
}
