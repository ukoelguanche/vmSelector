package engine

import (
	"apodeiktikos.com/fbtest/interfaces"
)

func SetCurrentSequenceByName(b interfaces.Animatable, name string) {
	sequence := b.GetSprite().Sequences[name]
	b.SetCurrentSequence(sequence)
	b.SetCurrentSequencePosition(0.0)
}
