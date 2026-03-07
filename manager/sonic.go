package manager

import (
	"math/rand"
	"time"

	"apodeiktikos.com/fbtest/engine"
	"github.com/ukoelguanche/graphicsengine/core"
)

var sonic *engine.Character
var lastIdleTime time.Time = time.Now()
var boredInterval time.Duration = 3
var boredAnimations = []string{"stare", "RaiseEyebrows", "FootTap"}
var jumping = false
var jumpHeight float64 = 70
var jumpDuration = 350 * time.Millisecond

func SetupSonic(sprites core.Sprites) *engine.Character {
	sonic = engine.BuildCharacter(sprites, "Sonic", "idle", core.Point{X: 39, Y: 132})
	sonic.SetOnAnimationComplete(SonicIdleComplete)

	return sonic
}

func SonicIdleComplete(sonic *engine.Character) {
	if time.Since(lastIdleTime) < boredInterval*time.Second {
		return
	}

	selectedAnim := boredAnimations[rand.Intn(len(boredAnimations))]

	sonic.CurrentSequence = sonic.Sprite.Sequences[selectedAnim]
	sonic.SetOnAnimationComplete(SonicBoredCompleted)

}

func SonicBoredCompleted(sonic *engine.Character) {
	lastIdleTime = time.Now()
	sonic.CurrentSequence = sonic.Sprite.Sequences["idle"]
	sonic.SetOnAnimationComplete(SonicIdleComplete)
}

func SonicStartJump() {
	if jumping {
		return
	}
	jumping = true
	sonic.OnAnimationComplete = nil
	sonic.CurrentSequence = sonic.Sprite.Sequences["jump"]
	sonic.SetEaseFunction(engine.EaseOutQuad)
	sonic.OnMovementComplete = SonicJump1
	sonic.MoveTo(sonic.GetPosition().IncY(-jumpHeight), jumpDuration)
}

func SonicJump1(engine.Renderable) {
	sonic.SetOnMovementComplete(SonicJump2)
	sonic.SetEaseFunction(engine.EaseInQuad)
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
