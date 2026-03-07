package main

import (
	_ "image/png"
	"time"

	"apodeiktikos.com/fbtest/engine"
	"apodeiktikos.com/fbtest/manager"
	"github.com/ukoelguanche/graphicsengine/core"
	"github.com/ukoelguanche/graphicsengine/drivers"
	"github.com/ukoelguanche/graphicsengine/loaders"
	//"apodeiktikos.com/fbtest/render"
	"apodeiktikos.com/fbtest/util"
)

const TARGET_FPS = 25
const FRAME_DELAY = time.Second / TARGET_FPS
const SPRITES_FILE = "./assets/sprites/Sprites.json"

var renderables []engine.Renderable

func Init() {
	util.LoadContext()
	drivers.GlobalDisplay = drivers.InitDisplay(drivers.VW, drivers.VH)
	drivers.GlobalKeyboard = drivers.InitKeyboard()

	var sprites core.Sprites
	loaders.LoadSprites(SPRITES_FILE, &sprites)

	renderables = make([]engine.Renderable, 0)

	renderables = manager.SetupClouds(sprites, renderables)
	renderables = manager.SetupGreenHillBackground(sprites, renderables)
	renderables = manager.SetupGreenHillForeground(sprites, renderables)

	renderables = append(renderables, manager.SetupSonic(sprites))
	renderables = manager.SetupHud(sprites, renderables)
}

func main() {
	Init()

	defer drivers.GlobalDisplay.Close()
	defer drivers.GlobalDisplay.Clear()
	defer drivers.GlobalKeyboard.Close()

	for {
		start := time.Now()

		if handleKeyboardInput() {
			break
		}

		Loop()

		elapsed := time.Since(start)
		//log.Println("Elapsed time: ", elapsed)
		if elapsed < FRAME_DELAY {
			time.Sleep(FRAME_DELAY - elapsed)
		}
	}
}

func handleKeyboardInput() bool {
	kbd := drivers.GlobalKeyboard.GetInput()
	var inc int

	if kbd == drivers.KBD_ESCAPE {
		return true
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

	return false
}

func Loop() {
	drivers.GlobalDisplay.Clear()
	for _, renderable := range renderables {
		engine.RenderEntity(renderable)
		renderable.NextFrame()
	}
	drivers.GlobalDisplay.Present()
}
