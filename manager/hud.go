package manager

import (
	"time"

	"apodeiktikos.com/fbtest/engine"
	"apodeiktikos.com/fbtest/engine/colorProcessors"
	"apodeiktikos.com/fbtest/interfaces"
	"apodeiktikos.com/fbtest/model"
	"apodeiktikos.com/fbtest/util"
	"github.com/ukoelguanche/graphicsengine/core"
	"github.com/ukoelguanche/graphicsengine/drivers"
)

const hudOffset float64 = 95

var ring *engine.SpriteInstance
var keyboardIcon *engine.SpriteInstance
var gpuString string
var centinelVM *model.VM
var vms []model.VM
var selectedVMIndex = 0
var texts []*engine.Text
var menuStatus = "waiting"

var textHeight float64
var verticalCenter float64
var pixelFade *colorProcessors.PixelFade
var blackScreen *colorProcessors.BlackScreen

func SetupHud(sprites core.Sprites, renderables []interfaces.Renderable) []interfaces.Renderable {
	for y := 0; y < 13; y++ {
		renderables = append(renderables, engine.BuildSpriteInstance(sprites, "ZigZag", "idle", core.Point{X: hudOffset, Y: float64(y * 16)}))
	}

	textHeight = sprites.Sprites["GenesisLetters"].Frames[0].Size.H
	verticalCenter = drivers.VH/2 - textHeight/2

	texts = SetupHUDTexts(sprites, hudOffset)
	for _, t := range texts {
		renderables = append(renderables, t)
	}

	ring = engine.BuildSpriteInstance(sprites, "Ring", "idle", core.Point{X: hudOffset + 28, Y: verticalCenter})
	ring.SetEaseFunction(engine.EaseInQuad)
	renderables = append(renderables, ring)

	keyboardIcon = engine.BuildSpriteInstance(sprites, "Keyboard", "idle", core.Point{X: drivers.VW/2 - 30, Y: drivers.VH/2 - 23})
	renderables = append(renderables, keyboardIcon)

	drivers.GlobalKeyboard.EventsHandler = KeyboardEventsHandler{}

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

		textInstance := engine.BuildTextInstance(sprites.Sprites["GenesisLetters"], text, getTextPositionForVMIndex(i))
		textInstance.SetEaseFunction(engine.EaseInOutQuad)
		texts = append(texts, textInstance)
	}

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
	selectedVMIndex = max(0, min(len(vms)-1, selectedVMIndex+value))
	for i, _ := range vms {
		position := getTextPositionForVMIndex(i)
		texts[i].MoveTo(position, transitionDuaration)
	}
}

func getTextPositionForVMIndex(index int) core.Point {
	var textOffset float64 = 35
	if index == selectedVMIndex {
		textOffset = 50
	}
	return core.Point{X: hudOffset + textOffset, Y: verticalCenter + float64((index-selectedVMIndex)*25)}
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

func OnMovementComplete(renderabe interfaces.Movable) {
	for _, text := range texts {
		if renderabe == text {
			text.SetEaseFunction(engine.EaseInOutQuad)
			text.MoveTo(core.Point{X: 350, Y: text.Position.Y}, 1000*time.Millisecond)
		}
	}
}

func OnRingAnimationComplete(renderable interfaces.Animatable) {
	ring.SetOnAnimationComplete(VMSelected)
	ring.SetCurrentSequence(ring.GetSequences("end"))
	devMode := util.ContextStorage.DevMode
	if devMode == "true" {
		return
	}
	model.SwitchToVM(centinelVM, vms[selectedVMIndex])
}

func VMSelected(renderable interfaces.Animatable) {
	ring.SetOnAnimationComplete(nil)
	pixelFade = colorProcessors.BuildPixelFade(8, 700*time.Millisecond)
	//pixelFade = engine.BuildReversePixelFade(8, 700*time.Millisecond)
	drivers.GlobalDisplay.AddTransformer(pixelFade)
	// pixelFade.OnComplete = PixelFadeComplete
}

func PixelFadeComplete() {
	drivers.GlobalDisplay.RemoveTransformer(pixelFade)
	blackScreen = colorProcessors.BuildBlackScreen(700 * time.Millisecond)
	blackScreen.OnComplete = BlackScreenComplete
	drivers.GlobalDisplay.AddTransformer(blackScreen)
}

func BlackScreenComplete() {
	drivers.GlobalDisplay.RemoveTransformer(blackScreen)
}
