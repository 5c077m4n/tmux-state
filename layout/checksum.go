package layout

import "fmt"

// Based on https://github.com/tmux/tmux/blob/851c5a933d4838c32ad06c248b2ba975d106149c/layout-custom.c#L47
func calculateTmuxChecksum(layout string) uint16 {
	checksum := uint16(0)
	for i := 0; i < len(layout); i++ {
		checksum = (checksum >> 1) + ((checksum & 1) << 15) + uint16(layout[i])
	}
	return checksum
}

func FormatTmuxLayout(layoutBody string) string {
	checksum := calculateTmuxChecksum(layoutBody)
	return fmt.Sprintf("%04x,%s", checksum, layoutBody)
}
