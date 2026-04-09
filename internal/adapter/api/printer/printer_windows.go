//go:build windows

package printer

import (
	"fmt"

	"github.com/alexbrainman/printer"
)

// FindAllAvailablePrinters mendeteksi semua printer yang tersedia di sistem operasi.
// Metode ini dapat mendeteksi printer dari berbagai koneksi (USB, LAN, IPP, dll)
// karena library printer mengandalkan driver OS yang sudah terpasang.
func (a *Printer) FindAllAvailablePrinters() ([]string, error) {
	printerList, err := printer.ReadNames()
	if err != nil {
		return nil, fmt.Errorf("gagal mendeteksi printer di Windows: %w", err)
	}
	if len(printerList) == 0 {
		return nil, fmt.Errorf("tidak ada printer yang terdeteksi di Windows")
	}
	return printerList, nil
}

func (a *Printer) Print(content []byte, printerName string) error {
	p, err := printer.Open(printerName)
	if err != nil {
		return fmt.Errorf("gagal membuka printer %s: %w", printerName, err)
	}
	defer p.Close()

	if err := p.StartDocument("ESC/POS Print Job", "RAW"); err != nil {
		return fmt.Errorf("gagal memulai dokumen: %w", err)
	}
	defer p.EndDocument()

	if err := p.StartPage(); err != nil {
		return fmt.Errorf("gagal memulai halaman: %w", err)
	}

	_, err = p.Write(content)
	if err != nil {
		p.EndPage()
		return fmt.Errorf("gagal menulis ke printer: %w", err)
	}

	if err := p.EndPage(); err != nil {
		return fmt.Errorf("gagal mengakhiri halaman: %w", err)
	}

	return nil
}
