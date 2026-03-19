package manager

import (
	"time"

	"apodeiktikos.com/fbtest/engine"
	"apodeiktikos.com/fbtest/interfaces"
	"github.com/ukoelguanche/graphicsengine/core"
	"github.com/ukoelguanche/graphicsengine/drivers"
)

type stressMotion struct {
	targetA  core.Point
	targetB  core.Point
	duration time.Duration
	toB      bool
}

var stressMotions = map[*engine.SpriteInstance]*stressMotion{}

func SetupStressSprites(sprites core.Sprites, renderables []interfaces.Renderable, count int) []interfaces.Renderable {
	if count <= 0 {
		return renderables
	}

	stressMotions = make(map[*engine.SpriteInstance]*stressMotion, count)

	for i := 0; i < count; i++ {
		spriteName, sequenceName := stressSpriteSpec(i)
		position := stressSpritePosition(i)
		sprite := engine.BuildSpriteInstance(sprites, spriteName, sequenceName, position)

		if i%3 == 0 {
			attachStressMotion(sprite, i)
		}

		renderables = append(renderables, sprite)
	}

	return renderables
}

func stressSpriteSpec(index int) (string, string) {
	switch index % 4 {
	case 0:
		return "Ring", "idle"
	case 1:
		return "Flower1", "idle"
	case 2:
		return "Flower2", "idle"
	default:
		return "Sonic", "idle"
	}
}

func stressSpritePosition(index int) core.Point {
	const columns = 8
	cellW := float64(drivers.VW / columns)
	cellH := 22.0

	col := index % columns
	row := index / columns

	x := float64(col)*cellW + float64((index%3)*4)
	y := 8.0 + float64(row%8)*cellH

	return core.Point{X: x, Y: y}
}

func attachStressMotion(sprite *engine.SpriteInstance, index int) {
	start := sprite.GetPosition()
	target := start

	if index%2 == 0 {
		target = target.IncY(10)
	} else {
		target = target.IncX(12)
	}

	motion := &stressMotion{
		targetA:  start,
		targetB:  target,
		duration: time.Duration(350+(index%5)*80) * time.Millisecond,
		toB:      true,
	}

	stressMotions[sprite] = motion
	sprite.SetEaseFunction(engine.EaseInOutQuad)
	sprite.SetOnMovementComplete(onStressMotionComplete)
	sprite.MoveTo(motion.targetB, motion.duration)
}

func onStressMotionComplete(movable interfaces.Movable) {
	sprite, ok := movable.(*engine.SpriteInstance)
	if !ok {
		return
	}

	motion, ok := stressMotions[sprite]
	if !ok {
		return
	}

	nextTarget := motion.targetA
	if !motion.toB {
		nextTarget = motion.targetB
	}

	motion.toB = !motion.toB
	sprite.MoveTo(nextTarget, motion.duration)
}
