package printer

type Printer interface {
	FindAllAvailablePrinters() ([]string, error)
	Print(content []byte, printerName string) error
}
