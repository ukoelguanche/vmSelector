package manager

import (
	"apodeiktikos.com/fbtest/engine"
	"github.com/ukoelguanche/graphicsengine/core"
)

func SetupRing(sprites core.Sprites) *engine.SpriteInstance {
	ring = engine.BuildSpriteInstance(sprites, "Ring", "idle", core.Point{X: hudOffset + 20, Y: 56})
	ring.SetEaseFunction(engine.EaseInOutCubic)
	ring.OnAnimationComplete = OnRingAnimationComplete

	return ring
}
