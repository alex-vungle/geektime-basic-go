package main

import "fmt"

type Printer interface {
	PrintMessage() string
}

type OldPrinter struct{}

func (p *OldPrinter) Print() string {
	return "Old Printer: PrintMessage"
}

type PrinterAdapter struct {
	OldPrinter *OldPrinter
}

func (p *PrinterAdapter) PrintMessage() string {
	if p.OldPrinter != nil {
		return p.OldPrinter.Print()
	}
	return "PrinterAdapter: No old printer"
}

func invoke(p Printer) {
	fmt.Println(p.PrintMessage())
}

func main() {
	oldPrinter := &OldPrinter{}
	adapter := &PrinterAdapter{OldPrinter: oldPrinter}

	var printer Printer = adapter

	invoke(printer)
}
