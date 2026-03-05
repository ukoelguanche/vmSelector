package drivers

var GlobalKeyboard *Keyboard

type KeyboardInput int

const (
	KBD_NONE KeyboardInput = iota
	KBD_RETURN
	KBD_ESCAPE
	KBD_SPACE
	KBD_UP
	KBD_DOWN
	KBD_LEFT
	KBD_RIGHT
)
