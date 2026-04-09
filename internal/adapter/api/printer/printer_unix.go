//go:build linux || darwin

package printer

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// FindAllAvailablePrinters mendeteksi semua printer yang tersedia di sistem operasi.
// Metode ini dapat mendeteksi printer dari berbagai koneksi (USB, LAN, IPP, dll)
// karena library printer mengandalkan driver OS yang sudah terpasang.
func (a *Printer) FindAllAvailablePrinters() ([]string, error) {
	cmd := exec.Command("lpstat", "-a")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("gagal menjalankan perintah lpstat: %w", err)
	}
	output := out.String()
	lines := strings.Split(output, "\n")
	var printerList []string
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) > 0 {
			printerName := strings.Trim(parts[0], ":")
			printerList = append(printerList, printerName)
		}
	}
	if len(printerList) == 0 {
		return nil, fmt.Errorf("tidak ada printer yang terdeteksi di Unix")
	}
	return printerList, nil
}

func (a *Printer) Print(content []byte, printerName string) error {
	cmd := exec.Command("lp", "-d", printerName, "-o", "raw")
	cmd.Stdin = bytes.NewReader(content)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gagal print ke %s: %w", printerName, err)
	}
	return nil
}
