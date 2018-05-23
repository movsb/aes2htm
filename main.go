package main

import (
	"os"
)

func main() {
	var err error
	ah := NewAes2Htm(os.Stdout)
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
