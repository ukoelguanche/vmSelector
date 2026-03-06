package model

import (
	"math"
	"time"
)

type SpriteInstance struct {
	Sprite *Sprite

	FrameIdx                int
	CurrentSequence         []int
	SequenceOffset          float32
	CurrentSequencePosition float32
	SequenceLength          int
	Scale                   float64
	Moving                  bool
	OnAnimationComplete     func(*SpriteInstance)
	OnMovementComplete      func(Renderable)
	totalDistance           float64
	movementFrameCount      float64
	movementFrame           float64
	easeFunc                func(float64) float64

	CurrentPalleteSwapOffset   float32
	CurrentPalleteSwapPosition float32

	PaletteSwapIndex int
	Speed            Size
	AbsSpeed         float64

	Position       Point
	StartPosition  Point
	TargetPosition Point
	StartTime      time.Time
	Duration       time.Duration
	TotalDistance  float64
}

func (si *SpriteInstance) GetSprite() *Sprite {
	return si.Sprite
}
func (si *SpriteInstance) GetBitmap() *Bitmap {
	return si.Sprite.Bitmap
}
func (si *SpriteInstance) SetEaseFunction(f func(float64) float64) { si.easeFunc = f }
func (si *SpriteInstance) GetStartTime() time.Time                 { return si.StartTime }
func (si *SpriteInstance) GetDuration() time.Duration              { return si.Duration }
func (si *SpriteInstance) GetStartPosition() Point                 { return si.StartPosition }
func (si *SpriteInstance) MoveTo(target Point, duration time.Duration) {
	si.StartPosition = si.Position
	si.TargetPosition = target
	si.StartTime = time.Now()
	si.Duration = duration
	si.Moving = true
}
func (si *SpriteInstance) SetOnMovementComplete(f func(Renderable)) { si.OnMovementComplete = f }
func (si *SpriteInstance) GetEaseFunction() func(float64) float64   { return si.easeFunc }
func (si *SpriteInstance) GetMovementFrameCount() float64           { return si.movementFrameCount }
func (si *SpriteInstance) GetMovementFrame() float64                { return si.movementFrame }
func (si *SpriteInstance) GetTotalDistance() float64                { return si.totalDistance }
func (si *SpriteInstance) GetPosition() Point                       { return si.Position }
func (si *SpriteInstance) GetTargetPosition() Point                 { return si.TargetPosition }
func (si *SpriteInstance) GetSpeed() Size                           { return si.Speed }
func (si *SpriteInstance) IsMoving() bool                           { return si.Moving }

func (si *SpriteInstance) SetPosition(position Point) {
	si.Position = position
}
func (si *SpriteInstance) SetTargetPosition(targetPosition Point) {

	si.TargetPosition = targetPosition

	si.Moving = true
	si.totalDistance = math.Sqrt(math.Pow(targetPosition.X-si.Position.X, 2) + math.Pow(targetPosition.Y-si.Position.Y, 2))
	return
}

func (si *SpriteInstance) SetSpeed(absSpeed float64) {
	dx := si.TargetPosition.X - si.Position.X
	dy := si.TargetPosition.Y - si.Position.Y
	angle := math.Atan2(dy, dx)

	si.movementFrameCount = si.totalDistance / absSpeed
	si.movementFrame = 0

	si.AbsSpeed = absSpeed
	si.Speed = Size{W: absSpeed * math.Cos(angle), H: absSpeed * math.Sin(angle)}
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

func (si *SpriteInstance) ProcessColor(color []byte) []byte {
	if si.Sprite.PaletteSwap.TargetPalette != nil {
		return ReplacePalette(color,
			si.Sprite.PaletteSwap.SourcePalette,
			si.Sprite.PaletteSwap.TargetPalette,
			si.CurrentSwapPaletteIndex())
	}
	return color
}

func (s *SpriteInstance) NextFrame() {
	UpdatePosition(s)

	// Update Frame
	s.CurrentSequencePosition += s.SequenceOffset
	if s.CurrentSequencePosition >= 1 {
		// ToDo: avoid loop if not needed
		s.CurrentSequencePosition = 0

		if s.OnAnimationComplete != nil {
			s.OnAnimationComplete(s)
		}

	}

	// Swap palettes
	if s.Sprite.PaletteSwap.TargetPalette == nil {
		return
	}

	s.CurrentPalleteSwapPosition += s.CurrentPalleteSwapOffset
	if s.CurrentPalleteSwapPosition >= 1 {
		s.CurrentPalleteSwapPosition = 0
	}
}

func (s *SpriteInstance) CurrentFrame() Rect {
	frame := int(float32(len(s.CurrentSequence)) * s.CurrentSequencePosition)

	return s.Sprite.Frames[s.CurrentSequence[frame]]
}

func (s *SpriteInstance) CurrentSwapPaletteIndex() int {
	return int(float32(len(*s.Sprite.PaletteSwap.TargetPalette)) * s.CurrentPalleteSwapPosition)
}

func BuildSpriteInstance(sprites Sprites, name string, sequenceName string, position Point) *SpriteInstance {
	sequence := sprites.Sprites[name].Sequences[sequenceName]
	relativeSeqenceSpeed := float32(0.5)
	relativePaletteSwapSpeed := float32(0.07)
	spriteInstance := &SpriteInstance{
		Sprite:                     sprites.Sprites[name],
		Position:                   position,
		TargetPosition:             position,
		FrameIdx:                   0,
		CurrentSequence:            sequence,
		SequenceOffset:             1 / float32(len(sequence)) * relativeSeqenceSpeed,
		CurrentSequencePosition:    0.0,
		SequenceLength:             len(sequence),
		CurrentPalleteSwapOffset:   1 / float32(len(sequence)) * relativePaletteSwapSpeed,
		CurrentPalleteSwapPosition: 0.0,
		Moving:                     false,
	}

	return spriteInstance

}
