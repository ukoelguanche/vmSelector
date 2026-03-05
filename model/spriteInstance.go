package model

import "math"

type SpriteInstance struct {
	Sprite         *Sprite
	Position       Point
	TargetPosition Point

	FrameIdx                int
	CurrentSequence         []int
	SequenceOffset          float32
	CurrentSequencePosition float32
	SequenceLength          int
	Scale                   float64
	moving                  bool
	OnAnimationComplete     func(*SpriteInstance)
	OnMovementComplete      func(*SpriteInstance)
	totalDistance           float64
	easeFunc                func(float64) float64

	CurrentPalleteSwapOffset   float32
	CurrentPalleteSwapPosition float32

	PaletteSwapIndex int
	Speed            Size
}

func (si *SpriteInstance) GetSprite() *Sprite {
	return si.Sprite
}
func (si *SpriteInstance) GetBitmap() *Bitmap {
	return si.Sprite.Bitmap
}
func (si *SpriteInstance) SetPosition(position Point) {
	si.Position = position
}
func (si *SpriteInstance) SetTargetPosition(targetPosition Point, speed Size) {
	si.TargetPosition = targetPosition
	si.Speed = speed
	si.moving = true
	si.totalDistance = math.Sqrt(math.Pow(targetPosition.X-si.Position.X, 2) + math.Pow(targetPosition.Y-si.Position.Y, 2))
	return
}
func (si *SpriteInstance) SetEaseFunction(f func(float64) float64) { si.easeFunc = f }
func (si *SpriteInstance) GetEaseFunction() func(float64) float64  { return si.easeFunc }
func (si *SpriteInstance) GetTotalDistance() float64               { return si.totalDistance }
func (si *SpriteInstance) GetPosition() Point                      { return si.Position }
func (si *SpriteInstance) GetTargetPosition() Point                { return si.TargetPosition }
func (si *SpriteInstance) GetSpeed() Size                          { return si.Speed }
func (si *SpriteInstance) IsMoving() bool                          { return si.moving }
func (si *SpriteInstance) EndMovement() {
	if !si.moving {
		return
	}
	si.moving = false
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
		moving:                     false,
	}

	return spriteInstance

}
