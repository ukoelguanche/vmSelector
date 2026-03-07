package engine

import (
	"math"
	"time"

	"apodeiktikos.com/fbtest/interfaces"
	"github.com/ukoelguanche/graphicsengine/core"
)

type Character struct {
	Sprite *core.Sprite

	FrameIdx                int
	CurrentSequence         []int
	SequenceOffset          float32
	CurrentSequencePosition float32
	SequenceLength          int
	Scale                   float64
	Moving                  bool
	OnAnimationComplete     func(*Character)
	OnMovementComplete      func(interfaces.Renderable)
	movementFrameCount      float64
	movementFrame           float64
	easeFunc                func(float64) float64

	PaletteSwapIndex int
	// Speed            core.Size
	AbsSpeed float64

	Position       core.Point
	StartPosition  core.Point
	TargetPosition core.Point
	StartTime      time.Time
	Duration       time.Duration
	TotalDistance  float64

	Acceleration float64
	MaxSpeed     float64
	Speed        core.Size
}

func (c *Character) GetSprite() *core.Sprite {
	return c.Sprite
}
func (c *Character) SetEaseFunction(f func(float64) float64) { c.easeFunc = f }
func (c *Character) GetStartTime() time.Time                 { return c.StartTime }
func (c *Character) GetDuration() time.Duration              { return c.Duration }
func (c *Character) GetStartPosition() core.Point            { return c.StartPosition }
func (c *Character) MoveTo(target core.Point, duration time.Duration) {
	c.StartPosition = c.Position
	c.TargetPosition = target
	c.StartTime = time.Now()
	c.Duration = duration
	c.Moving = true
}

func (c *Character) SetOnMovementComplete(f func(interfaces.Renderable)) { c.OnMovementComplete = f }
func (c *Character) SetOnAnimationComplete(f func(*Character))           { c.OnAnimationComplete = f }
func (c *Character) GetEaseFunction() func(float64) float64              { return c.easeFunc }
func (c *Character) GetMovementFrameCount() float64                      { return c.movementFrameCount }
func (c *Character) GetMovementFrame() float64                           { return c.movementFrame }
func (c *Character) GetPosition() core.Point                             { return c.Position }
func (c *Character) GetTargetPosition() core.Point                       { return c.TargetPosition }

func (c *Character) GetSpeed() core.Size { return c.Speed }
func (c *Character) IsMoving() bool      { return c.Moving }

func (c *Character) SetPosition(position core.Point) {
	c.Position = position
}
func (c *Character) SetTargetPosition(targetPosition core.Point) {

	c.TargetPosition = targetPosition

	c.Moving = true
	return
}

func (c *Character) SetSpeed(absSpeed float64) {
	dx := c.TargetPosition.X - c.Position.X
	dy := c.TargetPosition.Y - c.Position.Y
	angle := math.Atan2(dy, dx)

	//c.movementFrameCount = c.totalDistance / absSpeed
	c.movementFrame = 0

	c.AbsSpeed = absSpeed
	c.Speed = core.Size{W: absSpeed * math.Cos(angle), H: absSpeed * math.Sin(angle)}
}

func (c *Character) EndMovement() {
	if !c.Moving {
		return
	}
	c.Moving = false
	if c.OnMovementComplete != nil {
		c.OnMovementComplete(c)
	}
}

func (s *Character) Draw(d interfaces.Drawer) {
	d.DrawSpriteRect(s.Sprite, s.CurrentFrame(), s.Position)
}

func (s *Character) NextFrame() {
	interfaces.UpdatePosition(s)

	// Update Frame
	s.CurrentSequencePosition += s.SequenceOffset
	if s.CurrentSequencePosition >= 1 {
		// ToDo: avoid loop if not needed
		s.CurrentSequencePosition = 0

		if s.OnAnimationComplete != nil {
			s.OnAnimationComplete(s)
		}

	}

}

func (s *Character) CurrentFrame() core.Rect {
	frame := int(float32(len(s.CurrentSequence)) * s.CurrentSequencePosition)

	return s.Sprite.Frames[s.CurrentSequence[frame]]
}

func BuildCharacter(sprites core.Sprites, name string, sequenceName string, position core.Point) *Character {
	sequence := sprites.Sprites[name].Sequences[sequenceName]
	relativeSeqenceSpeed := float32(0.5)
	//relativePaletteSwapSpeed := float32(0.07)
	spriteInstance := &Character{
		Sprite:                  sprites.Sprites[name],
		Position:                position,
		TargetPosition:          position,
		FrameIdx:                0,
		CurrentSequence:         sequence,
		SequenceOffset:          1 / float32(len(sequence)) * relativeSeqenceSpeed,
		CurrentSequencePosition: 0.0,
		SequenceLength:          len(sequence),
		Moving:                  false,
	}

	return spriteInstance
}
