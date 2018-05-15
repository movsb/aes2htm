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
}

func (o *Aes2Htm) Input(r io.Reader) error {
	var er error
	var c byte

	var (
		bold      = false
		italic    = false
		underline = false
		fgcolor   = -1
		bgcolor   = -1

		oldBold      = false
		oldItalic    = false
		oldUnderline = false
		oldFgcolor   = -1
		oldBgcolor   = -1

		openTags = 0
	)

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
			oldBold = bold
			oldItalic = italic
			oldUnderline = underline
			oldFgcolor = fgcolor
			oldBgcolor = bgcolor

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
					for _, n := range ns {
						if n == 0 {
							bold = false
							italic = false
							underline = false
							fgcolor = -1
							bgcolor = -1
						} else if n == 1 {
							bold = true
						} else if n == 3 {
							italic = true
						} else if n == 4 {
							underline = true
						} else if 30 <= n && n <= 37 {
							fgcolor = n
						} else if 40 <= n && n <= 48 {
							bgcolor = n
						} else if n == 39 || n == 49 {
							// Default foreground color
							// Default background color
							fgcolor = -1
							bgcolor = -1
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
			if bold != oldBold || italic != oldItalic || underline != oldUnderline || fgcolor != oldFgcolor || bgcolor != oldBgcolor {
				// 如果原来有，现在没有。则应关闭重新打开，以去掉没有了的属性
				if oldBold || oldItalic || oldUnderline ||
					(oldFgcolor != -1 && oldFgcolor != fgcolor) ||
					(oldBgcolor != -1 && oldBgcolor != bgcolor) {
					fmt.Print(strings.Repeat("</span>", openTags))
					openTags = 0
				}
				// 重新输出新的属性
				if bold || italic || underline || fgcolor != -1 || bgcolor != -1 {
					fmt.Print("<span style=\"")
					if bold {
						fmt.Print("font-weight:bold;")
					}
					if italic {
						fmt.Print("font-style:italic;")
					}
					if underline {
						fmt.Print("text-decoration:underline;")
					}
					if fgcolor != oldFgcolor {
						switch fgcolor - 30 {
						case 0:
							fmt.Print("color:black;")
						case 1:
							fmt.Print("color:red;")
						case 2:
							fmt.Print("color:green;")
						case 3:
							fmt.Print("color:yellow;")
						case 4:
							fmt.Print("color:blue;")
						case 5:
							fmt.Print("color:magenta;")
						case 6:
							fmt.Print("color:cyan;")
						case 7:
							fmt.Print("color:white;")
						}
					}
					if bgcolor != oldBgcolor {
						switch bgcolor - 40 {
						case 0:
							fmt.Print("background-color:black;")
						case 1:
							fmt.Print("background-color:red;")
						case 2:
							fmt.Print("background-color:green;")
						case 3:
							fmt.Print("background-color:yellow;")
						case 4:
							fmt.Print("background-color:blue;")
						case 5:
							fmt.Print("background-color:magenta;")
						case 6:
							fmt.Print("background-color:cyan;")
						case 7:
							fmt.Print("background-color:white;")
						}
					}
					fmt.Print("\">")
					openTags++
				}
			}
		} else {
			if c < 128 {
				switch c {
				case '<':
					fmt.Print("&lt;")
				case '>':
					fmt.Print("&gt;")
				default:
					fmt.Printf("%c", c)
				}
				continue
			}

			br.UnreadByte()
			r, s, er := br.ReadRune()
			if r == utf8.RuneError || s < 1 || er != nil {
				return errors.New("invalid rune")
			}

			fmt.Printf("%c", r)
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
		err = ah.Input(f)
	} else {
		err = ah.Input(os.Stdin)
	}
	if err != nil {
		panic(err)
	}
}
