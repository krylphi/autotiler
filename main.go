package main

import (
	"fmt"
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

	g := NewUnpacker(img, 4)
	if err := g.Init(); err != nil {
		log.Fatal(err)
	}

	canvas := g.From6to16()
	file4x4, err := os.Create(fmt.Sprintf("4x4_%s", outputFile))
	if err != nil {
		panic(err)
	}
	defer file4x4.Close()

	err = png.Encode(file4x4, canvas)
	if err != nil {
		panic(err)
	}

	canvas = g.From16to51(canvas)
	file6x9, err := os.Create(fmt.Sprintf("6x9_%s", outputFile))
	if err != nil {
		panic(err)
	}
	defer file6x9.Close()

	err = png.Encode(file6x9, canvas)
	if err != nil {
		panic(err)
	}
}
