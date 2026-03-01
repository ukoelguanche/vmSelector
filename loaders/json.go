package loaders

import (
	"encoding/json"
	"os"

	"apodeiktikos.com/fbtest/model"
)

func LoadJSON(path string, destination *model.SpriteDefinition) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &destination)
}
