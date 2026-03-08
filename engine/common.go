package engine

import (
	"apodeiktikos.com/fbtest/interfaces"
	"github.com/ukoelguanche/graphicsengine/drivers"
)

func RenderEntity(r interfaces.Renderable) {
	r.Draw(drivers.GlobalDisplay)
}
