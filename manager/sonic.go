package manager

import (
	"math/rand"
	"time"

	"apodeiktikos.com/fbtest/model"
	"apodeiktikos.com/fbtest/util"
)

var sonic *model.SpriteInstance
var lastIdleTime time.Time = time.Now()
var boredInterval time.Duration = 10
var boredAnimations = []string{"stare", "RaiseEyebrows", "FootTap"}

func SetupSonic(sprites model.Sprites) *model.SpriteInstance {
	sonic = model.BuildSpriteInstance(sprites, "Sonic", "idle", model.Point{X: 39, Y: 132})
	sonic.OnAnimationComplete = SonicIddleComplete

	return sonic
}

func SonicIddleComplete(sprite *model.SpriteInstance) {
	if time.Since(lastIdleTime) < boredInterval*time.Second {
		return
	}

	selectedAnim := boredAnimations[rand.Intn(len(boredAnimations))]

	sonic.CurrentSequence = sonic.Sprite.Sequences[selectedAnim]
	sonic.OnAnimationComplete = SonicBoredCompleted

}

func SonicBoredCompleted(sprite *model.SpriteInstance) {
	lastIdleTime = time.Now()
	sonic.CurrentSequence = sonic.Sprite.Sequences["idle"]
	sonic.OnAnimationComplete = SonicIddleComplete
}

func SonicStartJump() {
	sonic.OnAnimationComplete = nil
	sonic.CurrentSequence = sonic.Sprite.Sequences["jump"]
	sonic.SetEaseFunction(util.EaseOutQuad)
	sonic.OnMovementComplete = SonicJump1
	sonic.MoveTo(sonic.GetPosition().IncY(-80), 400*time.Millisecond)
}

func SonicJump1(sprite model.Renderable) {
	sonic.SetOnMovementComplete(SonicJump2)
	sonic.MoveTo(sonic.GetPosition().IncY(80), 200*time.Millisecond)

}

func SonicJump2(sprite model.Renderable) {
	sonic.SetOnMovementComplete(nil)
	sonic.CurrentSequence = sonic.Sprite.Sequences["idle"]
}
