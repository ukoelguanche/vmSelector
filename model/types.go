package model

import (
	"log"
)

type Point struct {
	X, Y int
}

type Size struct {
	W, H int
}

type Rect struct {
	Point Point
	Size  Size
}

type SpriteDataSection map[string]Rect
type SpriteAnimation []Rect
type SpriteDefinition struct {
	SourceImage string
	Sections    map[string]SpriteDataSection
}

// ToDo: Change W, H to Size type
type Bitmap struct {
	W, H   int
	Pixels []byte
}

type Sprite struct {
	SourceImage string
	Sections    map[string]SpriteDataSection
	Animations  map[string]SpriteAnimation
	Bitmap      *Bitmap
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
		log.Fatalf("Section %s not found", animationName)
	}

	return animationFrame
}

func (section SpriteDataSection) GetSprite(name string) Rect {
	rect, ok := section[name]
	if !ok {
		log.Fatalf("Sprite %s not found", name)
	}
	return rect
}
