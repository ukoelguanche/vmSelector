package loaders

import (
	//"fmt"

	"log"

	"apodeiktikos.com/fbtest/model"
)

func LoadSprites(definitionPath string, sprites *model.Sprites) {
	LoadJson(definitionPath, sprites)

	// Load all bitmaps
	bitmaps := make(model.Bitmaps)
	for name, path := range sprites.BitmapSources {
		log.Printf("Loading bitmap source: %s %s", name, path)
		bitmaps[name] = LoadBitmap(path)
		bitmaps[name].Name = name
	}

	// Assign bitmap pointers to sprites
	for name, sprite := range sprites.Sprites {
		sprite.Name = name
		sprite.Bitmap = bitmaps[sprite.BitmapSource]
		log.Printf("Bitmap [%s] assigned to sprite: [%s]", sprite.BitmapSource, sprite.Name)
	}
	return

}
