package main

import (
	"io"
	"log"
	"os"

	"github.com/romiras/aes2htm/pkg"
)

func main() {
	tohtml := len(os.Args) == 2 && os.Args[1] == "--html"
	if tohtml {
		ah, err := pkg.NewAes2Htm(os.Stdout)
		if err != nil {
			log.Fatalln(err)
		}

		err = ah.WriteHTML(os.Stdin)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		_, err := io.Copy(os.Stdout, os.Stdin)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
