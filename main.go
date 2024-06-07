/*
 * MIT License
 *
 * Copyright (c) 2024 The autotiler authors
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"strconv"

	"github.com/krylphi/autotiler/internal/unpack"
)

func main() {

	inputFile := os.Args[1]

	outputFile := os.Args[2]
	var err error
	var padding = 0
	if len(os.Args) > 3 {
		padding, err = strconv.Atoi(os.Args[3])
		if err != nil {
			panic(err)
		}
	}

	imgFile, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}

	g := unpack.NewUnpacker(img, 2, 3, padding)
	if err := g.Init(2); err != nil {
		panic(err)
	}

	// todo cleanup and parallel

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
