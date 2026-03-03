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

type Sprite struct {
	Name         string
	BitmapSource string `json:"BitmapSource"`
	Bitmap       *Bitmap
	Frames       []Rect           `json:"Frames"`
	Sequences    map[string][]int `json:"Sequences"`
	Characters   map[string]int   `json:"Characters"`
}

type Sprites struct {
	BitmapSources map[string]string  `json:"BitmapSources"`
	Sprites       map[string]*Sprite `json:"sprites"`
}

type SpriteInstance struct {
	Sprite                  *Sprite
	Position                Point
	FrameIdx                int
	CurrentSequence         []int
	CurrentSequenceOffset   float32
	CurrentSequencePosition float32
	SequenceLength          int
	Scale                   float64
}

type Text struct {
	Sprite   *Sprite
	Position Point
	Text     string
}

func BuildSpriteInstance(sprites Sprites, name string, sequenceName string, position Point) *SpriteInstance {
	sequence := sprites.Sprites[name].Sequences[sequenceName]
	relativeSpeed := float32(0.5)
	spriteInstance := &SpriteInstance{
		Sprite:                  sprites.Sprites[name],
		Position:                position,
		FrameIdx:                0,
		CurrentSequence:         sequence,
		CurrentSequenceOffset:   1 / float32(len(sequence)) * relativeSpeed,
		CurrentSequencePosition: 0.0,
		SequenceLength:          len(sequence),
	}

	return spriteInstance

}

func (s *SpriteInstance) NextFrame() {
	s.CurrentSequencePosition += s.CurrentSequenceOffset
	if s.CurrentSequencePosition >= 1 {
		s.CurrentSequencePosition = 0
	}
	//s.FrameIdx = (s.FrameIdx + 1) % s.SequenceLength
}

func (s *SpriteInstance) CurrentFrame() Rect {
	frame := int(float32(len(s.CurrentSequence)) * s.CurrentSequencePosition)
	return s.Sprite.Frames[s.CurrentSequence[frame]]
}

// ToDo: Change W, H to Size type
type Bitmap struct {
	Name   string
	Size   Size
	Pixels []byte
}

type Bitmaps map[string]*Bitmap

type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

func (c Color) Byte() []byte {
	return []byte{c.R, c.G, c.B, c.A}
}

type Gradient []Color

/*

func (g Gradient) GradientIndex(color []byte) int {
	for i, c := range g {
		if color[0] == c.R && color[1] == c.G && color[2] == c.B {
			return i
		}
	}

	return -1
}

func (s Sprite) GetSection(sectionName string) SpriteDataSection {
	rects, ok := s.Sections[sectionName]
	if !ok {
		log.Fatalf("Section %s not found", sectionName)
	}

	return rects
}

func (s Sprite) GetAnimation(animationName string) SpriteAnimation {
	animationFrame, ok := s.Animations[animationName]
	if !ok {
		log.Fatalf("Animation %s not found", animationName)
	}

	return animationFrame
}

func (s Sprite) GetAnimationRects(animationSectionName string) []Rect {
	animationRects, ok := s.AnimationSections[animationSectionName]
	if !ok {
		log.Fatalf("Animation rect %s not found", animationSectionName)
	}

	return animationRects
}

func (section SpriteDataSection) GetSprite(name string) Rect {
	rect, ok := section[name]
	if !ok {
		log.Fatalf("Sprite %s not found", name)
	}
	return rect
}
*/
