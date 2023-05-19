// Package cwgenerator implements a Morse code player to an io.Reader interface
package cwgenerator

import (
	"fmt"
	"io"
	"math"
	"sync"
	"time"
	"unicode"

	"github.com/ncw/cwtool/cw"
)

// Generator contains state for the Morse generation
type Generator struct {
	opt          *cw.Options
	sequenceMu   sync.Mutex // hold mutex when adding/removing things from sequence
	sequence     []byte     // sequence to play samples in
	sampleLength int        // length of sample in bytes
	samples      [2][]byte  // samples to play
	sampleIndex  byte       // index of sample we are playing now
	sampleOffset int        // how far we've got through that sample
	extraDits    int        // extra dits after each letter
}

// New makes a new player with the Options passed in
func New(opt *cw.Options) *Generator {
	cw := &Generator{
		opt: opt,
	}

	ditTimeSeconds := wpmToDitTime(opt.WPM)
	cyclesPerDit := opt.Frequency * ditTimeSeconds
	if cw.opt.Debug {
		fmt.Printf("cyclesPerDit = %.3f at %.1f Hz\n", cyclesPerDit, opt.Frequency)
	}
	// Round cycles per dit to an exact number to avoid clicks
	// this changes the frequency slightly
	cyclesPerDit = math.Round(cyclesPerDit)
	newFrequency := cyclesPerDit / ditTimeSeconds
	if cw.opt.Debug {
		fmt.Printf("cyclesPerDit = %.3f at %.1f Hz\n", cyclesPerDit, newFrequency)
	}

	// Compute number of extra dit times to meet Farnsworth target
	if opt.Farnsworth > 0 && opt.Farnsworth < opt.WPM {
		// Time to send one PARIS word is
		wordTimeNormal := 60 / opt.WPM
		// Time to send one PARIS word at the Farnsworth speed is
		wordTimeFarnsworth := 60 / opt.Farnsworth
		// So we need to slow each word down by this much
		wordDelay := wordTimeFarnsworth - wordTimeNormal
		// We put a Farnsworth space in after each letter and
		// one extra at the end of the word for 6 in PARIS so
		// the Farnsworth delay must be
		delay := wordDelay / 6
		// Calculate what this is in dits
		extraDits := delay / ditTimeSeconds
		// Round to nearest dit
		cw.extraDits = int(extraDits + 0.5)
		if cw.extraDits <= 0 {
			cw.extraDits = 1
		}
		if cw.opt.Debug {
			fmt.Printf("Farnsworth at %.1f WPM using %.1f WPM needs %.1f extra dits\n", opt.Farnsworth, opt.WPM, extraDits)
			actualWordTime := wordTimeNormal + 6*ditTimeSeconds*float64(cw.extraDits)
			actualWPM := 60 / actualWordTime
			fmt.Printf("This rounds to %d extra dits which makes an actual Farnsworth of %.1f WPM\n", cw.extraDits, actualWPM)
		}
	}

	samplesPerDit := int(math.Round(float64(opt.SampleRate) * ditTimeSeconds))
	sampleWidth := opt.Channels * opt.BitDepthInBytes
	cw.sampleLength = samplesPerDit * sampleWidth
	cw.samples[0] = make([]byte, cw.sampleLength)
	cw.samples[1] = make([]byte, cw.sampleLength)
	dit := cw.samples[1]
	for i := 0; i < samplesPerDit; i++ {
		b := int16(math.Sin(2*math.Pi*float64(i)/float64(samplesPerDit)*cyclesPerDit) * 0.3 * float64(opt.MaxSampleValue))
		for ch := 0; ch < opt.Channels; ch++ {
			dit[sampleWidth*i+2*ch] = byte(b)
			dit[sampleWidth*i+1+2*ch] = byte(b >> 8)
		}
	}
	return cw
}

// Read a symbol from the sequence or return not found
func (cw *Generator) in() (symbol byte, found bool) {
	cw.sequenceMu.Lock()
	defer cw.sequenceMu.Unlock()
	if len(cw.sequence) <= 0 {
		return 0, false
	}
	symbol, cw.sequence = cw.sequence[0], cw.sequence[1:]
	return symbol, true
}

// Add things to output sequence, call with lock held
func (cw *Generator) _out(symbols ...byte) {
	cw.sequence = append(cw.sequence, symbols...)
}

// Clear empties the sequence and resets the state
func (cw *Generator) Clear() {
	cw.sequence = cw.sequence[:0]
	cw.sampleOffset = 0
}

// Time it should take to play the Morse
func (cw *Generator) duration() time.Duration {
	return time.Duration((float64(len(cw.sequence)) * wpmToDitTime(cw.opt.WPM)) * float64(time.Second))
}

// Read implements the io.Reader interface for the sound data
func (cw *Generator) Read(buf []byte) (n int, err error) {
	for len(buf) > 0 {
		if cw.sampleOffset == 0 {
			// What sample are we supposed to be playing
			var found bool
			cw.sampleIndex, found = cw.in()
			if !found {
				if cw.opt.Continuous {
					err = nil
				} else {
					err = io.EOF
				}
				break
			}
		}

		// Get its waveform
		sample := cw.samples[cw.sampleIndex]
		nn := copy(buf, sample[cw.sampleOffset:])
		if nn == 0 {
			break
		}
		n += nn
		cw.sampleOffset += nn
		if cw.sampleOffset >= cw.sampleLength {
			cw.sampleOffset = 0
		}
		buf = buf[nn:]
	}
	// fmt.Printf("n=%d, err=%v, sampleOffset=%d, sampleIndex=%d\n", n, err, cw.sampleOffset, cw.sampleIndex)
	return n, err
}

// Add Farnsworth spacing
func (cw *Generator) _extraDits() {
	for i := 0; i < cw.extraDits; i++ {
		cw._out(0)
	}
}

// Adds the rune to the output
func (cw *Generator) Rune(r rune) {
	cw.sequenceMu.Lock()
	defer cw.sequenceMu.Unlock()

	r = unicode.ToUpper(r)
	code := morseCode[r]
	if code == "" {
		if cw.opt.Debug {
			fmt.Printf("Don't know how to play '%c'\n", r)
		}
		return
	}
	for _, c := range code {
		switch c {
		case '-':
			cw._out(1, 1, 1, 0)
		case '.':
			cw._out(1, 0)
		case ' ':
			// word space is 7 dits
			// we've written 1 on the last dit/dah
			// and we'll write 2 after this
			// so need 4 more
			cw._out(0, 0, 0, 0)
			// And we have an extra Farnsworth space every word
			cw._extraDits()
		default:
			panic("Bad symbol in code")
		}
	}
	// write letter gap of 3 dits - have written one already
	cw._out(0, 0)
	cw._extraDits()
}

// Adds the string to the output
func (cw *Generator) String(s string) {
	for _, r := range s {
		cw.Rune(r)
	}
}

// check interfaces
var _ io.Reader = (*Generator)(nil)
