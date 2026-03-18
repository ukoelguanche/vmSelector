package manager

import "apodeiktikos.com/fbtest/engine"

type KeyboardEventsHandler struct {
}

func (KeyboardEventsHandler) OnKeyboardLost() {
	engine.SetCurrentSequenceByName(keyboardIcon, "blink")
}

func (KeyboardEventsHandler) OnKeyboardPlugged() {
	engine.SetCurrentSequenceByName(keyboardIcon, "idle")
}
