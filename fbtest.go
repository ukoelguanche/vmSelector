package main

import (
	_ "image/png"
	"time"

	"apodeiktikos.com/fbtest/core"
	"apodeiktikos.com/fbtest/drivers"
	"apodeiktikos.com/fbtest/engine"
	"apodeiktikos.com/fbtest/loaders"
	"apodeiktikos.com/fbtest/manager"
	"apodeiktikos.com/fbtest/render"
	"apodeiktikos.com/fbtest/util"
)

const targetFPS = 25
const frameDelay = time.Second / targetFPS

var renderables []engine.Renderable

func Init() {
	util.LoadContext()
	drivers.GlobalDisplay = drivers.InitDisplay(drivers.VW, drivers.VH)
	drivers.GlobalKeyboard = drivers.InitKeyboard()

	var sprites core.Sprites
	loaders.LoadSprites("./resources/sprites/Sprites.json", &sprites)

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
		if elapsed < frameDelay {
			time.Sleep(frameDelay - elapsed)
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
		render.RenderEntity(renderable)
		renderable.NextFrame()
	}
	drivers.GlobalDisplay.Present()
}
