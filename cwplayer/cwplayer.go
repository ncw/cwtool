package cwplayer

import (
	"time"

	"github.com/hajimehoshi/oto/v2"
	"github.com/ncw/cwtool/cw"
	"github.com/ncw/cwtool/cwgenerator"
)

// Player contains state for the morse generation
type Player struct {
	generator *cwgenerator.Generator
	opt       *cw.Options
	context   *oto.Context
	player    oto.Player
}

func New(opt *cw.Options) (*Player, error) {
	context, ready, err := oto.NewContext(opt.SampleRate, opt.ChannelNum, opt.BitDepthInBytes)
	if err != nil {
		return nil, err
	}
	<-ready
	generator := cwgenerator.New(opt)
	p := &Player{
		generator: generator,
		context:   context,
		player:    context.NewPlayer(generator),
	}
	p.player.Reset()
	return p, nil
}

// kick the player into action
func (p *Player) kick() {
	// There is a race condition here if the player stops playing
	// after the test, however we kick again in Sync which should
	// pick it up.
	if !p.player.IsPlaying() {
		p.player.Reset()
	}
	p.player.Play()
}

// Rune adds r to the output
func (p *Player) Rune(r rune) {
	p.generator.Rune(r)
	p.kick()
}

// String adds s to the output
func (p *Player) String(s string) {
	p.generator.String(s)
	p.kick()
}

// Sync by waiting for all the morse to be played
func (p *Player) Sync() {
	p.kick()
	for p.player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}
}

// Close the output
func (p *Player) Close() error {
	return nil
}

// Check interface
var _ cw.CW = (*Player)(nil)
