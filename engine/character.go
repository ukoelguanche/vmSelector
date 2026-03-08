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
func (c *BaseMovable) MoveTo(target core.Point, duration time.Duration) {
	c.startPosition = c.position
	c.targetPosition = target
	c.startTime = time.Now()
	c.duration = duration
	c.moving = true
}

func (c *Character) GetEaseFunction() func(float64) float64 { return c.easeFunc }
func (c *Character) GetPosition() core.Point                { return c.position }
func (c *Character) GetTargetPosition() core.Point          { return c.targetPosition }

func (c *Character) GetSpeed() core.Size { return c.Speed }
func (c *Character) IsMoving() bool      { return c.moving }

func (c *Character) SetPosition(position core.Point) {
	c.position = position
}
func (c *Character) SetTargetPosition(targetPosition core.Point) {

	c.targetPosition = targetPosition

	c.moving = true
	return
}

func (c *Character) EndMovement() {
	if !c.moving {
		return
	}
	c.moving = false
	if c.onMovementComplete != nil {
		c.onMovementComplete(c)
	}
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
	//relativePaletteSwapSpeed := float32(0.07)
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
