package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const margin = 2

type printer struct {
	writer   io.Writer
	headlen  int
	sections []*section
}

func newPrinter() *printer {
	p := new(printer)
	p.setOutput(os.Stdout)
	return p
}

func (p *printer) setOutput(w io.Writer) {
	p.writer = w
}

func (p *printer) append(s *section) {
	if s.block == nil || len(s.block.contents) == 0 {
		return
	}
	p.sections = append(p.sections, s)
	if p.headlen < len(s.name)+margin {
		p.headlen = len(s.name) + margin
	}
}

func (p printer) show() {
	for _, s := range p.sections {
		s.show(p.writer, p.headlen)
	}
}

type section struct {
	name  string
	block *block
}

func newSection(name string) *section {
	s := new(section)
	s.name = name + ":"
	s.block = new(block)
	return s
}

func (s section) show(w io.Writer, indent int) {
	if s.block == nil {
		return
	}
	if len(s.block.contents) == 0 {
		return
	}

	fmtr := "%-" + fmt.Sprintf("%d", indent) + "s "
	fmt.Fprintf(w, fmtr, s.name)
	s.block.root = true
	s.block.show(w, indent+1)
}

type block struct {
	root     bool
	indent   int
	contents []interface{}
}

func newGroupedBlock(title string, blk ...*block) *block {
	g := new(block)
	g.append(title)

	var hasContents bool
	for _, b := range blk {
		if b == nil {
			continue
		}

		g.append(b)
		if !hasContents && b.hasContents() {
			hasContents = true
		}
	}

	if hasContents {
		return g
	}
	return nil
}

func newIndentedBlock(lines []string) *block {
	g := new(block)
	for _, l := range lines {
		g.append(l)
	}
	return g
}

func (b block) hasContents() bool {
	return (len(b.contents) > 0)
}

func (b *block) appendf(fmtr string, str ...interface{}) {
	b.contents = append(b.contents, fmt.Sprintf(fmtr, str...))
}

func (b *block) append(c interface{}) {
	switch c.(type) {
	case string:
		if c.(string) == "" {
			return
		}
	case []string:
		for _, m := range c.([]string) {
			if m == "" {
				continue
			}
			b.contents = append(b.contents, m)
		}
		return
	case *block:
		if c.(*block) == nil {
			return
		}
		if !c.(*block).hasContents() {
			return
		}
	}
	b.contents = append(b.contents, c)
}

func (b block) show(w io.Writer, indent int) {
	for i, c := range b.contents {
		switch c.(type) {
		case string:
			if (b.root && i != 0) || !b.root {
				fmt.Fprintf(w, strings.Repeat(" ", indent))
			}
			fmt.Fprintln(w, c.(string))
		case *block:
			c.(*block).show(w, indent+2)
		}
	}
}
