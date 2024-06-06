package main

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/krylphi/autotiler/internal/unpack"
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

	g := unpack.NewUnpacker(img, 2, 3)
	if err := g.Init(2); err != nil {
		panic(err)
	}

	// todo parallel

	canvas, err := g.From6to48Terrain1()
	if err != nil {
		panic(err)
	}
	file48t1, err := os.Create(fmt.Sprintf("out/12x4_terrain1_%s", outputFile))
	if err != nil {
		panic(err)
	}
	defer file48t1.Close()

	err = png.Encode(file48t1, canvas)
	if err != nil {
		panic(err)
	}

	canvas, err = g.From6to48Terrain2()
	if err != nil {
		panic(err)
	}
	file48t2, err := os.Create(fmt.Sprintf("out/12x4_terrain2_%s", outputFile))
	if err != nil {
		panic(err)
	}
	defer file48t2.Close()

	err = png.Encode(file48t2, canvas)
	if err != nil {
		panic(err)
	}

	canvas, err = g.From6to16Terrain1()
	if err != nil {
		panic(err)
	}
	file15t1, err := os.Create(fmt.Sprintf("out/16x1_terrain1_%s", outputFile))
	if err != nil {
		panic(err)
	}
	defer file15t1.Close()

	err = png.Encode(file15t1, canvas)
	if err != nil {
		panic(err)
	}

	canvas, err = g.From6to16Terrain2()
	if err != nil {
		panic(err)
	}
	file15t2, err := os.Create(fmt.Sprintf("out/16x1_terrain2_%s", outputFile))
	if err != nil {
		panic(err)
	}
	defer file15t2.Close()

	err = png.Encode(file15t2, canvas)
	if err != nil {
		panic(err)
	}

	canvas, err = g.From6to28()
	if err != nil {
		panic(err)
	}
	file1, err := os.Create(fmt.Sprintf("out/14x2_%s", outputFile))
	if err != nil {
		panic(err)
	}
	defer file1.Close()

	err = png.Encode(file1, canvas)
	if err != nil {
		panic(err)
	}
}
