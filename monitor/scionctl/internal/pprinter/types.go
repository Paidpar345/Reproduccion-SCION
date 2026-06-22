package pprinter

import "github.com/jedib0t/go-pretty/v6/table"

type TableInput struct {
	header  table.Row
	content []table.Row
}
