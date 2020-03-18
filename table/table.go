package table

import (
	"fmt"
	"github.com/alexeyco/simpletable"
)

func MakeTable() *simpletable.Table {
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "HOST"},
			{Align: simpletable.AlignCenter, Text: "IP"},
		},
	}

	table.SetStyle(simpletable.StyleCompactLite)
	return table
}

func InsertHostData(table simpletable.Table, i int, hostName string, ip string) {
	r := []*simpletable.Cell{
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", i)},
		{Text: hostName},
		{Align: simpletable.AlignRight, Text: ip},
	}

	table.Body.Cells = append(table.Body.Cells, r)
	i += 1
}
