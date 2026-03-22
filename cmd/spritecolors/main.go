package main

import (
	_ "image/png"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"github.com/ukoelguanche/graphicsengine/core"
	"github.com/ukoelguanche/graphicsengine/loaders"
)

const defaultSpritesFile = "./assets/sprites/Sprites.json"

type spriteReport struct {
	Name       string
	BitmapName string
	ColorCount int
}

func main() {
	spritesPath := flag.String("sprites", defaultSpritesFile, "path to Sprites.json")
	verbose := flag.Bool("verbose", false, "print loader logs")
	flag.Parse()

	if !*verbose {
		log.SetOutput(io.Discard)
	}

	var sprites core.Sprites
	loaders.LoadSprites(*spritesPath, &sprites)

	reports, globalColors, err := analyzeSprites(sprites)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sprite analysis failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Sprite file: %s\n", *spritesPath)
	fmt.Printf("Sprites analyzed: %d\n", len(reports))
	fmt.Printf("Minimum global palette size for all referenced sprite frames: %d colors\n", len(globalColors))
	fmt.Println()
	fmt.Println("Per-sprite color counts:")

	for _, report := range reports {
		fmt.Printf("- %-24s %4d colors (%s)\n", report.Name, report.ColorCount, report.BitmapName)
	}
}

func analyzeSprites(sprites core.Sprites) ([]spriteReport, map[uint32]struct{}, error) {
	names := make([]string, 0, len(sprites.Sprites))
	for name := range sprites.Sprites {
		names = append(names, name)
	}
	sort.Strings(names)

	reports := make([]spriteReport, 0, len(names))
	globalColors := make(map[uint32]struct{})

	for _, name := range names {
		sprite := sprites.Sprites[name]
		if sprite == nil {
			return nil, nil, fmt.Errorf("sprite %q is nil", name)
		}
		if sprite.Bitmap == nil {
			return nil, nil, fmt.Errorf("sprite %q has no bitmap assigned", name)
		}

		spriteColors := collectSpriteColors(sprite)
		for color := range spriteColors {
			globalColors[color] = struct{}{}
		}

		reports = append(reports, spriteReport{
			Name:       name,
			BitmapName: sprite.Bitmap.Name,
			ColorCount: len(spriteColors),
		})
	}

	return reports, globalColors, nil
}

func collectSpriteColors(sprite *core.Sprite) map[uint32]struct{} {
	colors := make(map[uint32]struct{})
	bitmap := sprite.Bitmap
	bitmapW := int(bitmap.W)
	bitmapH := int(bitmap.H)

	for _, frame := range sprite.Frames {
		startX := maxInt(0, int(frame.Point.X))
		startY := maxInt(0, int(frame.Point.Y))
		endX := minInt(bitmapW, int(frame.Point.X+frame.Size.W))
		endY := minInt(bitmapH, int(frame.Point.Y+frame.Size.H))

		for y := startY; y < endY; y++ {
			rowOffset := y * bitmapW * 4
			for x := startX; x < endX; x++ {
				offset := rowOffset + x*4
				color := packRGBA(
					bitmap.Pixels[offset],
					bitmap.Pixels[offset+1],
					bitmap.Pixels[offset+2],
					bitmap.Pixels[offset+3],
				)
				colors[color] = struct{}{}
			}
		}
	}

	return colors
}

func packRGBA(r, g, b, a byte) uint32 {
	return uint32(r)<<24 | uint32(g)<<16 | uint32(b)<<8 | uint32(a)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
