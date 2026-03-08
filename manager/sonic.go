package manager

import (
	"math/rand"
	"time"

	"apodeiktikos.com/fbtest/engine"
	"apodeiktikos.com/fbtest/interfaces"
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

func SonicIdleComplete(animatable interfaces.Animatable) {
	if time.Since(lastIdleTime) < boredInterval*time.Second {
		return
	}

	selectedAnim := boredAnimations[rand.Intn(len(boredAnimations))]

	engine.SetCurrentSequenceByName(sonic, selectedAnim)
	animatable.SetOnAnimationComplete(SonicBoredCompleted)
}

func SonicBoredCompleted(sonic interfaces.Animatable) {
	lastIdleTime = time.Now()
	engine.SetCurrentSequenceByName(sonic, "idle")
	sonic.SetOnAnimationComplete(SonicIdleComplete)
}

func SonicStartJump() {
	if jumping {
		return
	}
	jumping = true
	sonic.SetOnAnimationComplete(nil)
	engine.SetCurrentSequenceByName(sonic, "jump")
	sonic.SetEaseFunction(engine.EaseOutQuad)
	sonic.SetOnMovementComplete(SonicJump1)
	sonic.MoveTo(sonic.GetPosition().IncY(-jumpHeight), jumpDuration)
}

func SonicJump1(movable interfaces.Movable) {
	sonic.SetOnMovementComplete(SonicJump2)
	sonic.SetEaseFunction(engine.EaseInQuad)
	sonic.MoveTo(sonic.GetPosition().IncY(jumpHeight), jumpDuration)

}

func SonicJump2(saa interfaces.Movable) {
	// End jump
	jumping = false
	sonic.SetOnMovementComplete(nil)
	// Back to idle
	lastIdleTime = time.Now()
	engine.SetCurrentSequenceByName(sonic, "idle")
	sonic.SetOnAnimationComplete(SonicIdleComplete)
}
