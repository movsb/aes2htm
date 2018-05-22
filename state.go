package main

import (
	"fmt"
	"io"
)

type State struct {
	bold      bool
	italic    bool
	underline bool
	fgcolor   int
	bgcolor   int
}

func (o *State) Init() {
	o.bold = false
	o.italic = false
	o.underline = false
	o.fgcolor = -1
	o.bgcolor = -1
}

func (o *State) Empty() bool {
	s := State{}
	s.Init()
	return *o == s
}

func (o *State) AnyChange(s *State) bool {
	return o.bold != s.bold ||
		o.italic != s.italic ||
		o.underline != s.underline ||
		o.fgcolor != s.fgcolor ||
		o.bgcolor != s.bgcolor
}

func (o *State) AnyClose(s *State) bool {
	if s.bold && !o.bold {
		return true
	}
	if s.italic && !o.italic {
		return true
	}
	if s.underline && !o.underline {
		return true
	}
	if s.fgcolor != -1 && o.fgcolor != s.fgcolor {
		return true
	}
	if s.bgcolor != -1 && o.bgcolor != s.bgcolor {
		return true
	}
	return false
}

func (o *State) WriteStyles(w io.Writer, s *State) error {
	var err error

	f := func(b bool, s string) {
		if err == nil && b {
			_, err = w.Write([]byte(s))
		}
	}

	f(o.bold, "font-weight:bold;")
	f(o.italic, "font-style:italic;")
	f(o.underline, "text-decoration:underline;")

	if o.fgcolor != -1 && o.fgcolor != s.fgcolor {
		if color := o.fgcolor - 30; 0 <= color && color <= 7 {
			f(true, "color:"+Palette[color]+";")
		} else {
			panic("bad fg color: " + fmt.Sprint(o.fgcolor))
		}
	}

	if o.bgcolor != -1 && o.bgcolor != s.bgcolor {
		if color := o.bgcolor - 40; 0 <= color && color <= 7 {
			f(true, "background-color:"+Palette[color]+";")
		} else {
			panic("bad bg color: " + fmt.Sprint(o.bgcolor))
		}
	}

	return err
}
