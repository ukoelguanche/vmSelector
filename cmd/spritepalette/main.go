package main

import (
	_ "image/png"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ukoelguanche/graphicsengine/core"
	"github.com/ukoelguanche/graphicsengine/loaders"
)

const (
	defaultSpritesFile = "./assets/sprites/Sprites.json"
	defaultOutputDir   = "./assets/indexed"
)

type paletteEntry struct {
	Index int    `json:"index"`
	R     uint8  `json:"r"`
	G     uint8  `json:"g"`
	B     uint8  `json:"b"`
	A     uint8  `json:"a"`
	Note  string `json:"note,omitempty"`
}

type indexedBitmapMetadata struct {
	Name        string `json:"name"`
	SourcePath  string `json:"sourcePath"`
	Width       int32  `json:"width"`
	Height      int32  `json:"height"`
	IndexedPath string `json:"indexedPath"`
}

type outputMetadata struct {
	SpritesFile string                  `json:"spritesFile"`
	PalettePath string                  `json:"palettePath"`
	PaletteType string                  `json:"paletteType"`
	ColorCount  int                     `json:"colorCount"`
	Bitmaps     []indexedBitmapMetadata `json:"bitmaps"`
}

func main() {
	spritesPath := flag.String("sprites", defaultSpritesFile, "path to Sprites.json")
	outputDir := flag.String("out", defaultOutputDir, "output directory for palette and indexed bitmap files")
	verbose := flag.Bool("verbose", false, "print loader logs")
	flag.Parse()

	if !*verbose {
		log.SetOutput(io.Discard)
	}

	var sprites core.Sprites
	loaders.LoadSprites(*spritesPath, &sprites)

	palette, colorToIndex, err := buildGlobalPalette(sprites)
	if err != nil {
		fmt.Fprintf(os.Stderr, "palette generation failed: %v\n", err)
		os.Exit(1)
	}

	metadata, err := writeIndexedAssets(*outputDir, *spritesPath, sprites, palette, colorToIndex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "writing indexed assets failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Sprites file: %s\n", *spritesPath)
	fmt.Printf("Output directory: %s\n", *outputDir)
	fmt.Printf("Global palette size: %d entries\n", len(palette))
	fmt.Printf("Indexed bitmaps written: %d\n", len(metadata.Bitmaps))
	fmt.Printf("Palette metadata: %s\n", metadata.PalettePath)
}

func buildGlobalPalette(sprites core.Sprites) ([]paletteEntry, map[uint32]byte, error) {
	colors := collectBitmapColors(sprites)
	if len(colors) > 255 {
		return nil, nil, fmt.Errorf("referenced sprites need %d opaque colors plus transparency, which exceeds 256 palette entries", len(colors))
	}

	sortedColors := make([]uint32, 0, len(colors))
	for color := range colors {
		sortedColors = append(sortedColors, color)
	}
	sort.Slice(sortedColors, func(i, j int) bool { return sortedColors[i] < sortedColors[j] })

	palette := make([]paletteEntry, 0, len(sortedColors)+1)
	colorToIndex := make(map[uint32]byte, len(sortedColors))

	palette = append(palette, paletteEntry{
		Index: 0,
		R:     0,
		G:     0,
		B:     0,
		A:     0,
		Note:  "binary transparency",
	})

	for i, color := range sortedColors {
		r, g, b, _ := unpackRGBA(color)
		index := byte(i + 1)
		palette = append(palette, paletteEntry{
			Index: int(index),
			R:     r,
			G:     g,
			B:     b,
			A:     255,
		})
		colorToIndex[color] = index
	}

	return palette, colorToIndex, nil
}

func collectBitmapColors(sprites core.Sprites) map[uint32]struct{} {
	colors := make(map[uint32]struct{})
	seenBitmaps := make(map[string]struct{})

	for _, sprite := range sprites.Sprites {
		if sprite == nil || sprite.Bitmap == nil {
			continue
		}

		if _, seen := seenBitmaps[sprite.Bitmap.Name]; seen {
			continue
		}
		seenBitmaps[sprite.Bitmap.Name] = struct{}{}

		bitmap := sprite.Bitmap

		for offset := 0; offset < len(bitmap.Pixels); offset += 4 {
			a := bitmap.Pixels[offset+3]
			if a < 128 {
				continue
			}

			color := packRGBA(
				bitmap.Pixels[offset],
				bitmap.Pixels[offset+1],
				bitmap.Pixels[offset+2],
				255,
			)
			colors[color] = struct{}{}
		}
	}

	return colors
}

