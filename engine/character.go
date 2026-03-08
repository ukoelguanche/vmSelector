package engine

import (
	"time"

	"apodeiktikos.com/fbtest/interfaces"
	"github.com/ukoelguanche/graphicsengine/core"
)

type Character struct {
	BaseMovable
	BaseAnimatable

	easeFunc func(float64) float64

	PaletteSwapIndex int
	AbsSpeed         float64
}

func (c *Character) SetEaseFunction(f func(float64) float64) { c.easeFunc = f }
func (c *Character) GetEaseFunction() func(float64) float64  { return c.easeFunc }

func (c *BaseMovable) MoveTo(target core.Point, duration time.Duration) {
	c.startPosition = c.position
	c.targetPosition = target
	c.startTime = time.Now()
	c.duration = duration
	c.moving = true
}

func (s *Character) Draw(d interfaces.Drawer) {
	d.DrawSpriteRect(s.GetSprite(), s.GetCurrentFrame(), s.position)
}

func (c *Character) NextFrame() {
	c.UpdatePosition(c)
	c.UpdateFrame(c)
}

func BuildCharacter(sprites core.Sprites, name string, sequenceName string, position core.Point) *Character {
	sequence := sprites.Sprites[name].Sequences[sequenceName]
	relativeSeqenceSpeed := float32(0.5)
	character := &Character{
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

	return character
}
