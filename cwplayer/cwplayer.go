package cwplayer

import (
	"time"

	"github.com/hajimehoshi/oto/v2"
	"github.com/ncw/ncwtester/cwgenerator"
)

// Player contains state for the morse generation
type Player struct {
	*cwgenerator.Generator
	context *oto.Context
	player  oto.Player
}

func New(opt cwgenerator.Options) (*Player, error) {
	context, ready, err := oto.NewContext(opt.SampleRate, opt.ChannelNum, opt.BitDepthInBytes)
	if err != nil {
		return nil, err
	}
	<-ready
	generator := cwgenerator.New(opt)
	return &Player{
		Generator: generator,
		context:   context,
		player:    context.NewPlayer(generator),
	}, nil
}

// Resets the audio
func (p *Player) Reset() {
	p.player.Reset()
}

// Plays what we have so far and syncs
func (p *Player) SyncPlay() {
	p.player.Play()
	for p.player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}
}
