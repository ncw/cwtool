// Package cw describes the implementation of CW generators and players
package cw

// CW is an interface to cover several implementations
type CW interface {
	// Rune adds r to the output
	Rune(r rune)

	// String adds s to the output
	String(s string)

	// Sync by waiting for all the morse to be played
	Sync()

	// Close the file
	Close() error
}

// Options to configure the CW generator and player
type Options struct {
	WPM             float64 // WPM to send morse at
	Frequency       float64 // Frequency to generate morse at
	SampleRate      int     // samples per second to generate
	ChannelNum      int
	BitDepthInBytes int
	MaxSampleValue  int
	Continuous      bool   // generates CW continously, never returns EOF from Read
	OutputFile      string // file to send output to
	Debug           bool   // print info messages to stdout
}
