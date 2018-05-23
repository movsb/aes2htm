package main

import (
	"fmt"
	"io"
)

type Color struct {
	typ   int // -1: no, 0: index, 1: absolute
	color uint32
}

func (o *Color) Equal(c *Color) bool {
	return o.typ == -1 && c.typ == -1 ||
		o.typ == c.typ && o.color == c.color
}

func (o *Color) HasColor() bool {
	return o.typ != -1
}

func (o *Color) SetNone() {
	o.typ = -1
}

func (o *Color) SetRGB(r, g, b int) {
	o.typ = 1
	o.color = uint32((r&0xff)<<16 + (g&0xff)<<8 + (b & 0xff))
}

func (o *Color) SetIndex(i int) {
	o.typ = 0
	o.color = uint32(i & 0xff)
}

func (o *Color) String() string {
	if o.typ == 0 {
		return Palette[o.color]
	} else if o.typ == 1 {
		return fmt.Sprintf("#%06x", o.color)
	} else {
		panic("invalid color")
	}
}

type State struct {
	bold      bool
	italic    bool
	underline bool
	blink     bool
	fgcolor   Color
	bgcolor   Color
}

func (o *State) Init() {
	o.bold = false
	o.italic = false
	o.underline = false
	o.blink = false
	o.fgcolor.SetNone()
	o.bgcolor.SetNone()
}

func (o *State) Empty() bool {
	return !o.bold && !o.italic && !o.underline && !o.blink &&
		!o.fgcolor.HasColor() && !o.bgcolor.HasColor()
}

func (o *State) AnyChange(s *State) bool {
	return o.bold != s.bold ||
		o.italic != s.italic ||
		o.underline != s.underline ||
		o.blink != s.blink ||
		!o.fgcolor.Equal(&s.fgcolor) ||
		!o.bgcolor.Equal(&s.bgcolor)
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
	if s.blink && !o.blink {
		return true
	}
	if s.fgcolor.HasColor() && !o.fgcolor.HasColor() {
		return true
	}
	if s.bgcolor.HasColor() && !o.bgcolor.HasColor() {
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

	if o.fgcolor.HasColor() && !o.fgcolor.Equal(&s.fgcolor) {
		f(true, "color:"+o.fgcolor.String()+";")
	}

	if o.bgcolor.HasColor() && !o.bgcolor.Equal(&s.bgcolor) {
		f(true, "background-color:"+o.bgcolor.String()+";")
	}

	return err
}

func (o *State) WriteClasses(w io.Writer, s *State) error {
	var err error

	f := func(b bool, s string) {
		if err == nil && b {
			_, err = w.Write([]byte(s))
		}
	}

	f(o.blink, "aes2htm-blink")

	return err
}
