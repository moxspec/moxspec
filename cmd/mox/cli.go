package main

import (
	"flag"
	"fmt"
	"path/filepath"
)

type app struct {
	name  string
	cmd   string
	args  []string
	fset  *flag.FlagSet
	flags []*genericFlag
	fmap  map[string]*genericFlag
}

func (a *app) appendFlag(key string, def interface{}, desc string) {
	f := newGenericFlag(key, def, desc)
	a.flags = append(a.flags, f)
	a.fmap[key] = f

	switch f.def.(type) {
	case string:
		f.value = a.fset.String(f.key, f.def.(string), f.desc)
	case bool:
		f.value = a.fset.Bool(f.key, f.def.(bool), f.desc)
	}
}

func (a *app) getString(key string) string {
	f, ok := a.fmap[key]
	if !ok {
		return ""
	}
	if _, ok := f.value.(*string); ok {
		return *(f.value.(*string))
	}
	return ""
}

func (a *app) getBool(key string) bool {
	f, ok := a.fmap[key]
	if !ok {
		return false
	}
	if _, ok := f.value.(*bool); ok {
		return *(f.value.(*bool))
	}
	return false
}

func (a *app) parse() error {
	return a.fset.Parse(a.args)
}

func newApp(args []string) *app {
	if len(args) < 2 {
		return nil
	}
	a := new(app)
	a.name = filepath.Base(args[0])
	a.cmd = args[1]
	a.args = args[2:]
	a.fset = flag.NewFlagSet(fmt.Sprintf("%s %s", a.name, a.cmd), flag.ExitOnError)
	a.fmap = make(map[string]*genericFlag)
	return a
}

func newAppWithoutCmd(args []string) *app {
	a := new(app)
	a.name = filepath.Base(args[0])
	a.cmd = ""
	a.args = args[1:]
	a.fset = flag.NewFlagSet(fmt.Sprintf("%s %s", a.name, a.cmd), flag.ExitOnError)
	a.fmap = make(map[string]*genericFlag)
	return a
}

type genericFlag struct {
	key   string
	desc  string
	def   interface{}
	value interface{}
}

func newGenericFlag(key string, def interface{}, desc string) *genericFlag {
	g := new(genericFlag)
	g.key = key
	g.def = def
	g.desc = desc
	return g
}
