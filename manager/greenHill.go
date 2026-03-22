package manager

import (
	"fmt"
	"time"

	"apodeiktikos.com/fbtest/engine"
	"apodeiktikos.com/fbtest/interfaces"
	"apodeiktikos.com/fbtest/util"
	"github.com/ukoelguanche/graphicsengine/core"
)

var clouds []*engine.SpriteInstance

func SetupGreenHillBackground(sprites core.Sprites, renderables []interfaces.Renderable) []interfaces.Renderable {
	renderables = append(renderables, buildGreenHillLayer("GreenHillStaticBackground", []engine.CachedSpriteDraw{
		{
			Sprite:   sprites.Sprites["GreenHillBackgroundLayer4"],
			Frame:    sprites.Sprites["GreenHillBackgroundLayer4"].GetFrame(0),
			Position: core.Point{X: 0, Y: 64},
		},
	}, true))
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
		y += cloudSprite.GetSprite().GetFrame(0).Size.H

		clouds = append(clouds, cloudSprite)
	}

	return renderables
}

func SetupGreenHillForeground(sprites core.Sprites, renderables []interfaces.Renderable) []interfaces.Renderable {
	renderables = append(renderables, buildGreenHillLayer("GreenHillStaticForeground", []engine.CachedSpriteDraw{
		{
			Sprite:   sprites.Sprites["GreenHillForeground"],
			Frame:    sprites.Sprites["GreenHillForeground"].GetFrame(0),
			Position: core.Point{X: 0, Y: 0},
		},
	}, false))

	renderables = append(renderables, engine.BuildSpriteInstance(sprites, "Flower1", "idle", core.Point{X: 154, Y: 90}))
	renderables = append(renderables, engine.BuildSpriteInstance(sprites, "Flower2", "idle", core.Point{X: -5, Y: 115}))
	renderables = append(renderables, engine.BuildSpriteInstance(sprites, "Flower2", "idle", core.Point{X: 220, Y: 115}))
	renderables = append(renderables, engine.BuildSpriteInstance(sprites, "Flower2", "idle", core.Point{X: 250, Y: 115}))

	return renderables
}

func buildGreenHillLayer(name string, draws []engine.CachedSpriteDraw, isStatic bool) interfaces.Renderable {
	if util.ContextStorage.UseCachedLayers {
		return engine.BuildCachedLayer(name, draws, isStatic)
	}
	return engine.BuildDrawLayer(draws, isStatic)
}

func OnCloudMovementComplete(sprite interfaces.Movable) {
	cloud, ok := sprite.(*engine.SpriteInstance)
	if !ok {
		return
	}

	startPosition := cloud.GetPosition().SetX(0)
	cloud.SetPosition(startPosition)
	cloud.MoveTo(startPosition.SetX(-980), cloud.GetDuration())
}
