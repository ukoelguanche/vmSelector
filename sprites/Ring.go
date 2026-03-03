package sprites

import (
	"apodeiktikos.com/fbtest/model"
)

type Ring struct {
	Sprite model.Sprite
	Point  model.Point
}

func (r *Ring) Render() {
	//	drivers.DrawSprite(hud, "items", "zigZagBG", 90+xOffset, int32(i*16))
}
