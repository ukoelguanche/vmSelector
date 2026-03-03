package model

type Point struct {
	X, Y int32
}

type Size struct {
	W, H int32
}

type Rect struct {
	Point Point
	Size  Size
}

type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

type Palette []Color

type PaletteSwap struct {
	SourcePaletteName string
	TargetPaletteName string
	SourcePalette     *Palette
	TargetPalette     *Palette
}

type Sprite struct {
	Name         string
	BitmapSource string `json:"BitmapSource"`
	Bitmap       *Bitmap
	Frames       []Rect           `json:"Frames"`
	Sequences    map[string][]int `json:"Sequences"`
	Characters   map[string]int   `json:"Characters"`
	PaletteSwap  PaletteSwap      `json:"PaletteSwap"`
}

type Sprites struct {
	BitmapSources map[string]string   `json:"BitmapSources"`
	Sprites       map[string]*Sprite  `json:"sprites"`
	Palettes      map[string]*Palette `json:"Palettes"`
}

type Renderable interface {
	GetBitmap() *Bitmap
	GetSprite() *Sprite
	ProcessColor(color []byte) []byte
	NextFrame()
}

type SpriteInstance struct {
	Sprite                  *Sprite
	Position                Point
	TargetPosition          Point
	FrameIdx                int
	CurrentSequence         []int
	SequenceOffset          float32
	CurrentSequencePosition float32
	SequenceLength          int
	Scale                   float64

	CurrentPalleteSwapOffset   float32
	CurrentPalleteSwapPosition float32

	PaletteSwapIndex int
	Speed            int32
}

func (si *SpriteInstance) GetSprite() *Sprite {
	return si.Sprite
}

func (si *SpriteInstance) GetBitmap() *Bitmap {
	return si.Sprite.Bitmap
}

func ReplacePalette(color []byte, sourcePalette *Palette, targetPalette *Palette, animationIndex int) []byte {
	if sourcePalette == nil || targetPalette == nil {
		return color
	}

	gradientIndex := sourcePalette.GradientIndex(color)

	if gradientIndex >= 0 {
		return (*targetPalette)[(gradientIndex+animationIndex)%len(*targetPalette)].Byte()
	}
	return color
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

type Text struct {
	Sprite         *Sprite
	Position       Point
	TargetPosition Point
	Text           string
}

func (t *Text) GetSprite() *Sprite {
	return t.Sprite
}

func (t *Text) GetBitmap() *Bitmap {
	return t.Sprite.Bitmap
}

func (t *Text) NextFrame() {
	UpdatePositionT(t)
}

func (t *Text) ProcessColor(color []byte) []byte {
	return color
}

func BuildSpriteInstance(sprites Sprites, name string, sequenceName string, position Point) *SpriteInstance {
	sequence := sprites.Sprites[name].Sequences[sequenceName]
	relativeSeqenceSpeed := float32(0.5)
	relativePaletteSwapSpeed := float32(0.2)
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
	}

	return spriteInstance

}

func BuildTextInstance(sprite *Sprite, text string, position Point) *Text {
	return &Text{Sprite: sprite, Text: text, Position: position, TargetPosition: position}
}

func (s *SpriteInstance) NextFrame() {
	UpdatePosition(s)

	// Update Frame
	s.CurrentSequencePosition += s.SequenceOffset
	if s.CurrentSequencePosition >= 1 {
		s.CurrentSequencePosition = 0
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

func UpdatePosition(s *SpriteInstance) {
	dx := s.TargetPosition.X - s.Position.X
	dy := s.TargetPosition.Y - s.Position.Y

	if dx > 0 {
		s.Position.X += s.Speed
	} else if dx < 0 {
		s.Position.X -= s.Speed
	}
	if dy > 0 {
		s.Position.Y += s.Speed
	} else if dy < 0 {
		s.Position.Y -= s.Speed
	}
}

func UpdatePositionT(s *Text) {
	dx := s.TargetPosition.X - s.Position.X
	dy := s.TargetPosition.Y - s.Position.Y

	speed := int32(2)

	if dx > 0 {
		s.Position.X += speed
	} else if dx < 0 {
		s.Position.X -= speed
	}
	if dy > 0 {
		s.Position.Y += speed
	} else if dy < 0 {
		s.Position.Y -= speed
	}
}

func (s *SpriteInstance) CurrentFrame() Rect {
	frame := int(float32(len(s.CurrentSequence)) * s.CurrentSequencePosition)
	return s.Sprite.Frames[s.CurrentSequence[frame]]
}

func (s *SpriteInstance) CurrentSwapPaletteIndex() int {
	return int(float32(len(*s.Sprite.PaletteSwap.TargetPalette)) * s.CurrentPalleteSwapPosition)
}

type Bitmap struct {
	Name   string
	Size   Size
	Pixels []byte
}

type Bitmaps map[string]*Bitmap

func (c Color) Byte() []byte {
	return []byte{c.R, c.G, c.B, c.A}
}

func (palette Palette) GradientIndex(color []byte) int {
	for i, c := range palette {
		if color[0] == c.R && color[1] == c.G && color[2] == c.B {
			return i
		}
	}

	return -1
}
