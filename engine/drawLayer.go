package engine

import (
	"apodeiktikos.com/fbtest/interfaces"
	"github.com/ukoelguanche/graphicsengine/core"
)

type DrawLayer struct {
	BaseMovable
	draws    []CachedSpriteDraw
	isStatic bool
	easeFunc func(float64) float64
}

func (l *DrawLayer) GetSprite() *core.Sprite                 { return nil }
func (l *DrawLayer) GetEaseFunction() func(float64) float64  { return l.easeFunc }
func (l *DrawLayer) SetEaseFunction(f func(float64) float64) { l.easeFunc = f }
func (l *DrawLayer) Update()                                 {}
func (l *DrawLayer) IsStatic() bool                          { return l.isStatic }
func (l *DrawLayer) Draw(d interfaces.Drawer) {
	for _, draw := range l.draws {
		d.DrawSpriteRect(draw.Sprite, draw.Frame, draw.Position)
	}
}

func BuildDrawLayer(draws []CachedSpriteDraw, isStatic bool) *DrawLayer {
	return &DrawLayer{
		BaseMovable: BaseMovable{
			position:       core.Point{X: 0, Y: 0},
			startPosition:  core.Point{X: 0, Y: 0},
			targetPosition: core.Point{X: 0, Y: 0},
			moving:         false,
		},
		draws:    draws,
		isStatic: isStatic,
	}
}
