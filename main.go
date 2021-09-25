package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
	"unicode"

	"github.com/fatih/color"
	"github.com/hajimehoshi/oto/v2"
	"golang.org/x/term"
)

var (
	sampleRate = flag.Int("samplerate", 44100, "sample rate")
	wpm        = flag.Float64("wpm", 25.0, "WPM to send at")
	frequency  = flag.Float64("frequency", 600.0, "HZ of morse")
	logFile    = flag.String("log", "ncwtesterstats.csv", "CSV file to log attempts")
	timeCutoff = flag.Duration("cutoff", 0, "If set, ignore stats older than this")
)

const (
	channelNum      = 2
	bitDepthInBytes = 2
	maxSampleValue  = 32767
)

type cwPlayer struct {
	sequence      []byte    // sequence to play samples in
	sampleLength  int       // length of sample in bytes
	samples       [2][]byte // samples to play
	sequenceIndex int       // which sample we are playing
	sampleOffset  int       // how far we've got through that sample
	eofTime       time.Time
}

func newCWPlayer() *cwPlayer {
	cw := &cwPlayer{}

	ditTimeSeconds := wpmToDitTime(*wpm)
	cyclesPerDit := *frequency * ditTimeSeconds
	// fmt.Printf("cyclesPerDit = %.3f at %.1f Hz\n", cyclesPerDit, *frequency)
	// Round cycles per dit to an exact number to avoid clicks
	// this changes the frequency slightly
	cyclesPerDit = math.Round(cyclesPerDit)
	// newFrequency := cyclesPerDit / ditTimeSeconds
	// fmt.Printf("cyclesPerDit = %.3f at %.1f Hz\n", cyclesPerDit, newFrequency)

	samplesPerDit := int(math.Round(float64(*sampleRate) * ditTimeSeconds))
	sampleWidth := channelNum * bitDepthInBytes
	cw.sampleLength = samplesPerDit * sampleWidth
	cw.samples[0] = make([]byte, cw.sampleLength)
	cw.samples[1] = make([]byte, cw.sampleLength)
	dit := cw.samples[1]
	for i := 0; i < samplesPerDit; i++ {
		b := int16(math.Sin(2*math.Pi*float64(i)/float64(samplesPerDit)*cyclesPerDit) * 0.3 * maxSampleValue)
		for ch := 0; ch < channelNum; ch++ {
			dit[sampleWidth*i+2*ch] = byte(b)
			dit[sampleWidth*i+1+2*ch] = byte(b >> 8)
		}
	}
	return cw
}

// Empties the sequence and resets the state
func (cw *cwPlayer) reset() {
	cw.sequence = cw.sequence[:0]
	cw.sequenceIndex = 0
	cw.sampleOffset = 0
}

// Time it should take to play the morse
func (cw *cwPlayer) duration() time.Duration {
	return time.Duration((float64(len(cw.sequence)) * wpmToDitTime(*wpm)) * float64(time.Second))
}

// Convert a duration into milliseconds
func ms(t time.Duration) int64 {
	return t.Milliseconds()
}

func (cw *cwPlayer) Read(buf []byte) (n int, err error) {
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
func (cw *cwPlayer) rune(r rune) {
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
func (cw *cwPlayer) string(s string) {
	for _, r := range s {
		cw.rune(r)
	}
}

func shuffleString(s string) string {
	rs := []rune(s)
	rand.Shuffle(len(rs), func(i, j int) {
		rs[i], rs[j] = rs[j], rs[i]
	})
	return string(rs)
}

// Returns whether the character is an exit character, eg CTRL-C or ESC
func isExit(r rune) bool {
	return r == 0x03 || r == 0x1B
}

// Reads a single character from the terminal
func getChar() (r rune) {
	s, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Failed to MakeRaw: %v", err)
	}
	defer func() {
		err := term.Restore(int(os.Stdin.Fd()), s)
		if err != nil {
			log.Fatalf("Failed to Restore: %v", err)
		}
	}()
	var buf [1]byte
	n, err := os.Stdin.Read(buf[:])
	if err != nil {
		log.Fatalf("Failed to Read: %v", err)
	}
	if n != 1 {
		log.Fatalf("Didn't read exactly 1 character")
	}
	return unicode.ToLower(rune(buf[0]))
}

func yorn(prompt string) bool {
	fmt.Printf("%s (y/n)> ", prompt)
	var c rune
	for {
		c = getChar()
		if c == 'y' || c == 'n' {
			break
		} else if isExit(c) {
			fmt.Println("...bye\n")
			os.Exit(0)
		}
	}
	fmt.Println(string(c))
	return c == 'y'
}

// Plays the player and waits for it to finish
func syncPlay(p oto.Player) {
	p.Play()
	for p.IsPlaying() {
		time.Sleep(time.Millisecond)
	}
}

func run() error {
	c, ready, err := oto.NewContext(*sampleRate, channelNum, bitDepthInBytes)
	if err != nil {
		return err
	}
	<-ready

	cw := newCWPlayer()
	p := c.NewPlayer(cw)
	cw.string(" vvv")
	syncPlay(p)

	csvLog := NewCSVLog(*logFile)

	letters := "abcdefghijklmnopqrstuvwxyz0123456789.=/,?"

outer:
	for {
		if !yorn(fmt.Sprintf("Start test round with %d letters?", len(letters))) {
			break outer
		}
		for i, tx := range shuffleString(letters) {
			p.Reset()
			cw.reset()
			cw.rune(' ')
			cw.rune(tx)
			// cwDuration := cw.duration()
			// startPlaying := time.Now()
			syncPlay(p)
			finishedPlaying := time.Now()
			// fmt.Printf("time to play %dms, expected %dms, diff=%dms\n", ms(finishedPlaying.Sub(startPlaying)), ms(cwDuration), ms(finishedPlaying.Sub(startPlaying)-cwDuration))

			rx := getChar()
			if isExit(rx) {
				break outer
			}
			reactionTime := time.Since(finishedPlaying)
			ok := rx == tx
			fmt.Printf("%2d/%2d: %c: reaction time %5dms: ", i+1, len(letters), tx, ms(reactionTime))
			if ok {
				color.Green("OK\n")
			} else {
				color.Red(fmt.Sprintf("BAD %c\n", rx))
			}
			csvLog.Add(tx, rx, reactionTime)
		}
	}

	stats := NewStats(*logFile, *timeCutoff)
	stats.Summary()

	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()
	if err := run(); err != nil {
		panic(err)
	}
}
