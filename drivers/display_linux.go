package drivers

import (
	"encoding/binary"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/term"
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
var oldState *term.State

type Display struct {
	file         *os.File
	keyboardFile *os.File
	pixels       []byte
	buffer       []byte
	LineLength   int // <--- Añade esto
	VW, VH       int // <--- Y esto
}

func InitDisplay(vw, vh int) *Display {
	sw, sh = getDisplaySize()
	f, err := os.OpenFile("/dev/fb0", os.O_RDWR, 0)
	if err != nil {
		panic(err)
	}

	// OBTENER LINE LENGTH REAL
	// Esto evita que la imagen se vea "hacia un lado" o "torcida"
	var fixInfo struct {
		id                            [16]byte
		smem_start                    uintptr
		smem_len                      uint32
		type_                         uint32
		type_aux                      uint32
		visual                        uint32
		xpanstep, ypanstep, ywrapstep uint16
		line_length                   uint32 // Este es el que nos importa
	}
	// FBIOGET_FSCREENINFO = 0x4602
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), 0x4602, uintptr(unsafe.Pointer(&fixInfo)))

	lineLen := int(fixInfo.line_length)
	size := lineLen * sh // Tamaño real de la memoria de video

	data, _ := syscall.Mmap(int(f.Fd()), 0, size, syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	backBuffer := make([]byte, size)

	kbdPath := findKeyboardDevice()
	log.Printf("Keyboard file is: %s", kbdPath)
	keyboardFile, err := os.OpenFile(kbdPath, os.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		log.Fatalf("Failed to open keyboard file: %s", err)
	}

	fd := int(os.Stdin.Fd())

	oldState, err = term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}

	return &Display{
		file:         f,
		pixels:       data,
		buffer:       backBuffer,
		LineLength:   lineLen,
		VW:           vw,
		VH:           vh,
		keyboardFile: keyboardFile,
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
	// Usamos las constantes o variables VW y VH que tengas definidas
	if vx < 0 || vx >= int32(VW) || vy < 0 || vy >= int32(VH) {
		return
	}

	// Proyección dinámica: calculamos el área real que ocupa el píxel virtual
	// Esto reparte los píxeles sobrantes automáticamente
	xStart := int(float64(vx) * float64(sw) / float64(VW))
	xEnd := int(float64(vx+1) * float64(sw) / float64(VW))

	yStart := int(float64(vy) * float64(sh) / float64(VH))
	yEnd := int(float64(vy+1) * float64(sh) / float64(VH))

	r, g, b, a := c[0], c[1], c[2], c[3]

	// Dibujamos el bloque estirado
	for py := yStart; py < yEnd; py++ {
		for px := xStart; px < xEnd; px++ {
			// Importante: No olvides que si Alpine usa LineLength,
			// deberías usar d.LineLength en lugar de sw aquí.
			offset := (py*sw + px) * 4

			if offset+3 < len(d.buffer) {
				d.buffer[offset] = b
				d.buffer[offset+1] = g
				d.buffer[offset+2] = r
				d.buffer[offset+3] = a
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

func (d *Display) GetInput() (int, bool, bool) {
	if d.keyboardFile == nil {
		return 0, false, false
	}

	buffer := make([]byte, 256)
	n, err := syscall.Read(int(d.keyboardFile.Fd()), buffer)

	// Si no hay datos, devolvemos todo en falso/cero inmediatamente
	if err != nil || n < 24 {
		return 0, false, false
	}

	var lastCode uint16
	var isQuit, isEnter bool
	foundKey := false

	// Recorremos los eventos que han llegado
	for i := 0; i+24 <= n; i += 24 {
		chunk := buffer[i : i+24]

		typ := binary.LittleEndian.Uint16(chunk[16:18])
		code := binary.LittleEndian.Uint16(chunk[18:20])
		val := binary.LittleEndian.Uint32(chunk[20:24])

		if typ == 1 { // EV_KEY
			if val == 1 || val == 2 { // Pulsado o mantenido
				lastCode = code
				foundKey = true
				if code == 1 {
					isQuit = true
				}
				if code == 28 {
					isEnter = true
				}
			} else if val == 0 {
				// Si quieres que Sonic se pare al soltar la tecla:
				// foundKey = false
			}
		}
	}

	if foundKey {
		return int(lastCode), isQuit, isEnter
	}

	return 0, false, false
}

func findKeyboardDevice() string {
	data, err := os.ReadFile("/proc/bus/input/devices")
	if err != nil {
		log.Printf("error reading /proc/bus/input/devices: %v", err)
		return "/dev/input/event2" // Tu sospechoso principal
	}

	sections := strings.Split(string(data), "\n\n")
	for _, section := range sections {
		// 1. Que tenga el nombre de tu teclado
		// 2. Y que en Handlers aparezca "kbd" (esto descarta los que son solo ratón o control)
		if strings.Contains(section, "Gaming KB") && strings.Contains(section, "kbd") {
			lines := strings.Split(section, "\n")
			for _, line := range lines {
				if strings.Contains(line, "Handlers=") {
					// Buscamos el eventX que esté en esta línea
					parts := strings.Fields(line)
					for _, p := range parts {
						if strings.HasPrefix(p, "event") {
							log.Printf("Returning /dev/input/%s", p)
							return "/dev/input/" + p
						}
					}
				}
			}
		}
	}
	log.Printf("fallback to /dev/input/event2")
	return "/dev/input/event2"
}

func (d *Display) Close() {
	d.Clear()
	for i := range d.pixels {
		d.pixels[i] = 0
	}

	syscall.Munmap(d.pixels)

	if d.file != nil {
		d.file.Close()
	}

	if oldState != nil {
		term.Restore(int(os.Stdin.Fd()), oldState)
	}

	syscall.Munmap(d.pixels)
}
