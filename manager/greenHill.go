package manager

import (
	"fmt"
	"time"

	"apodeiktikos.com/fbtest/engine"
	"apodeiktikos.com/fbtest/interfaces"
	"github.com/ukoelguanche/graphicsengine/core"
)

var clouds []*engine.SpriteInstance

func SetupGreenHillBackground(sprites core.Sprites, renderables []interfaces.Renderable) []interfaces.Renderable {
	renderables = append(renderables, engine.BuildSpriteInstance(sprites, "GreenHillBackgroundLayer4", "idle", core.Point{X: 0, Y: 64}))
	renderables = append(renderables, engine.BuildSpriteInstance(sprites, "GreenHillBackgroundLayer5", "idle", core.Point{X: 0, Y: 112}))
	renderables = append(renderables, engine.BuildSpriteInstance(sprites, "GreenHillBackgroundLayer6", "idle", core.Point{X: 0, Y: 152}))

	return renderables
}

func SetupClouds(sprites core.Sprites, renderables []interfaces.Renderable) []interfaces.Renderable {
	var cloudSprite *engine.SpriteInstance
	clouds = make([]*engine.SpriteInstance, 0)
	var y float64 = 0
	for i := 0; i < 3; i++ {
		cloudSprite = engine.BuildSpriteInstance(sprites, fmt.Sprintf("GreenHillBackgroundLayer%d", i+1), "idle", core.Point{X: 0, Y: float64(y)})

		cloudSprite.MoveTo(cloudSprite.GetPosition().SetX(-980), time.Duration(90000+i*10000)*time.Millisecond)

		cloudSprite.SetOnMovementComplete(OnCloudMovementComplete)
		renderables = append(renderables, cloudSprite)
		y += cloudSprite.GetFrame(0).Size.H

		clouds = append(clouds, cloudSprite)
	}

	return renderables
}

func SetupGreenHillForeground(sprites core.Sprites, renderables []interfaces.Renderable) []interfaces.Renderable {
	renderables = append(renderables, engine.BuildSpriteInstance(sprites, "Flower1", "idle", core.Point{X: 154, Y: 90}))
	renderables = append(renderables, engine.BuildSpriteInstance(sprites, "Flower2", "idle", core.Point{X: -5, Y: 115}))
	renderables = append(renderables, engine.BuildSpriteInstance(sprites, "Flower2", "idle", core.Point{X: 220, Y: 115}))
	renderables = append(renderables, engine.BuildSpriteInstance(sprites, "Flower2", "idle", core.Point{X: 250, Y: 115}))
	renderables = append(renderables, engine.BuildSpriteInstance(sprites, "GreenHillForeground", "idle", core.Point{X: 0, Y: 0}))

	return renderables
}

func OnCloudMovementComplete(sprite interfaces.Renderable) {
	spritePosition := sprite.GetPosition()
	sprite.SetPosition(spritePosition.SetX(0))
	sprite.SetTargetPosition(spritePosition.SetX(-980))

}
