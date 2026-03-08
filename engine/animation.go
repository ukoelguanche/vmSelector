package engine

import (
	"apodeiktikos.com/fbtest/interfaces"
	"github.com/ukoelguanche/graphicsengine/core"
)

func GetFrame(saa interfaces.Animatable, index int32) core.Frame {
	return saa.GetSprite().Frames[index]
}
func GetCurrentFrame(saa interfaces.Animatable) core.Frame {
	sequence := saa.GetCurrentSequence()
	frame := int(float32(len(sequence)) * saa.GetCurrentSequencePosition())

	return saa.GetSprite().Frames[sequence[frame]]
}

func SetCurrentSequenceByName(b interfaces.Animatable, name string) {
	sequence := b.GetSprite().Sequences[name]
	b.SetCurrentSequence(sequence)
}
