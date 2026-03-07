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

	currentSequence         []int
	SequenceOffset          float32
	currentSequencePosition float32
	moving                  bool
	onAnimationComplete     func(interfaces.Renderable)
	onMovementComplete      func(interfaces.Renderable)
	totalDistance           float64
	easeFunc                func(float64) float64

	paletteSwapIndex int
	speed            core.Size
	AbsSpeed         float64

	position       core.Point
	startPosition  core.Point
	targetPosition core.Point
	StartTime      time.Time
	duration       time.Duration
}

func (si *SpriteInstance) GetSprite() *core.Sprite {
	return si.sprite
}
func (si *SpriteInstance) SetEaseFunction(f func(float64) float64) { si.easeFunc = f }
func (si *SpriteInstance) GetStartTime() time.Time                 { return si.StartTime }
func (si *SpriteInstance) GetDuration() time.Duration              { return si.duration }
func (si *SpriteInstance) GetStartPosition() core.Point            { return si.startPosition }
func (si *SpriteInstance) MoveTo(target core.Point, duration time.Duration) {
	si.startPosition = si.position
	si.targetPosition = target
	si.StartTime = time.Now()
	si.duration = duration
	si.moving = true
}
func (si *SpriteInstance) SetCurrentSequence(sequence []int) {
	si.currentSequence = sequence
	si.currentSequencePosition = 0
}

func (si *SpriteInstance) SetOnAnimationComplete(f func(interfaces.Renderable)) {
	si.onAnimationComplete = f
}

func (si *SpriteInstance) SetOnMovementComplete(f func(interfaces.Renderable)) {
	si.onMovementComplete = f
}
func (si *SpriteInstance) GetEaseFunction() func(float64) float64 { return si.easeFunc }
func (si *SpriteInstance) GetPosition() core.Point                { return si.position }
func (si *SpriteInstance) GetTargetPosition() core.Point          { return si.targetPosition }
func (si *SpriteInstance) GetSpeed() core.Size                    { return si.speed }
func (si *SpriteInstance) IsMoving() bool                         { return si.moving }

func (si *SpriteInstance) GetFrame(index int32) core.Rect { return si.sprite.Frames[index] }
func (si *SpriteInstance) GetSequences(sequenceName string) []int {
	return si.sprite.Sequences[sequenceName]
}

func (si *SpriteInstance) SetPosition(position core.Point) {
	si.position = position
}
func (si *SpriteInstance) SetTargetPosition(targetPosition core.Point) {

	si.targetPosition = targetPosition

	si.moving = true
	si.totalDistance = math.Sqrt(math.Pow(targetPosition.X-si.position.X, 2) + math.Pow(targetPosition.Y-si.position.Y, 2))
	return
}

func (si *SpriteInstance) EndMovement() {
	if !si.moving {
		return
	}
	si.moving = false
	if si.onMovementComplete != nil {
		si.onMovementComplete(si)
	}
}

func (s *SpriteInstance) Draw(d interfaces.Drawer) {
	d.DrawSpriteRect(s.sprite, s.CurrentFrame(), s.position)
}

func (s *SpriteInstance) NextFrame() {
	s.UpdatePosition(s)
	s.UpdateFrame(s)
}

func (s *SpriteInstance) CurrentFrame() core.Rect {
	frame := int(float32(len(s.currentSequence)) * s.currentSequencePosition)

	return s.sprite.Frames[s.currentSequence[frame]]
}

func (s *SpriteInstance) GetCurrentSequencePosition() float32 {
	return s.currentSequencePosition
}

func (s *SpriteInstance) SetCurrentSequencePosition(currentSequencePosition float32) {
	s.currentSequencePosition = currentSequencePosition
}

func (s *SpriteInstance) GetSequenceOffset() float32 {
	return s.SequenceOffset
}

func (s *SpriteInstance) ExecOnAnimationComplete() {
	if s.onAnimationComplete != nil {
		s.onAnimationComplete(s)
	}
}

func BuildSpriteInstance(sprites core.Sprites, name string, sequenceName string, position core.Point) *SpriteInstance {
	sequence := sprites.Sprites[name].Sequences[sequenceName]
	relativeSeqenceSpeed := float32(0.5)
	spriteInstance := &SpriteInstance{
		sprite:                  sprites.Sprites[name],
		position:                position,
		targetPosition:          position,
		currentSequence:         sequence,
		SequenceOffset:          1 / float32(len(sequence)) * relativeSeqenceSpeed,
		currentSequencePosition: 0.0,
		moving:                  false,
	}

	return spriteInstance
}
