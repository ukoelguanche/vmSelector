package engine

import (
	"github.com/ukoelguanche/graphicsengine/drivers"
)

func RenderEntity(r Renderable) {
	r.Draw(drivers.GlobalDisplay)
}
