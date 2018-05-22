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

type Aes2Htm struct {
	w   io.Writer
	br  *bufio.Reader
	out func(s string)
}

func NewAes2Htm(w io.Writer, r io.Reader) *Aes2Htm {
	ah := &Aes2Htm{}
	ah.w = w
	ah.br = bufio.NewReader(r)
	ah.out = func(s string) {
		ah.w.Write([]byte(s))
	}
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

func (o *Aes2Htm) Input(w io.Writer, r io.Reader) error {
	var er error
	var c byte

	var openTags = 0
	var st State
	var stb State

	st.Init()
	stb.Init()

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
					o.out(strings.Repeat("</span>", openTags))
					openTags = 0
				}
				// 重新输出新的属性
				if !st.Empty() {
					o.out("<span style=\"")
					st.WriteStyles(w, &stb)
					o.out("\">")
					openTags++
				}
			}
		} else {
			if er = o.inputPlain(c); er != nil {
				return er
			}
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
