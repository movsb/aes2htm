package pkg

import "fmt"

// https://en.wikipedia.org/wiki/ANSI_escape_code#3/4_bit

var Palette = [256]string{
	// 0-7 standard colors (ubuntu version)
	`rgb(1,1,1)`, `rgb(222,56,43)`, `rgb(57,181,74)`, `rgb(255,199,6)`, `rgb(0,111,184)`, `rgb(118,38,113)`, `rgb(44,181,233)`, `rgb(204,204,204)`,

	// 8-15 high-intensity colors (ubuntu version)
	`rgb(128,128,128)`, `rgb(255,0,0)`, `rgb(0,255,0)`, `rgb(255,255,0)`, `rgb(0,0,255)`, `rgb(255,0,255)`, `rgb(0,255,255)`, `rgb(255,255,255)`,
}

func init() {
	// 16-231
	// 总共216种颜色：6*6*6
	// 每种颜色范围：0 ≤ r,g,b ≤ 5
	// 索引算法：16 + 36*r + 6*g + b
	// 生成算法：红色每36次变化一次，绿色每6次变化一次，蓝色每次变化一次
	// 颜色递进：256种颜色被分成6次递进：{0x00, 0x57, 0x87, 0xaf, 0xd7, 0xff}
	var steps1 = [6]uint32{0x00, 0x57, 0x87, 0xaf, 0xd7, 0xff}
	for r := 0; r < 6; r++ {
		for g := 0; g < 6; g++ {
			for b := 0; b < 6; b++ {
				index := 16 + 36*r + 6*g + b
				color := steps1[r]<<16 + steps1[g]<<8 + steps1[b]
				str := fmt.Sprintf("#%06x", color)
				Palette[index] = str
			}
		}
	}

	// 232-255
	// 24阶灰阶颜色：00x ~ 23x
	// 从8开始，即 x=8
	for i := 0; i < 24; i++ {
		index := 232 + i
		color := uint32(8 + i*10)
		color = color<<16 + color<<8 + color
		Palette[index] = fmt.Sprintf("#%06x", color)
	}
}
