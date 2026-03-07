package engine

import (
	"math"
	"time"

	"apodeiktikos.com/fbtest/interfaces"
	"github.com/ukoelguanche/graphicsengine/core"
)

type SpriteInstance struct {
	BaseMovable
	BaseAnimatable

	sprite *core.Sprite

	CurrentSequence         []int
	SequenceOffset          float32
	CurrentSequencePosition float32
	SequenceLength          int
	Scale                   float64
	Moving                  bool
	OnAnimationComplete     func(interfaces.Renderable)
	OnMovementComplete      func(interfaces.Renderable)
	totalDistance           float64
	easeFunc                func(float64) float64

	PaletteSwapIndex int
	Speed            core.Size
	AbsSpeed         float64

	Position       core.Point
	StartPosition  core.Point
	TargetPosition core.Point
	StartTime      time.Time
	Duration       time.Duration
	TotalDistance  float64
}

func (si *SpriteInstance) GetSprite() *core.Sprite {
	return si.sprite
}
func (si *SpriteInstance) SetEaseFunction(f func(float64) float64) { si.easeFunc = f }
func (si *SpriteInstance) GetStartTime() time.Time                 { return si.StartTime }
func (si *SpriteInstance) GetDuration() time.Duration              { return si.Duration }
func (si *SpriteInstance) GetStartPosition() core.Point            { return si.StartPosition }
func (si *SpriteInstance) MoveTo(target core.Point, duration time.Duration) {
	si.StartPosition = si.Position
	si.TargetPosition = target
	si.StartTime = time.Now()
	si.Duration = duration
	si.Moving = true
}
func (si *SpriteInstance) SetCurrentSequence(sequence []int) {
	si.CurrentSequence = sequence
	si.CurrentSequencePosition = 0
}

func (si *SpriteInstance) SetOnAnimationComplete(f func(interfaces.Renderable)) {
	si.OnAnimationComplete = f
}

func (si *SpriteInstance) SetOnMovementComplete(f func(interfaces.Renderable)) {
	si.OnMovementComplete = f
}
func (si *SpriteInstance) GetEaseFunction() func(float64) float64 { return si.easeFunc }
func (si *SpriteInstance) GetPosition() core.Point                { return si.Position }
func (si *SpriteInstance) GetTargetPosition() core.Point          { return si.TargetPosition }
func (si *SpriteInstance) GetSpeed() core.Size                    { return si.Speed }
func (si *SpriteInstance) IsMoving() bool                         { return si.Moving }

func (si *SpriteInstance) GetFrame(index int32) core.Rect { return si.sprite.Frames[index] }
func (si *SpriteInstance) GetSequences(sequenceName string) []int {
	return si.sprite.Sequences[sequenceName]
}

func (si *SpriteInstance) SetPosition(position core.Point) {
	si.Position = position
}
func (si *SpriteInstance) SetTargetPosition(targetPosition core.Point) {

	si.TargetPosition = targetPosition

	si.Moving = true
	si.totalDistance = math.Sqrt(math.Pow(targetPosition.X-si.Position.X, 2) + math.Pow(targetPosition.Y-si.Position.Y, 2))
	return
}

func (si *SpriteInstance) EndMovement() {
	if !si.Moving {
		return
	}
	si.Moving = false
	if si.OnMovementComplete != nil {
		si.OnMovementComplete(si)
	}
}

func (s *SpriteInstance) Draw(d interfaces.Drawer) {
	d.DrawSpriteRect(s.sprite, s.CurrentFrame(), s.Position)
}

func (s *SpriteInstance) NextFrame() {
	s.UpdatePosition(s)
	s.UpdateFrame(s)
}

func (s *SpriteInstance) CurrentFrame() core.Rect {
	frame := int(float32(len(s.CurrentSequence)) * s.CurrentSequencePosition)

	return s.sprite.Frames[s.CurrentSequence[frame]]
}

func (s *SpriteInstance) GetCurrentSequencePosition() float32 {
	return s.CurrentSequencePosition
}

func (s *SpriteInstance) SetCurrentSequencePosition(currentSequencePosition float32) {
	s.CurrentSequencePosition = currentSequencePosition
}

func (s *SpriteInstance) GetSequenceOffset() float32 {
	return s.SequenceOffset
}

func (s *SpriteInstance) ExecOnAnimationComplete() {
	if s.OnAnimationComplete != nil {
		s.OnAnimationComplete(s)
	}
}

func BuildSpriteInstance(sprites core.Sprites, name string, sequenceName string, position core.Point) *SpriteInstance {
	sequence := sprites.Sprites[name].Sequences[sequenceName]
	relativeSeqenceSpeed := float32(0.5)
	spriteInstance := &SpriteInstance{
		sprite:                  sprites.Sprites[name],
		Position:                position,
		TargetPosition:          position,
		CurrentSequence:         sequence,
		SequenceOffset:          1 / float32(len(sequence)) * relativeSeqenceSpeed,
		CurrentSequencePosition: 0.0,
		SequenceLength:          len(sequence),
		Moving:                  false,
	}

	return spriteInstance
}
