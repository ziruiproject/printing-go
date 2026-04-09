package escpos

import (
	"fmt"
	"strconv"
	"strings"
)

// ESC/POS Commands
const (
	// Control Commands
	ESC = "\x1B"
	GS  = "\x1D"
	NUL = "\x00"

	// Formatting Commands
	BoldOn  = ESC + "\x45\x01" // ESC E 1
	BoldOff = ESC + "\x45\x00" // ESC E 0

	UnderlineOn    = ESC + "\x2D\x01" // ESC - 1
	UnderlineOff   = ESC + "\x2D\x00" // ESC - 0
	UnderlineThick = ESC + "\x2D\x02" // ESC - 2

	FontA = ESC + "\x4D\x00" // ESC M 0
	FontB = ESC + "\x4D\x01" // ESC M 1

	// Justification/Alignment
	AlignLeft   = ESC + "\x61\x00" // ESC a 0
	AlignCenter = ESC + "\x61\x01" // ESC a 1
	AlignRight  = ESC + "\x61\x02" // ESC a 2

	// Line Feed
	LF = "\x0A" // Line Feed

	// Paper Control
	CutFull    = GS + "\x56\x00" // GS V 0
	CutPartial = GS + "\x56\x01" // GS V 1
)

func DoubleSize(width, height bool) string {
	var size byte = 0
	if width {
		size |= 0x10
	}
	if height {
		size |= 0x01
	}
	return GS + "!" + string(size)
}

type PrinterTextBuilder struct {
	content strings.Builder
}

func NewBuilder() *PrinterTextBuilder {
	return &PrinterTextBuilder{}
}

// AddText menambahkan teks biasa.
func (b *PrinterTextBuilder) AddText(s string) *PrinterTextBuilder {
	b.content.WriteString(s)
	return b
}

// AddCommand menambahkan perintah non-cetak.
func (b *PrinterTextBuilder) AddCommand(s string) *PrinterTextBuilder {
	b.content.WriteString(s)
	return b
}

// AddLine menambahkan teks dan diikuti baris baru.
func (b *PrinterTextBuilder) AddLine(s string) *PrinterTextBuilder {
	b.content.WriteString(s)
	b.content.WriteString(LF)
	return b
}

// JustifyLeft mengatur perataan teks ke kiri.
func (b *PrinterTextBuilder) JustifyLeft() *PrinterTextBuilder {
	return b.AddCommand(AlignLeft)
}

// JustifyCenter mengatur perataan teks ke tengah.
func (b *PrinterTextBuilder) JustifyCenter() *PrinterTextBuilder {
	return b.AddCommand(AlignCenter)
}

// JustifyRight mengatur perataan teks ke kanan.
func (b *PrinterTextBuilder) JustifyRight() *PrinterTextBuilder {
	return b.AddCommand(AlignRight)
}

// BoldOn mengaktifkan cetak tebal.
func (b *PrinterTextBuilder) BoldOn() *PrinterTextBuilder {
	return b.AddCommand(BoldOn)
}

// BoldOff menonaktifkan cetak tebal.
func (b *PrinterTextBuilder) BoldOff() *PrinterTextBuilder {
	return b.AddCommand(BoldOff)
}

// CutFull melakukan pemotongan penuh pada kertas.
func (b *PrinterTextBuilder) CutFull() *PrinterTextBuilder {
	return b.AddCommand(CutFull)
}

// padOrTrim padding atau crop teks sesuai width
func (b *PrinterTextBuilder) PadOrTrim(s string, width int) string {
	if len(s) > width {
		return s[:width]
	}
	return s + strings.Repeat(" ", width-len(s))
}

// PadOrTrimLeft -> teks rata kiri
func (b *PrinterTextBuilder) PadOrTrimLeft(s string, width int) string {
	if len(s) > width {
		return s[:width]
	}
	return s + strings.Repeat(" ", width-len(s))
}

// PadOrTrimRight -> teks rata kanan
func (b *PrinterTextBuilder) PadOrTrimRight(s string, width int) string {
	if len(s) > width {
		return s[:width]
	}
	return strings.Repeat(" ", width-len(s)) + s
}

// formatItem ke tabel
func (b *PrinterTextBuilder) FormatItem(desc string, qty int, price, subtotal float64) string {
	col1 := b.PadOrTrimLeft(desc, 18)                           // Deskripsi rata kiri
	col2 := b.PadOrTrimRight(fmt.Sprintf("%d", qty), 5)         // Qty rata kanan
	col3 := b.PadOrTrimRight(FormatRupiah(price, false), 10)    // Harga rata kanan
	col4 := b.PadOrTrimRight(FormatRupiah(subtotal, false), 12) // Total rata kanan
	return fmt.Sprintf("%s %s %s %s", col1, col2, col3, col4)
}

// Build mengembalikan hasil akhir sebagai byte slice.
func (b *PrinterTextBuilder) Build() []byte {
	return []byte(b.content.String())
}

func FormatRupiah(amount float64, withRupiah bool) string {
	// Convert float ke integer (tanpa desimal)
	n := int64(amount)

	// Ubah ke string
	s := strconv.FormatInt(n, 10)

	// Sisipkan titik setiap 3 digit dari belakang
	var result strings.Builder
	length := len(s)
	for i, digit := range s {
		if (length-i)%3 == 0 && i != 0 {
			result.WriteRune('.')
		}
		result.WriteRune(digit)
	}

	if withRupiah {
		return "Rp" + result.String()
	} else {
		return result.String()
	}
}

// FormatLineJustify kiri teks, kanan angka, dengan lebar total line
func FormatLineJustify(left string, right string, width int) string {
	left = strings.TrimSpace(left)
	right = strings.TrimSpace(right)

	// sisa ruang = total width - (len(left) + len(right))
	spaces := width - (len(left) + len(right))
	if spaces < 1 {
		spaces = 1 // minimal 1 spasi
	}

	return left + strings.Repeat(" ", spaces) + right
}
