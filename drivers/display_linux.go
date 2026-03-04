package drivers

import (
	"encoding/binary"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
)

type inputEvent struct {
	Type  uint16
	Code  uint16
	Value int32
}

const (
	KeyEnter = 28
	KeyUp    = 103
	KeyDown  = 108
)

var sw, sh int

type Display struct {
	file         *os.File
	keyboardFile *os.File
	pixels       []byte
	buffer       []byte
}

func InitDisplay(vw, vh int) *Display {
	sw, sh = getDisplaySize()
	f, err := os.OpenFile("/dev/fb0", os.O_RDWR, 0)
	if err != nil {
		panic(err)
	}

	size := sw * sh * 4
	log.Printf("Screen size: %dx%d (%d)", sw, sh, size)

	data, err := syscall.Mmap(int(f.Fd()), 0, size, syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	backBuffer := make([]byte, size)

	if err != nil {
		panic("Error en Mmap: " + err.Error())
	}

	return &Display{
		file:   f,
		pixels: data,
		buffer: backBuffer,
	}
}

func getDisplaySize() (int, int) {
	vsBytes, err := os.ReadFile("/sys/class/graphics/fb0/virtual_size")
	if err != nil {
		panic(err)
	}

	parts := strings.Split(strings.TrimSpace(string(vsBytes)), ",")
	realWidth, _ := strconv.Atoi(parts[0])
	realHeight, _ := strconv.Atoi(parts[1])

	return realWidth, realHeight
}

func (d *Display) DrawPixel(vx, vy int32, c []byte) {
	if vx < 0 || vx >= VW || vy < 0 || vy >= VH {
		return
	}

	scaleX, scaleY := sw/VW, sh/VH

	r, g, b, a := c[0], c[1], c[2], c[3]

	for py := 0; py < scaleY; py++ {
		for px := 0; px < scaleX; px++ {
			rx, ry := int(vx)*scaleX+px, int(vy)*scaleY+py

			if rx >= 0 && rx < sw && ry >= 0 && ry < sh {
				offset := (ry*sw + rx) * 4

				// AQUÍ ESTÁ EL CAMBIO: Escribimos en orden B, G, R, A
				d.buffer[offset] = b   // Azul primero
				d.buffer[offset+1] = g // Verde igual
				d.buffer[offset+2] = r // Rojo al final
				d.buffer[offset+3] = a // Alpha
			}
		}
	}
}

func (d *Display) Clear() {
	for i := range d.buffer {
		d.buffer[i] = 0
	}
}

func (d *Display) Present() {
	copy(d.pixels, d.buffer)
}

func (d *Display) GetInput() (int32, bool, bool) {
	if d.keyboardFile == nil {
		return 0, false, false
	}

	b := make([]byte, 24)
	n, err := d.keyboardFile.Read(b)

	if err != nil || n < 24 {
		return 0, false, false
	}

	evType := binary.LittleEndian.Uint16(b[16:18])
	evCode := binary.LittleEndian.Uint16(b[18:20])
	evValue := int32(binary.LittleEndian.Uint32(b[20:24]))

	if evType == 1 && (evValue == 1 || evValue == 2) {
		switch evCode {
		case 103: // Flecha ARRIBA
			return -1, true, false
		case 108: // Flecha ABAJO
			return 1, true, false
		case 28: // ENTER
			return 0, false, true
		}
	}

	return 0, false, false
}

func (d *Display) Close() {
	for i := range d.pixels {
		d.pixels[i] = 0
	}

	syscall.Munmap(d.pixels)

	if d.file != nil {
		d.file.Close()
	}
}
