package manager

import (
	"math/rand"
	"time"

	"apodeiktikos.com/fbtest/core"
	"apodeiktikos.com/fbtest/engine"
	"apodeiktikos.com/fbtest/util"
)

var sonic *engine.SpriteInstance
var lastIdleTime time.Time = time.Now()
var boredInterval time.Duration = 3
var boredAnimations = []string{"stare", "RaiseEyebrows", "FootTap"}
var jumping = false
var jumpHeight float64 = 70
var jumpDuration = 350 * time.Millisecond

func SetupSonic(sprites core.Sprites) *engine.SpriteInstance {
	sonic = engine.BuildSpriteInstance(sprites, "Sonic", "idle", core.Point{X: 39, Y: 132})
	sonic.OnAnimationComplete = SonicIdleComplete

	return sonic
}

func SonicIdleComplete(*engine.SpriteInstance) {
	if time.Since(lastIdleTime) < boredInterval*time.Second {
		return
	}

	selectedAnim := boredAnimations[rand.Intn(len(boredAnimations))]

	sonic.CurrentSequence = sonic.Sprite.Sequences[selectedAnim]
	sonic.OnAnimationComplete = SonicBoredCompleted

}

func SonicBoredCompleted(*engine.SpriteInstance) {
	lastIdleTime = time.Now()
	sonic.CurrentSequence = sonic.Sprite.Sequences["idle"]
	sonic.OnAnimationComplete = SonicIdleComplete
}

func SonicStartJump() {
	if jumping {
		return
	}
	jumping = true
	sonic.OnAnimationComplete = nil
	sonic.CurrentSequence = sonic.Sprite.Sequences["jump"]
	sonic.SetEaseFunction(util.EaseOutQuad)
	sonic.OnMovementComplete = SonicJump1
	sonic.MoveTo(sonic.GetPosition().IncY(-jumpHeight), jumpDuration)
}

func SonicJump1(engine.Renderable) {
	sonic.SetOnMovementComplete(SonicJump2)
	sonic.SetEaseFunction(util.EaseInQuad)
	sonic.MoveTo(sonic.GetPosition().IncY(jumpHeight), jumpDuration)

}

func SonicJump2(engine.Renderable) {
	// End jump
	jumping = false
	sonic.SetOnMovementComplete(nil)
	// Back to idle
	lastIdleTime = time.Now()
	sonic.CurrentSequence = sonic.Sprite.Sequences["idle"]
	sonic.OnAnimationComplete = SonicIdleComplete
}
