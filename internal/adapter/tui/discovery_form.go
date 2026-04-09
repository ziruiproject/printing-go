package tui

import (
	"context"
	"print-agent/internal/core/printer"
	"print-agent/internal/core/setting"

	"github.com/charmbracelet/huh"
)

type DiscoveryForm struct {
	Printer        printer.Printer
	SettingUsecase setting.Usecase
}

func NewDiscoveryForm(printer printer.Printer, setting setting.Usecase) *DiscoveryForm {
	return &DiscoveryForm{
		Printer:        printer,
		SettingUsecase: setting,
	}
}

func (tui *DiscoveryForm) Show() {
	printers, err := tui.Printer.FindAllAvailablePrinters()
	if err != nil {
		return
	}

	var printerOptions []huh.Option[string]
	var choosenPrinter string

	for _, p := range printers {
		printerOptions = append(printerOptions, huh.NewOption(p, p))
	}

	err = huh.NewSelect[string]().
		Title("Choose your printer").
		Options(printerOptions...).
		Value(&choosenPrinter).
		Run()

	if err != nil {
		return
	}

	request := setting.Content{
		ChoosenPrinter: choosenPrinter,
	}
	_, err = tui.SettingUsecase.Save(context.Background(), request)
	if err != nil {
		return
	}
}
