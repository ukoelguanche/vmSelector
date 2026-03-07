package manager

import (
	"time"

	"apodeiktikos.com/fbtest/engine"
	"apodeiktikos.com/fbtest/interfaces"
	"apodeiktikos.com/fbtest/model"
	"apodeiktikos.com/fbtest/util"
	"github.com/ukoelguanche/graphicsengine/core"
)

const hudOffset float64 = 175

var ring *engine.SpriteInstance
var gpuString string
var centinelVM *model.VM
var vms []model.VM
var selectedVMIndex = 0
var texts []*engine.Text
var menuStatus = "waiting"

func SetupHud(sprites core.Sprites, renderables []interfaces.Renderable) []interfaces.Renderable {
	for y := 0; y < 13; y++ {
		renderables = append(renderables, engine.BuildSpriteInstance(sprites, "ZigZag", "idle", core.Point{X: hudOffset, Y: float64(y * 16)}))
	}

	texts = SetupHUDTexts(sprites, hudOffset)
	for _, t := range texts {
		renderables = append(renderables, t)
	}

	ring = engine.BuildSpriteInstance(sprites, "Ring", "idle", core.Point{X: hudOffset + 20, Y: 56})
	renderables = append(renderables, ring)

	return renderables
}

func SetupHUDTexts(sprites core.Sprites, hudOffset float64) []*engine.Text {
	texts := make([]*engine.Text, 0)
	gpuString = util.ContextStorage.GpuString
	centinelVM = model.GetVMByName(util.ContextStorage.CentineVMName)
	vms = model.GetVMsWithGPU(gpuString, centinelVM)
	for i, vm := range vms {
		var text string
		if vm.Equals(centinelVM) {
			text = "POWER OFF"
		} else {
			text = vm.Name
		}

		textInstance := engine.BuildTextInstance(sprites.Sprites["BoldLetters"], text, core.Point{X: hudOffset + 30, Y: float64(i)*16 + 60})
		texts = append(texts, textInstance)
	}

	// ToDo uncomment next lines after rebugging
	texts[0].Position.X += 12
	texts[len(texts)-1].Position.Y += 12

	texts = append(texts, engine.BuildTextInstance(sprites.Sprites["GenesisLetters"], centinelVM.Name, core.Point{X: hudOffset + 20, Y: 30}))

	return texts
}

func IncrementVMIndex(value int) {
	if value == 0 || selectedVMIndex >= len(vms) || selectedVMIndex < 0 {
		return
	}

	if ring.IsMoving() {
		return
	}

	const transitionDuaration = 200 * time.Millisecond
	texts[selectedVMIndex].MoveTo(texts[selectedVMIndex].Position.SetX(hudOffset+30), transitionDuaration)
	selectedVMIndex = max(0, min(len(vms)-1, selectedVMIndex+value))
	texts[selectedVMIndex].MoveTo(texts[selectedVMIndex].Position.SetX(hudOffset+42), transitionDuaration)

	ring.MoveTo(ring.GetPosition().SetY(texts[selectedVMIndex].Position.Y-4), transitionDuaration)
	ring.SetOnMovementComplete(OnMovementComplete)
}

func SelectMenuOption() {
	if menuStatus != "waiting" {
		return
	}
	menuStatus = "ending"

	ring.SetCurrentSequence(ring.GetSequences("fade"))
	ring.SetOnAnimationComplete(OnRingAnimationComplete)

	ci := 0
	for i, text := range texts[:len(texts)-1] {
		if i == selectedVMIndex || i == len(texts)-1 {
			continue
		}
		text.SetOnMovementComplete(OnMovementComplete)
		text.SetEaseFunction(engine.EaseInOutQuad)
		text.MoveTo(text.Position, time.Duration(ci*100)*time.Millisecond)
		ci++
	}

}

func OnMovementComplete(renderabe interfaces.Renderable) {
	for _, text := range texts {
		if renderabe == text {
			text.SetEaseFunction(engine.EaseInOutQuad)
			text.MoveTo(core.Point{X: 350, Y: text.Position.Y}, 1000*time.Millisecond)
		}
	}
}

func OnRingAnimationComplete(renderable interfaces.Renderable) {
	ring.SetCurrentSequence(ring.GetSequences("end"))
}
