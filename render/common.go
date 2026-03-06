package render

import (
	"apodeiktikos.com/fbtest/drivers"
	"apodeiktikos.com/fbtest/engine"
)

func RenderEntity(r engine.Renderable) {
	r.Draw(drivers.GlobalDisplay)
	/*
		pos := r.GetPosition()
		img := r.GetTexture() // Supongamos que Renderable tiene este método

		// El Renderer decide CÓMO usar el driver para pintar el modelo
		drivers.FB.Draw(img, int(pos.X), int(pos.Y))
	*/
}
