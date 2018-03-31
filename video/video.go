package main

import (
	"fmt"
	"image/png"
	"io"
	"log"
	"os"
)

func main() {
	for true {
		img, err := png.Decode(os.Stdin)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(img.At(0, 0))
	}

}
