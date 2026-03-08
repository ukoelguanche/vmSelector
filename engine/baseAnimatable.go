package engine

import (
	"apodeiktikos.com/fbtest/interfaces"
	"github.com/ukoelguanche/graphicsengine/core"
)

type BaseAnimatable struct {
	interfaces.Animatable

	sprite *core.Sprite

	frameIdx                int
	currentSequence         []int
	sequenceOffset          float32
	currentSequencePosition float32
	onAnimationComplete     func(interfaces.Animatable)
}

func (b *BaseAnimatable) GetSprite() *core.Sprite        { return b.sprite }
func (b *BaseAnimatable) GetSequences(name string) []int { return b.sprite.Sequences[name] }

func (b *BaseAnimatable) GetCurrentSequence() []int { return b.currentSequence }
func (s *BaseAnimatable) GetCurrentSequencePosition() float32 {
	return s.currentSequencePosition
}
func (s *BaseAnimatable) SetCurrentSequencePosition(csp float32) { s.currentSequencePosition = csp }

func (s *BaseAnimatable) GetSequenceOffset() float32 {
	return s.sequenceOffset
}

func (c *BaseAnimatable) SetOnAnimationComplete(f func(animatable interfaces.Animatable)) {
	c.onAnimationComplete = f
}
func (s *BaseAnimatable) ExecOnAnimationComplete() {
	if s.onAnimationComplete != nil {
		s.onAnimationComplete(s)
	}
}

func (si *BaseAnimatable) SetCurrentSequence(sequence []int) {
	si.currentSequence = sequence
	si.currentSequencePosition = 0
}

func (b *BaseAnimatable) GetCurrentFrame() core.Frame {
	frame := int(float32(len(b.currentSequence)) * b.currentSequencePosition)

	return b.sprite.Frames[b.currentSequence[frame]]
}

func (b *BaseAnimatable) UpdateFrame(a interfaces.Animatable) {
	currentSequencePosition := a.GetCurrentSequencePosition()
	currentSequencePosition += a.GetSequenceOffset()

	if currentSequencePosition >= 1 {
		currentSequencePosition = 0
		a.ExecOnAnimationComplete()
	}

	a.SetCurrentSequencePosition(currentSequencePosition)
}
