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

type SpriteAnimation struct {
	Section string
	Frames  []int
}

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
	Bitmap            *Bitmap
	SourceImage       string
	Sections          map[string]SpriteDataSection
	AnimationSections map[string][]Rect
	Animations        map[string]SpriteAnimation
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
