package render

import (
	"apodeiktikos.com/fbtest/drivers"
	"apodeiktikos.com/fbtest/engine"
)

func RenderEntity(r engine.Renderable) {
	r.Draw(drivers.GlobalDisplay)
}
