package invoice

import (
	"context"
	"fmt"
	"print-agent/internal/core/printer"
	"print-agent/internal/core/setting"
	escpos "print-agent/pkg/ecspos"
	"strings"

	"gorm.io/gorm"
)

type UsecaseDependency struct {
	DB             *gorm.DB
	SettingUsecase setting.Usecase
	Printer        printer.Printer
}

type UsecaseImpl struct {
	UsecaseDependency
}

func NewUsecase(deps UsecaseDependency) Usecase {
	return &UsecaseImpl{
		deps,
	}
}

func (u UsecaseImpl) BuildInvoice(ctx context.Context, request Content) error {
	settings, err := u.SettingUsecase.Find(ctx)
	if err != nil {
		return err
	}

	builder := escpos.NewBuilder()

	// Header
	builder.JustifyCenter().BoldOn().AddLine("NOTA PEMBAYARAN").BoldOff()
	builder.AddLine(strings.Repeat("=", 48))

	// Detail Transaksi
	builder.JustifyLeft()
	builder.AddLine(fmt.Sprintf("Dilayani Oleh: %s", request.CashierName))
	builder.AddLine(fmt.Sprintf("Nama Pelanggan: %s", request.PartnerName))
	builder.AddLine(fmt.Sprintf("No. Faktur    : %s", request.InvoiceID))
	builder.AddLine(strings.Repeat("-", 48))

	// Item
	builder.AddLine(fmt.Sprintf("%-18s %s %s %s", "Nama Barang",
		builder.PadOrTrimRight("Qty", 5),
		builder.PadOrTrimRight("Harga", 10),
		builder.PadOrTrimRight("Total", 12)))

	builder.AddLine(strings.Repeat("-", 48))

	// Isi tabel
	for _, item := range request.Items {
		builder.AddLine(builder.FormatItem(item.Description, item.Qty, item.Price, item.Subtotal))
	}

	builder.AddLine(strings.Repeat("-", 48))

	// Total
	//builder.AddLine(fmt.Sprintf("%-25s %s", "Belanja:", escpos.FormatRupiah(request.Total, true)))
	//builder.AddLine(fmt.Sprintf("%-25s %s", fmt.Sprintf("PPN (%.0f%%)", request.TaxPerc), escpos.FormatRupiah(request.Tax, true)))
	//builder.AddLine(fmt.Sprintf("%-25s %s", "Diskon:", escpos.FormatRupiah(request.Discount, true)))
	//builder.AddLine(strings.Repeat("-", 48))

	builder.BoldOn().AddLine(
		escpos.FormatLineJustify("Total:", escpos.FormatRupiah(request.GrandTotal, false), 48),
	).BoldOff()

	builder.BoldOn().AddLine(
		escpos.FormatLineJustify("Dibayar:", escpos.FormatRupiah(request.PayedTotal, false), 48),
	).BoldOff()

	builder.BoldOn().AddLine(
		escpos.FormatLineJustify("Kembali:", escpos.FormatRupiah(request.ReturnTotal, false), 48),
	).BoldOff()

	// Footer
	builder.AddLine(strings.Repeat("=", 48))
	builder.JustifyCenter().AddLine("Terima kasih telah berbelanja!")
	builder.JustifyCenter().AddLine(request.OrderDate)
	builder.AddLine("")      // Baris kosong
	builder.AddLine("")      // Baris kosong
	builder.AddLine("")      // Baris kosong
	builder.AddLine("")      // Baris kosong
	builder.AddLine("")      // Baris kosong
	builder.AddLine("")      // Baris kosong
	builder.AddLine("")      // Baris kosong
	builder.AddLine("")      // Baris kosong
	builder.AddLine("")      // Baris kosong
	builder.AddCommand("\f") // Perintah Form Feed

	// Dapatkan request cetak dalam bentuk byte slice
	content := builder.Build()

	u.Printer.Print(content, settings.ChoosenPrinter)

	return nil
}
