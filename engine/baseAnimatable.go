package engine

import "apodeiktikos.com/fbtest/interfaces"

type BaseAnimatable struct {
	interfaces.Animatable
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
