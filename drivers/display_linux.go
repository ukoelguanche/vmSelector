package drivers

import (
	"encoding/binary"
	"log"
	"os"
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

type Display struct {
	file         *os.File
	keyboardFile *os.File
	pixels       []byte
}

func InitDisplay(sw, sh, vw, vh int) *Display {
	f, err := os.OpenFile("/dev/fb0", os.O_RDWR, 0)
	if err != nil {
		panic(err)
	}

	size := sw * sh * 4
	log.Printf("Screen size: %dx%d (%d)", sw, sh, size)

	data, err := syscall.Mmap(int(f.Fd()), 0, size, syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		panic("Error en Mmap: " + err.Error())
	}

	return &Display{
		file:   f,
		pixels: data,
	}
}

func (d *Display) DrawPixel(vx, vy int32, c []byte) {
	if vx < 0 || vx >= VW || vy < 0 || vy >= VH {
		return
	}

	scaleX, scaleY := SW/VW, SH/VH

	r, g, b, a := c[0], c[1], c[2], c[3]

	for py := 0; py < scaleY; py++ {
		for px := 0; px < scaleX; px++ {
			rx, ry := int(vx)*scaleX+px, int(vy)*scaleY+py

			if rx >= 0 && rx < SW && ry >= 0 && ry < SH {
				offset := (ry*SW + rx) * 4

				// AQUÍ ESTÁ EL CAMBIO: Escribimos en orden B, G, R, A
				d.pixels[offset] = b   // Azul primero
				d.pixels[offset+1] = g // Verde igual
				d.pixels[offset+2] = r // Rojo al final
				d.pixels[offset+3] = a // Alpha
			}
		}
	}
}

func (d *Display) Clear() {
	for i := range d.pixels {
		d.pixels[i] = 0
	}
}

func (d *Display) Present() {
	// ¡CON MMAP NO HACE FALTA HACER NADA AQUÍ!
	// En cuanto escribes en d.pixels[i] = x, el píxel viaja a la pantalla.
	// Solo si ves parpadeo, usaremos un buffer intermedio más adelante.
}

func (d *Display) GetInput() (int32, int32, bool, bool) {
	if d.keyboardFile == nil {
		return 0, 0, false, false
	}

	// El tamaño de input_event en Linux 64-bit es de 24 bytes.
	// [0-15]: Time (Segundos y Microsegundos) -> No los necesitamos ahora
	// [16-17]: Type (EV_KEY, EV_REL, etc.)
	// [18-19]: Code (Código de la tecla)
	// [20-23]: Value (0: soltado, 1: presionado, 2: repetido)
	b := make([]byte, 24)
	n, err := d.keyboardFile.Read(b)

	// Si el archivo está vacío (EAGAIN) o hay error, salimos.
	if err != nil || n < 24 {
		return 0, 0, false, false
	}

	evType := binary.LittleEndian.Uint16(b[16:18])
	evCode := binary.LittleEndian.Uint16(b[18:20])
	evValue := int32(binary.LittleEndian.Uint32(b[20:24]))

	// EV_KEY es siempre 1 en el protocolo de entrada de Linux
	if evType == 1 && (evValue == 1 || evValue == 2) {
		switch evCode {
		case 103: // Flecha ARRIBA
			return 0, -1, true, false
		case 108: // Flecha ABAJO
			return 0, 1, true, false
		case 28: // ENTER
			return 0, 0, false, true
		}
	}

	return 0, 0, false, false
}

func (d *Display) Close() { d.file.Close() }
