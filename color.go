package main

// https://en.wikipedia.org/wiki/ANSI_escape_code#3/4_bit

var Palette = [256]string{
	// standard colors (ubuntu version)
	`rgb(1,1,1)`, `rgb(222,56,43)`, `rgb(57,181,74)`, `rgb(255,199,6)`, `rgb(0,111,184)`, `rgb(118,38,113)`, `rgb(44,181,233)`, `rgb(204,204,204)`,

	// high-intensity colors (ubuntu version)
	`rgb(128,128,128)`, `rgb(255,0,0)`, `rgb(0,255,0)`, `rgb(255,255,0)`, `rgb(0,0,255)`, `rgb(255,0,255)`, `rgb(0,255,255)`, `rgb(255,255,255)`,
}
