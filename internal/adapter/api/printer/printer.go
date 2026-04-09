package printer

import "print-agent/internal/core/printer"

type Printer struct{}

func NewPrinter() printer.Printer {
	return &Printer{}
}
