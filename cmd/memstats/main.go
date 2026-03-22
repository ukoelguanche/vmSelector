package main

import (
	"fmt"
	_ "image/png"
	"io"
	"log"
	"runtime"
	"runtime/debug"
	"unsafe"

	"apodeiktikos.com/fbtest/engine"
	"apodeiktikos.com/fbtest/interfaces"
	"apodeiktikos.com/fbtest/manager"
	"github.com/ukoelguanche/graphicsengine/core"
	"github.com/ukoelguanche/graphicsengine/loaders"
)

const (
	spritesFile         = "./assets/sprites/Sprites.json"
	indexedMetadataFile = "./assets/indexed/metadata.json"
)

type memoryTotals struct {
	bitmapsActual int
	bitmapsRGBA   int
	palettes      int
	count         int
}

func main() {
	log.SetOutput(io.Discard)
	runtime.GC()
	debug.FreeOSMemory()

	var sprites core.Sprites
	loaders.LoadIndexedSprites(spritesFile, indexedMetadataFile, &sprites)

	var renderables []interfaces.Renderable
	renderables = manager.SetupClouds(sprites, renderables)
	renderables = manager.SetupGreenHillBackground(sprites, renderables)
	renderables = manager.SetupGreenHillForeground(sprites, renderables)
	renderables = append(renderables, manager.SetupSonic(sprites))

	baseTotals := collectBitmapTotals(baseSpriteBitmaps(sprites))
	paletteTotals := collectBitmapTotals(paletteVariantBitmaps(sprites))
	cachedTotals := collectBitmapTotals(cachedLayerBitmaps(renderables))

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	fmt.Println("Memory report (safe subset, no Proxmox calls)")
	fmt.Printf("Base sprite bitmaps:   actual=%s  rgba-equivalent=%s  palettes=%s  unique-bitmaps=%d\n",
		formatBytes(baseTotals.bitmapsActual), formatBytes(baseTotals.bitmapsRGBA), formatBytes(baseTotals.palettes), baseTotals.count)
	fmt.Printf("Palette variants:      actual=%s  rgba-equivalent=%s  palettes=%s  unique-bitmaps=%d\n",
		formatBytes(paletteTotals.bitmapsActual), formatBytes(paletteTotals.bitmapsRGBA), formatBytes(paletteTotals.palettes), paletteTotals.count)
	fmt.Printf("Cached layers:         actual=%s  rgba-equivalent=%s  palettes=%s  unique-bitmaps=%d\n",
		formatBytes(cachedTotals.bitmapsActual), formatBytes(cachedTotals.bitmapsRGBA), formatBytes(cachedTotals.palettes), cachedTotals.count)

	totalActual := baseTotals.bitmapsActual + baseTotals.palettes + paletteTotals.bitmapsActual + paletteTotals.palettes + cachedTotals.bitmapsActual + cachedTotals.palettes
	totalRGBA := baseTotals.bitmapsRGBA + paletteTotals.bitmapsRGBA + cachedTotals.bitmapsRGBA
	fmt.Printf("Graphics subtotal:     actual=%s  rgba-equivalent=%s  savings=%s\n",
		formatBytes(totalActual), formatBytes(totalRGBA), formatBytes(totalRGBA-totalActual))

	fmt.Println()
	fmt.Printf("Go heap alloc:         %s\n", formatBytes(int(mem.Alloc)))
	fmt.Printf("Go heap sys:           %s\n", formatBytes(int(mem.HeapSys)))
	fmt.Printf("Go total sys:          %s\n", formatBytes(int(mem.Sys)))
}

func baseSpriteBitmaps(sprites core.Sprites) []*core.Bitmap {
	result := make([]*core.Bitmap, 0, len(sprites.Sprites))
	for _, sprite := range sprites.Sprites {
		if sprite == nil || sprite.Bitmap == nil {
			continue
		}
		result = append(result, sprite.Bitmap)
	}
	return result
}

func paletteVariantBitmaps(sprites core.Sprites) []*core.Bitmap {
	result := make([]*core.Bitmap, 0)
	for _, sprite := range sprites.Sprites {
		if sprite == nil || len(sprite.PaletteBitmaps) == 0 {
			continue
		}
		result = append(result, sprite.PaletteBitmaps...)
	}
	return result
}

func cachedLayerBitmaps(renderables []interfaces.Renderable) []*core.Bitmap {
	result := make([]*core.Bitmap, 0)
	for _, renderable := range renderables {
		layer, ok := renderable.(*engine.CachedLayer)
		if !ok {
			continue
		}
		if layer.GetSprite() == nil || layer.GetSprite().Bitmap == nil {
			continue
		}
		result = append(result, layer.GetSprite().Bitmap)
	}
	return result
}

func collectBitmapTotals(bitmaps []*core.Bitmap) memoryTotals {
	seenBitmaps := make(map[*core.Bitmap]struct{})
	seenPalettes := make(map[*core.Palette]struct{})
	totals := memoryTotals{}

	for _, bitmap := range bitmaps {
		if bitmap == nil {
			continue
		}
		if _, seen := seenBitmaps[bitmap]; seen {
			continue
		}
		seenBitmaps[bitmap] = struct{}{}
		totals.count++

		totals.bitmapsActual += len(bitmap.Pixels) + len(bitmap.IndexedPixels)
		totals.bitmapsRGBA += int(bitmap.W * bitmap.H * 4)

		if bitmap.Palette != nil {
			if _, seen := seenPalettes[bitmap.Palette]; !seen {
				seenPalettes[bitmap.Palette] = struct{}{}
				totals.palettes += len(*bitmap.Palette) * int(unsafe.Sizeof(core.Color{}))
			}
		}
	}

	return totals
}

func formatBytes(n int) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}

	units := []string{"B", "KiB", "MiB", "GiB", "TiB"}
	value := float64(n)
	unitIndex := 0
	for value >= unit && unitIndex < len(units)-1 {
		value /= unit
		unitIndex++
	}

	return fmt.Sprintf("%.2f %s", value, units[unitIndex])
}
