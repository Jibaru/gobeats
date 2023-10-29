package config

import (
	"github.com/ebitengine/oto/v3"
)

func NewOtoContext() (*oto.Context, chan struct{}, error) {
	op := &oto.NewContextOptions{}
	op.SampleRate = 44100
	op.ChannelCount = 2
	op.Format = oto.FormatSignedInt16LE

	return oto.NewContext(op)
}
