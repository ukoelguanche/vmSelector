package drivers

import (
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

type Display struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	tex      *sdl.Texture
	pixels   []byte
}

func InitDisplay(sw, sh, vw, vh int) *Display {
	sdl.Init(sdl.INIT_EVERYTHING)
	w, _ := sdl.CreateWindow("framebuffer", 100, 100, int32(sw), int32(sh), sdl.WINDOW_SHOWN)
	r, _ := sdl.CreateRenderer(w, -1, sdl.RENDERER_ACCELERATED)
	t, _ := r.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(vw), int32(vh))
	return &Display{w, r, t, make([]byte, vw*vh*4)}
}

func (d *Display) DrawPixel(x, y int32, c []byte) {
	if x < 0 || x >= VW || y < 0 || y >= VH {
		return
	}

	offset := (y*VW + x) * 4
	copy(d.pixels[offset:offset+4], c)

}

func (d *Display) Clear() {
	for i := range d.pixels {
		d.pixels[i] = 0
	}
}

func (d *Display) Present() {
	d.tex.Update(nil, unsafe.Pointer(&d.pixels[0]), VW*4)
	d.renderer.Copy(d.tex, nil, nil)
	d.renderer.Present()
}

func (d *Display) GetInput() (int32, int32, bool) {
	var dx, dy int32
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		if t, ok := event.(*sdl.KeyboardEvent); ok && t.Type == sdl.KEYDOWN {
			switch t.Keysym.Sym {
			case sdl.K_UP:
				dy = -4
			case sdl.K_DOWN:
				dy = 4
			case sdl.K_LEFT:
				dx = -4
			case sdl.K_RIGHT:
				dx = 4
			case sdl.K_ESCAPE:
				return 0, 0, true
			}
		}
		if _, ok := event.(*sdl.QuitEvent); ok {
			return 0, 0, true
		}
	}
	return dx, dy, false
}

func (d *Display) Close() {
	d.window.Destroy()
	sdl.Quit()
}