func writeIndexedAssets(outputDir string, spritesPath string, sprites core.Sprites, palette []paletteEntry, colorToIndex map[uint32]byte) (*outputMetadata, error) {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, err
	}

	palettePath := filepath.Join(outputDir, "palette.pal")
	if err := writePaletteFile(palettePath, palette); err != nil {
		return nil, err
	}

	bitmapNames := make([]string, 0, len(sprites.BitmapSources))
	for name := range sprites.BitmapSources {
		bitmapNames = append(bitmapNames, name)
	}
	sort.Strings(bitmapNames)

	metadata := &outputMetadata{
		SpritesFile: spritesPath,
		PalettePath: palettePath,
		PaletteType: "rgb24",
		ColorCount:  len(palette),
		Bitmaps:     make([]indexedBitmapMetadata, 0, len(bitmapNames)),
	}

	for _, bitmapName := range bitmapNames {
		bitmap := findBitmapByName(sprites, bitmapName)
		if bitmap == nil {
			return nil, fmt.Errorf("bitmap %q is not referenced by any sprite", bitmapName)
		}

		indexedPixels, err := buildIndexedBitmap(bitmap, colorToIndex)
		if err != nil {
			return nil, fmt.Errorf("bitmap %q: %w", bitmapName, err)
		}

		fileName := sanitizeFileName(bitmapName) + ".idx"
		indexedPath := filepath.Join(outputDir, fileName)
		if err := os.WriteFile(indexedPath, indexedPixels, 0o644); err != nil {
			return nil, err
		}

		metadata.Bitmaps = append(metadata.Bitmaps, indexedBitmapMetadata{
			Name:        bitmapName,
			SourcePath:  sprites.BitmapSources[bitmapName],
			Width:       bitmap.W,
			Height:      bitmap.H,
			IndexedPath: indexedPath,
		})
	}

	metadataPath := filepath.Join(outputDir, "metadata.json")
	if err := writeJSONFile(metadataPath, metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}

func findBitmapByName(sprites core.Sprites, bitmapName string) *core.Bitmap {
	for _, sprite := range sprites.Sprites {
		if sprite == nil || sprite.Bitmap == nil {
			continue
		}
		if sprite.Bitmap.Name == bitmapName {
			return sprite.Bitmap
		}
	}
	return nil
}

func buildIndexedBitmap(bitmap *core.Bitmap, colorToIndex map[uint32]byte) ([]byte, error) {
	indexedPixels := make([]byte, int(bitmap.W*bitmap.H))

	for i := 0; i < len(indexedPixels); i++ {
		pixelOffset := i * 4
		a := bitmap.Pixels[pixelOffset+3]
		if a < 128 {
			indexedPixels[i] = 0
			continue
		}

		color := packRGBA(
			bitmap.Pixels[pixelOffset],
			bitmap.Pixels[pixelOffset+1],
			bitmap.Pixels[pixelOffset+2],
			255,
		)
		index, ok := colorToIndex[color]
		if !ok {
			return nil, fmt.Errorf("pixel color %s was not present in the referenced sprite palette", colorString(color))
		}
		indexedPixels[i] = index
	}

	return indexedPixels, nil
}

func writeJSONFile(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func writePaletteFile(path string, palette []paletteEntry) error {
	data := make([]byte, 0, len(palette)*3)
	for _, entry := range palette {
		data = append(data, entry.R, entry.G, entry.B)
	}
	return os.WriteFile(path, data, 0o644)
}

func sanitizeFileName(name string) string {
	replacer := strings.NewReplacer("/", "_", "\\", "_", " ", "_")
	return replacer.Replace(name)
}

func colorString(color uint32) string {
	r, g, b, a := unpackRGBA(color)
	return fmt.Sprintf("rgba(%d,%d,%d,%d)", r, g, b, a)
}

func packRGBA(r, g, b, a byte) uint32 {
	return uint32(r)<<24 | uint32(g)<<16 | uint32(b)<<8 | uint32(a)
}

func unpackRGBA(color uint32) (byte, byte, byte, byte) {
	return byte(color >> 24), byte(color >> 16), byte(color >> 8), byte(color)
}
