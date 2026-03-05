package drivers

import (
	"encoding/binary"
	"log"
	"os"
	"strings"
	"syscall"
)

func (d *Display) GetInput() KeyboardInput {
	if d.keyboardFile == nil {
		return KBD_NONE
	}

	buffer := make([]byte, 256)
	n, err := syscall.Read(int(d.keyboardFile.Fd()), buffer)

	if err != nil || n < 24 {
		return KBD_NONE
	}

	for i := 0; i+24 <= n; i += 24 {
		chunk := buffer[i : i+24]

		typ := binary.LittleEndian.Uint16(chunk[16:18])
		code := binary.LittleEndian.Uint16(chunk[18:20])
		val := binary.LittleEndian.Uint32(chunk[20:24])

		if typ == 1 {
			if val == 1 || val == 2 { // Pulsado o mantenido
				switch code {
				case 1:
					return KBD_ESCAPE
				case 28:
					return KBD_RETURN
				case 75:
					return KBD_SPACE
				case 103:
					return KBD_UP
				case 108:
					return KBD_DOWN
				case 105:
					return KBD_LEFT
				case 106:
					return KBD_RIGHT
				default:
					return KBD_NONE
				}
			} else if val == 0 {
				return KBD_NONE
			}
		}
	}

	return KBD_NONE
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
