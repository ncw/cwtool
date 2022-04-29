// Package cwgenerator implements a morse code player to an io.Reader interface
package cwgenerator

import (
	"fmt"
	"io"
	"math"
	"time"
	"unicode"
)

// Options to configure a Generator
type Options struct {
	WPM             float64 // WPM to send morse at
	Frequency       float64 // Frequency to generate morse at
	SampleRate      int     // samples per second to generate
	ChannelNum      int
	BitDepthInBytes int
	MaxSampleValue  int
}

// Generator contains state for the morse generation
type Generator struct {
	opt           Options
	sequence      []byte    // sequence to play samples in
	sampleLength  int       // length of sample in bytes
	samples       [2][]byte // samples to play
	sequenceIndex int       // which sample we are playing
	sampleOffset  int       // how far we've got through that sample
	eofTime       time.Time
}

// New makes a new player with the Options passed in
func New(opt Options) *Generator {
	cw := &Generator{}

	ditTimeSeconds := wpmToDitTime(opt.WPM)
	cyclesPerDit := opt.Frequency * ditTimeSeconds
	// fmt.Printf("cyclesPerDit = %.3f at %.1f Hz\n", cyclesPerDit, opt.Frequency)
	// Round cycles per dit to an exact number to avoid clicks
	// this changes the frequency slightly
	cyclesPerDit = math.Round(cyclesPerDit)
	// newFrequency := cyclesPerDit / ditTimeSeconds
	// fmt.Printf("cyclesPerDit = %.3f at %.1f Hz\n", cyclesPerDit, newFrequency)

	samplesPerDit := int(math.Round(float64(opt.SampleRate) * ditTimeSeconds))
	sampleWidth := opt.ChannelNum * opt.BitDepthInBytes
	cw.sampleLength = samplesPerDit * sampleWidth
	cw.samples[0] = make([]byte, cw.sampleLength)
	cw.samples[1] = make([]byte, cw.sampleLength)
	dit := cw.samples[1]
	for i := 0; i < samplesPerDit; i++ {
		b := int16(math.Sin(2*math.Pi*float64(i)/float64(samplesPerDit)*cyclesPerDit) * 0.3 * float64(opt.MaxSampleValue))
		for ch := 0; ch < opt.ChannelNum; ch++ {
			dit[sampleWidth*i+2*ch] = byte(b)
			dit[sampleWidth*i+1+2*ch] = byte(b >> 8)
		}
	}
	return cw
}

// Clear empties the sequence and resets the state
func (cw *Generator) Clear() {
	cw.sequence = cw.sequence[:0]
	cw.sequenceIndex = 0
	cw.sampleOffset = 0
}

// Time it should take to play the morse
func (cw *Generator) duration() time.Duration {
	return time.Duration((float64(len(cw.sequence)) * wpmToDitTime(cw.opt.WPM)) * float64(time.Second))
}

// Read implements the io.Reader interface for the sound data
func (cw *Generator) Read(buf []byte) (n int, err error) {
	for len(buf) > 0 {
		// fmt.Printf("sequenceIndex = %d, sequenceOffset=%d\n", cw.sequenceIndex, cw.sampleOffset)
		if cw.sequenceIndex >= len(cw.sequence) {
			err = io.EOF
			if cw.eofTime.IsZero() {
				cw.eofTime = time.Now()
			}
			break
		}

		// What sample are we supposed to be playing
		sampleIndex := cw.sequence[cw.sequenceIndex]
		// Get its waveform
		sample := cw.samples[sampleIndex]
		nn := copy(buf, sample[cw.sampleOffset:])
		if nn == 0 {
			break
		}
		n += nn
		cw.sampleOffset += nn
		if cw.sampleOffset >= cw.sampleLength {
			cw.sequenceIndex++
			cw.sampleOffset = 0
		}
		buf = buf[nn:]
	}

	// fmt.Printf("Return %d, %v\n", n, err)
	return n, err
}

// Adds the rune to the output
func (cw *Generator) Rune(r rune) {
	r = unicode.ToUpper(r)
	code := morseCode[r]
	if code == "" {
		fmt.Printf("Don't know how to play %v\n", r)
		return
	}
	out := func(b ...byte) {
		cw.sequence = append(cw.sequence, b...)
	}
	for _, c := range code {
		switch c {
		case '-':
			out(1, 1, 1, 0)
		case '.':
			out(1, 0)
		case ' ':
			// word space is 7 dits
			// we've written 1 on the last dit/dah
			// and we'll write 2 after this
			// so need 4 more
			out(0, 0, 0, 0)
		default:
			panic("Bad symbol in code")
		}
	}
	// write letter gap of 3 dits - have written one already
	out(0, 0)
}

// Adds the string to the output
func (cw *Generator) String(s string) {
	for _, r := range s {
		cw.Rune(r)
	}
}

// check interfaces
var _ io.Reader = (*Generator)(nil)
