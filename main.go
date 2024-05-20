package main

import (
	"image"
	"image/png"
	"log"
	"os"
)

func main() {

	inputFile := os.Args[1]

	outputFile := os.Args[2]

	imgFile, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}

	g := NewUnpacker(img, 15)
	if err := g.Init(); err != nil {
		log.Fatal(err)
	}

	canvas := g.From6to15()
	oFile, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer oFile.Close()

	err = png.Encode(oFile, canvas)
	if err != nil {
		panic(err)
	}
}
