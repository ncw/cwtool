package cwfile

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/ncw/cwtool/cw"
	"github.com/ncw/cwtool/cwgenerator"
)

// Player contains state for the morse generation
type Player struct {
	generator *cwgenerator.Generator
	opt       *cw.Options
	out       io.WriteCloser
	encoder   *wav.Encoder
	buf       []byte          // raw data buffer
	abuf      []int           // int sample buffer
	aibuf     audio.IntBuffer // buffer to send to output
}

func New(opt *cw.Options) (*Player, error) {
	generator := cwgenerator.New(opt)

	// Destination file
	out, err := os.Create(opt.OutputFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't create wav output file: %w", err)
	}

	// setup the encoder
	encoder := wav.NewEncoder(out,
		opt.SampleRate,
		8*opt.BitDepthInBytes,
		opt.Channels,
		1, // PCM format
	)
	software := strings.Join(os.Args, " ")
	title := opt.Title
	if title == "" {
		title = software
	}
	encoder.Metadata = &wav.Metadata{
		Title:    title,
		Product:  title,
		Software: software,
		Artist:   "cwtool",
	}

	p := &Player{
		generator: generator,
		opt:       opt,
		out:       out,
		encoder:   encoder,
	}

	// Create buffers for writing to file
	const bufSize = 64 * 1024
	p.buf = make([]byte, p.opt.BitDepthInBytes*bufSize)
	p.abuf = make([]int, bufSize)
	p.aibuf = audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: p.opt.Channels,
			SampleRate:  p.opt.SampleRate,
		},
		Data:           p.abuf[:0],
		SourceBitDepth: 8 * p.opt.BitDepthInBytes,
	}

	return p, nil
}

// Rune adds r to the output
func (p *Player) Rune(r rune) {
	p.generator.Rune(r)
}

// String adds s to the output
func (p *Player) String(s string) {
	p.generator.String(s)
}

// Sync the morse so far to the file
func (p *Player) Sync() {
	for {
		n, err := p.generator.Read(p.buf)
		if err != io.EOF && err != nil {
			log.Printf("read audio failed: %v", err)
			break
		}
		isEOF := err == io.EOF
		// Convert into ints for encoding
		samples := n / p.opt.BitDepthInBytes
		for i := 0; i < samples; i++ {
			// FIXME assumes signed 16 bit
			p.abuf[i] = int(int16(binary.LittleEndian.Uint16(p.buf[2*i : 2*i+2])))
		}
		p.aibuf.Data = p.abuf[:samples]

		err = p.encoder.Write(&p.aibuf)
		if err != nil {
			log.Printf("write audio failed: %v", err)
			break
		}

		if isEOF {
			break
		}
	}
}

// Close the output
func (p *Player) Close() error {
	p.Sync()
	return errors.Join(
		p.encoder.Close(),
		p.out.Close(),
	)
}

// Check interface
var _ cw.CW = (*Player)(nil)
