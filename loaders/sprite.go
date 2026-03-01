package loaders

import (
	"apodeiktikos.com/fbtest/model"
)

func LoadSprite(definitionPath string) *model.Sprite {
	sprite := model.Sprite{}

	LoadJson(definitionPath, &sprite)
	sprite.Bitmap = LoadBitmap(sprite.SourceImage)

	return &sprite
}
