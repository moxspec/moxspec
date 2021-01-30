package main

import (
	"fmt"
	"strings"
)

type table struct {
	headers []string
	widths  []int
	rows    [][]string
}

func newTable(h string, hs ...string) *table {
	t := new(table)
	t.headers = append([]string{h}, hs...)
	for _, w := range t.headers {
		t.widths = append(t.widths, len(w))
	}
	return t
}

func (t *table) append(d string, ds ...string) {
	row := append([]string{d}, ds...)
	if len(t.headers) != len(row) {
		return
	}

	t.rows = append(t.rows, row)

	for i := 0; i < len(t.widths); i++ {
		if t.widths[i] < len(row[i]) {
			t.widths[i] = len(row[i])
		}
	}
}

func (t table) border() string {
	line := "+"
	for _, l := range t.widths {
		line += fmt.Sprintf("-%s-+", strings.Repeat("-", l))
	}
	line += "\n"

	return line
}

func (t table) print() {
	fmtr := "|"
	for _, l := range t.widths {
		fmtr += fmt.Sprintf(" %%-%ds |", l)
	}
	fmtr += "\n"

	border := t.border()

	fmt.Printf(border)
	fmt.Printf(fmtr, strSliceToIntfSlice(t.headers)...)
	fmt.Printf(border)
	for _, r := range t.rows {
		fmt.Printf(fmtr, strSliceToIntfSlice(r)...)
	}
	fmt.Printf(border)
}
