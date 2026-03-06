package manager

import (
	"apodeiktikos.com/fbtest/core"
	"apodeiktikos.com/fbtest/engine"
	"apodeiktikos.com/fbtest/util"
)

func SetupRing(sprites core.Sprites) *engine.SpriteInstance {
	ring = engine.BuildSpriteInstance(sprites, "Ring", "idle", core.Point{X: hudOffset + 20, Y: 56})
	ring.SetEaseFunction(util.EaseInOutCubic)
	ring.OnAnimationComplete = OnRingAnimationComplete

	return ring
}
