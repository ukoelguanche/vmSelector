package main

import (
	_ "image/png"
	"time"

	"apodeiktikos.com/fbtest/engine"
	"apodeiktikos.com/fbtest/interfaces"
	"apodeiktikos.com/fbtest/manager"
	"apodeiktikos.com/fbtest/util"
	"github.com/ukoelguanche/graphicsengine/core"
	"github.com/ukoelguanche/graphicsengine/drivers"
	"github.com/ukoelguanche/graphicsengine/loaders"
)

const TARGET_FPS = 25
const FRAME_DELAY = time.Second / TARGET_FPS
const SPRITES_FILE = "./assets/sprites/Sprites.json"

var renderables []interfaces.Renderable

func Init() {
	util.LoadContext()
	drivers.GlobalDisplay = drivers.InitDisplay("VM Selector", drivers.VW, drivers.VH)
	drivers.InitKeyboard()

	var sprites core.Sprites
	loaders.LoadSprites(SPRITES_FILE, &sprites)

	renderables = make([]interfaces.Renderable, 0)

	renderables = manager.SetupClouds(sprites, renderables)
	renderables = manager.SetupGreenHillBackground(sprites, renderables)
	renderables = manager.SetupGreenHillForeground(sprites, renderables)

	renderables = append(renderables, manager.SetupSonic(sprites))
	renderables = manager.SetupHud(sprites, renderables)

	drivers.SpriteColorProcessors = append(drivers.SpriteColorProcessors, &engine.PaletteSwapColorProcessor{
		SourcePalette: sprites.Palettes["WaterSource"],
		TargetPalette: sprites.Palettes["WaterTarget"],
		Sprite:        renderables[4].GetSprite(),
	})
}

func main() {
	Init()

	defer drivers.GlobalDisplay.Close()
	defer drivers.GlobalDisplay.Clear()
	defer drivers.GlobalKeyboard.Close()

	for Loop() {
	}
}

func Loop() bool {
	start := time.Now()

	drivers.GlobalDisplay.Clear()
	for _, renderable := range renderables {
		engine.RenderEntity(renderable)
		renderable.NextFrame()
	}

	elapsed := time.Since(start)
	//log.Println("Elapsed time: ", elapsed)
	if elapsed < FRAME_DELAY {
		time.Sleep(FRAME_DELAY - elapsed)
	}

	drivers.GlobalDisplay.Present()
	return handleKeyboardInput()
}

func handleKeyboardInput() bool {
	kbd := drivers.GlobalKeyboard.GetInput()
	var inc int

	if kbd == drivers.KBD_ESCAPE {
		return false
	} else if kbd == drivers.KBD_RETURN {
		manager.SelectMenuOption()
	} else if kbd == drivers.KBD_SPACE {
		manager.SonicStartJump()
	} else if kbd == drivers.KBD_UP {
		inc = -1
	} else if kbd == drivers.KBD_DOWN {
		inc = 1
	} else if kbd == drivers.KBD_LEFT {
	} else if kbd == drivers.KBD_RIGHT {
	}

	manager.IncrementVMIndex(inc)

	return true
}
