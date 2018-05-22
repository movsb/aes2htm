// https://en.wikipedia.org/wiki/ANSI_escape_code

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
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
		switch o.fgcolor - 30 {
		case 0:
			f(true, "color:black;")
		case 1:
			f(true, "color:red;")
		case 2:
			f(true, "color:green;")
		case 3:
			f(true, "color:yellow;")
		case 4:
			f(true, "color:blue;")
		case 5:
			f(true, "color:magenta;")
		case 6:
			f(true, "color:cyan;")
		case 7:
			f(true, "color:white;")
		default:
			panic("bad fg color: " + fmt.Sprint(o.fgcolor))
		}
	}

	if o.bgcolor != -1 && o.bgcolor != s.bgcolor {
		switch o.bgcolor - 40 {
		case 0:
			f(true, "background-color:black;")
		case 1:
			f(true, "background-color:red;")
		case 2:
			f(true, "background-color:green;")
		case 3:
			f(true, "background-color:yellow;")
		case 4:
			f(true, "background-color:blue;")
		case 5:
			f(true, "background-color:magenta;")
		case 6:
			f(true, "background-color:cyan;")
		case 7:
			panic("bad bg color: " + fmt.Sprint(o.bgcolor))
		}
	}

	return err
}

type Aes2Htm struct {
}

func (o *Aes2Htm) Input(w io.Writer, r io.Reader) error {
	var er error
	var c byte

	var openTags = 0
	var st State
	var stb State

	st.Init()
	stb.Init()

	out := func(s string) {
		w.Write([]byte(s))
	}

	br := bufio.NewReader(r)
	for {
		c, er = br.ReadByte()
		if er != nil {
			if er == io.EOF {
				return nil
			}
			return er
		}

		if c == '\033' {
			// backup state
			stb = st

			c, er = br.ReadByte()
			if er != nil {
				return er
			}
			hasNum := false
			n := 0
			ns := []int{}
			for {
				c, er = br.ReadByte()
				if er != nil {
					return er
				}
				if c >= '0' && c <= '9' {
					if !hasNum {
						hasNum = true
						n = 0
					}
					n *= 10
					n += int(c - '0')
					continue
				} else if c == ';' {
					if hasNum {
						ns = append(ns, n)
						hasNum = false
					}
					continue
				} else if c == 'm' {
					if hasNum {
						ns = append(ns, n)
						hasNum = false
					}
					if len(ns) == 0 {
						ns = append(ns, 0)
					}
					for _, n := range ns {
						if n == 0 {
							st.Init()
						} else if n == 1 {
							st.bold = true
						} else if n == 3 {
							st.italic = true
						} else if n == 4 {
							st.underline = true
						} else if 30 <= n && n <= 37 {
							st.fgcolor = n
						} else if 40 <= n && n <= 47 {
							st.bgcolor = n
						} else if n == 39 || n == 49 {
							// Default foreground color
							// Default background color
							st.fgcolor = -1
							st.bgcolor = -1
						} else {
							return errors.New(fmt.Sprintf("invalid code: %d", n))
						}
					}
					break
				} else if c == '?' {
					// skip
					for {
						c, er := br.ReadByte()
						if er != nil {
							return errors.New("invalid")
						}
						if c >= '0' && c <= '9' {
							continue
						} else {
							break
						}
					}

					break
				} else {
					s := fmt.Sprintf("invalid terminate: %c", c)
					return errors.New(s)
				}
			}

			// 如果发生了改变的话
			if st.AnyChange(&stb) {
				// 如果原来有，现在没有。则应关闭重新打开，以去掉没有了的属性
				if st.AnyClose(&stb) {
					out(strings.Repeat("</span>", openTags))
					openTags = 0
				}
				// 重新输出新的属性
				if !st.Empty() {
					out("<span style=\"")
					st.WriteStyles(w, &stb)
					out("\">")
					openTags++
				}
			}
		} else {
			if c < 128 {
				switch c {
				case '<':
					out("&lt;")
				case '>':
					out("&gt;")
				default:
					out(fmt.Sprintf("%c", c))
				}
				continue
			}

			br.UnreadByte()
			r, s, er := br.ReadRune()
			if r == utf8.RuneError || s < 1 || er != nil {
				return errors.New("invalid rune")
			}

			out(fmt.Sprintf("%c", r))
		}
	}

	return nil
}

func main() {
	var err error
	ah := &Aes2Htm{}
	if len(os.Args) >= 2 {
		var f *os.File
		f, err = os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer f.Close()
		err = ah.Input(os.Stdout, f)
	} else {
		err = ah.Input(os.Stdout, os.Stdin)
	}
	if err != nil {
		panic(err)
	}
}
