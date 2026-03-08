package engine

import (
	"time"

	"apodeiktikos.com/fbtest/interfaces"
	"github.com/ukoelguanche/graphicsengine/core"
)

type SpriteInstance struct {
	BaseMovable
	BaseAnimatable

	easeFunc func(float64) float64

	PaletteSwapIndex int
	AbsSpeed         float64
}

func (si *SpriteInstance) SetEaseFunction(f func(float64) float64) { si.easeFunc = f }
func (si *SpriteInstance) GetEaseFunction() func(float64) float64  { return si.easeFunc }

func (si *SpriteInstance) GetSprite() *core.Sprite {
	return si.sprite
}

func (si *SpriteInstance) MoveTo(target core.Point, duration time.Duration) {
	si.startPosition = si.position
	si.targetPosition = target
	si.startTime = time.Now()
	si.duration = duration
	si.moving = true
}

func (s *SpriteInstance) Draw(d interfaces.Drawer) {
	d.DrawSpriteRect(s.GetSprite(), s.GetCurrentFrame(), s.position)
}

func (s *SpriteInstance) NextFrame() {
	s.UpdatePosition(s)
	s.UpdateFrame(s)
}

func BuildSpriteInstance(sprites core.Sprites, name string, sequenceName string, position core.Point) *SpriteInstance {
	sequence := sprites.Sprites[name].Sequences[sequenceName]
	relativeSeqenceSpeed := float32(0.5)
	spriteInstance := &SpriteInstance{
		BaseAnimatable: BaseAnimatable{
			sprite:                  sprites.Sprites[name],
			currentSequence:         sequence,
			frameIdx:                0,
			sequenceOffset:          1 / float32(len(sequence)) * relativeSeqenceSpeed,
			currentSequencePosition: 0.0,
		},
		BaseMovable: BaseMovable{
			position:       position,
			startPosition:  position,
			targetPosition: position,
			moving:         false,
		},
	}

	return spriteInstance
}
