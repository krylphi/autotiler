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
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/krylphi/autotiler/internal/unpack"
)

const (
	inKey      = "in"
	paddingKey = "p"
	exportKey  = "e"
	outKey     = "o"
)

const (
	export16  = "16"
	export28  = "28"
	export48  = "48"
	exportAll = "all"
)

func main() {
	args := parseArgs()
	inFiles, ok := args[inKey]
	if !ok {
		log.Print("Missing input file")
		os.Exit(1)
	}
	inputFile := inFiles[0] // todo add handling for multiple inputs
	imgFile, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}

	padding, err := strconv.Atoi(args[paddingKey][0])
	if err != nil {
		panic(err)
	}

	unpacker := unpack.NewUnpacker(img, 2, 3, padding)
	if err := unpacker.Init(2); err != nil {
		panic(err)
	}

	exports, ok := args[exportKey]
	var exportTypes []string
	if !ok || len(exports) == 0 || exports[0] == exportAll {
		exportTypes = []string{export16, export28, export48}
	} else {
		exportTypes = exports
	}

	outputFile := "out.local.png"
	outFiles, ok := args[outKey]
	if ok {
		outputFile = outFiles[0]
	}

	for e := range exportTypes {
		exportType := exportTypes[e]
		switch exportType {
		case export16:
			err := produceTileset(unpacker.From6to16Terrain1, outputFile, "16x1_terrain1")
			if err != nil {
				panic(err)
			}
			err = produceTileset(unpacker.From6to16Terrain2, outputFile, "16x1_terrain2")
			if err != nil {
				panic(err)
			}
		case export28:
			err := produceTileset(unpacker.From6to28, outputFile, "14x2")
			if err != nil {
				panic(err)
			}
		case export48:
			err := produceTileset(unpacker.From6to48Terrain1, outputFile, "12x4_terrain1")
			if err != nil {
				panic(err)
			}
			err = produceTileset(unpacker.From6to48Terrain2, outputFile, "12x4_terrain2")
			if err != nil {
				panic(err)
			}
		}
	}
}

func parseArgs() map[string][]string {
	if len(os.Args) < 2 {
		log.Print(
			"Usage: autotiler -in <file_in> [-o <file_out>] [-p <padding>] [-e <export_type(16,28,48,all)>]\n" +
				"       -e can be repeated\n")
		os.Exit(1)
	}
	res := make(map[string][]string)
	allTilesets := false
	for i := 1; i < len(os.Args); i += 2 {
		key := strings.TrimPrefix(os.Args[i], "-")
		v, ok := res[key]
		value := os.Args[i+1]
		if key == exportKey && allTilesets {
			continue
		}
		if key == exportKey && strings.EqualFold(value, exportAll) {
			allTilesets = true
			v = []string{"all"}
			res[key] = v
			continue
		}
		if ok {
			v = append(v, value)
		} else {
			v = []string{value}
		}
		res[key] = v
	}
	return res
}

func produceTileset(unpackFunction func() (*image.NRGBA, error), outputPath, exportType string) error {
	canvas, err := unpackFunction()
	if err != nil {
		return err
	}
	cleanPath := filepath.Clean(outputPath)
	outputFile := fmt.Sprintf("%s_%s", exportType, filepath.Base(cleanPath))
	outputFile = filepath.Join(filepath.Dir(cleanPath), outputFile)
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)
	err = png.Encode(file, canvas)
	if err != nil {
		return err
	}
	return nil
}

//nolint:unused //debug function
func printArgs(args map[string][]string) {
	for key, params := range args {
		for _, param := range params {
			log.Printf("%s: %s\n", key, param)
		}
	}
}
