package model

type Point struct {
	X, Y float64
}

func (p1 Point) Equals(p2 Point) bool { return p1.X == p2.X && p1.Y == p2.Y }
func (p Point) SetX(x float64) Point  { return Point{X: x, Y: p.Y} }
func (p Point) SetY(y float64) Point  { return Point{X: p.X, Y: y} }

func (p Point) IncX(x float64) Point { return Point{X: p.X + x, Y: p.Y} }
func (p Point) IncY(y float64) Point { return Point{X: p.X, Y: p.Y + y} }

type Size struct {
	W, H float64
}

func (s Size) SetW(w float64) Size      { return Size{W: w, H: s.H} }
func (s Size) SetH(h float64) Size      { return Size{W: s.W, H: h} }
func (s Size) Mult(factor float64) Size { return Size{W: s.W * factor, H: s.H * factor} }

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

type Bitmap struct {
	Name   string
	W      int32
	H      int32
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
