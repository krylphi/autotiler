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

	// todo parallel

	canvas := g.From6to16Terrain1()
	file15t1, err := os.Create(fmt.Sprintf("16x1_terrain1_%s", outputFile))
	if err != nil {
		panic(err)
	}
	defer file15t1.Close()

	err = png.Encode(file15t1, canvas)
	if err != nil {
		panic(err)
	}

	canvas = g.From6to16Terrain2()
	file15t2, err := os.Create(fmt.Sprintf("16x1_terrain2_%s", outputFile))
	if err != nil {
		panic(err)
	}
	defer file15t2.Close()

	err = png.Encode(file15t2, canvas)
	if err != nil {
		panic(err)
	}

	canvas = g.From6to28()
	file1, err := os.Create(fmt.Sprintf("14x2_%s", outputFile))
	if err != nil {
		panic(err)
	}
	defer file1.Close()

	err = png.Encode(file1, canvas)
	if err != nil {
		panic(err)
	}

	/*
		// TODO: consider usability
		canvas = g.From28To92(canvas)
		file2, err := os.Create(fmt.Sprintf("8x12_%s", outputFile))
		if err != nil {
			panic(err)
		}
		defer file2.Close()

		err = png.Encode(file2, canvas)
		if err != nil {
			panic(err)
		}
	*/
}
