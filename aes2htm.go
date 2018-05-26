// https://en.wikipedia.org/wiki/ANSI_escape_code

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

type Aes2Htm struct {
	w        io.Writer
	br       *bufio.Reader
	out      func(s string)
	st       State
	stb      State
	openTags int
}

func NewAes2Htm(w io.Writer) *Aes2Htm {
	ah := &Aes2Htm{}
	ah.w = w
	ah.out = func(s string) {
		ah.w.Write([]byte(s))
	}
	ah.st.Init()
	ah.stb.Init()
	return ah
}

func (o *Aes2Htm) inputPlain(c byte) error {
	if c < 128 {
		switch c {
		case '<':
			o.out("&lt;")
		case '>':
			o.out("&gt;")
		default:
			o.out(fmt.Sprintf("%c", c))
		}
		return nil
	}

	o.br.UnreadByte()
	r, s, er := o.br.ReadRune()
	if r == utf8.RuneError || s < 1 || er != nil {
		return errors.New("invalid rune")
	}

	o.out(fmt.Sprintf("%c", r))

	return nil
}

func (o *Aes2Htm) handleCSI() error {
	var er error
	var c byte

	// backup state
	o.stb = o.st

	st := &o.st

	hasNum := false
	n := 0
	ns := []int{}

	for {
		c, er = o.br.ReadByte()
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
			for i := 0; i < len(ns); i++ {
				n = ns[i]

				if n == 0 {
					st.Init()
				} else if n == 1 {
					st.bold = true
				} else if n == 3 {
					st.italic = true
				} else if n == 4 {
					st.underline = true
				} else if n == 5 || n == 6 {
					st.blink = true
				} else if 30 <= n && n <= 37 {
					st.fgcolor.SetIndex(n - 30)
				} else if n == 39 {
					st.fgcolor.SetNone()
				} else if 40 <= n && n <= 47 {
					st.bgcolor.SetIndex(n - 40)
				} else if n == 49 {
					st.bgcolor.SetNone()
				} else if n == 38 || n == 48 {
					// 5;n or 2;r;g;b
					if i++; i >= len(ns) {
						return fmt.Errorf("expect color")
					}
					switch ns[i] {
					case 5:
						if i++; i >= len(ns) {
							return fmt.Errorf("expect color")
						}
						index := ns[i]
						if n == 38 {
							st.fgcolor.SetIndex(index)
						} else if n == 48 {
							st.bgcolor.SetIndex(index)
						}
					case 2:
						if i+3 >= len(ns) {
							return fmt.Errorf("expect color")
						}
						r := ns[i+1]
						g := ns[i+2]
						b := ns[i+3]
						if n == 38 {
							st.fgcolor.SetRGB(r, g, b)
						} else if n == 48 {
							st.bgcolor.SetRGB(r, g, b)
						}
						i += 3
					}
				} else if 90 <= n && n <= 97 {
					st.fgcolor.SetIndex(n - 90 + 8)
				} else if 100 <= n && n <= 107 {
					st.bgcolor.SetIndex(n - 100 + 8)
				} else {
					return fmt.Errorf("invalid code: %d", n)
				}
			}
			break
		} else if c == '?' {
			// skip
			for {
				c, er := o.br.ReadByte()
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
		} else if c == 'r' { // Set Scrolling Region [top;bottom]
			break
		} else if c == 'l' { // Reset Mode
			break
		} else if c == 'H' { // Cursor Position
			break
		} else if c == 'J' { // Erase in Display
			break
		} else if c == 'd' { // Line Position Absolute
			break
		} else if c == 'X' { // Erase Characters
			break
		} else {
			s := fmt.Sprintf("invalid terminate: %c", c)
			return errors.New(s)
		}
	}

	// 如果发生了改变的话
	if st.AnyChange(&o.stb) {
		// 如果原来有，现在没有。则应关闭重新打开，以去掉没有了的属性
		if st.AnyClose(&o.stb) {
			o.out(strings.Repeat("</span>", o.openTags))
			o.openTags = 0
		}
		// 重新输出新的属性
		if !st.Empty() {
			o.out("<span")
			st.WriteStyles(o.w, &o.stb)
			st.WriteClasses(o.w, &o.stb)
			o.out(">")
			o.openTags++
		}
	}

	return nil
}

func (o *Aes2Htm) handleLeftParentheses() error {
	return nil
}

func (o *Aes2Htm) Input(r io.Reader) error {
	var er error
	var c byte

	o.openTags = 0

	o.st.Init()
	o.stb.Init()

	o.br = bufio.NewReader(r)
	br := o.br

	for {
		c, er = br.ReadByte()
		if er != nil {
			if er == io.EOF {
				return nil
			}
			return er
		}

		if c == '\033' {
			c, er = br.ReadByte()
			if er != nil {
				return er
			}
			switch c {
			case '[':
				if er = o.handleCSI(); er != nil {
					return er
				}
			case '(':
				if er = o.handleLeftParentheses(); er != nil {
					return er
				}
			case '=': // Application Keypad
				break
			default:
				return fmt.Errorf("unhandled char after escape: %c", c)
			}
		} else {
			if er = o.inputPlain(c); er != nil {
				return er
			}
		}
	}
}
