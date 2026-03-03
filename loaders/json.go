package loaders

import (
	"encoding/json"
	"log"
	"os"

	"apodeiktikos.com/fbtest/model"
)

func LoadJson(path string, sprites *model.Sprites) error {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Error loading sprite definition %s %v", path, err)
	}
	return json.Unmarshal(data, &sprites)
}
