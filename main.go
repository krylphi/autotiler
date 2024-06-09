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
	"slices"
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
	missingTerrain2 = "missing-terrain-two"
)

const (
	export16  = "16"
	export28  = "28"
	export48  = "48"
	exportAll = "all"
)

type param struct {
	value   string
	options []string
}

func main() {
	args, err := parseArgs()
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	printArgs(args)

	inFiles, ok := args[inKey]
	if !ok {
		log.Print("Missing input file")
		os.Exit(1)
	}
	outFiles, outs := args[outKey]
	for i, inFile := range inFiles {
		outputFile := fmt.Sprintf("%d.local.png", i)
		if outs {
			if len(outFiles) > i {
				outputFile = outFiles[i].value
			}
		}
		err := export(args, inFile, outputFile)
		if err != nil {
			log.Print(err)
			continue
		}
	}
}

func export(args map[string][]param, inputFile param, outputFile string) error { // todo wrap errors
	imgFile, err := os.Open(inputFile.value)
	if err != nil {
		return err
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return err
	}

	padding := 0
	if paddings, ok := args[paddingKey]; ok {
		padding, err = strconv.Atoi(paddings[0].value)
		if err != nil {
			return err
		}
	}

	options := unpack.Options{
		Padding:           padding,
		MissingTerrainTwo: slices.Contains(inputFile.options, missingTerrain2),
	}

	unpacker := unpack.NewUnpacker(img, 2, 3, options)
	if err := unpacker.Init(2); err != nil {
		return err
	}

	exports, ok := args[exportKey]
	var exportTypes []string
	if !ok || len(exports) == 0 || exports[0].value == exportAll {
		exportTypes = []string{export16, export28, export48}
	} else {
		exportTypes = make([]string, len(exports))
		for e := range exports {
			exportTypes[e] = exports[e].value
		}
	}

	for e := range exportTypes {
		exportType := exportTypes[e]
		switch exportType {
		case export16:
			err := produceTileset(unpacker.From6to16Terrain1, outputFile, "16x1_terrain1")
			if err != nil {
				return err
			}
			err = produceTileset(unpacker.From6to16Terrain2, outputFile, "16x1_terrain2")
			if err != nil {
				return err
			}
		case export28:
			err := produceTileset(unpacker.From6to28, outputFile, "14x2")
			if err != nil {
				return err
			}
		case export48:
			err := produceTileset(unpacker.From6to48Terrain1, outputFile, "12x4_terrain1")
			if err != nil {
				return err
			}
			err = produceTileset(unpacker.From6to48Terrain2, outputFile, "12x4_terrain2")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func parseArgs() (map[string][]param, error) {
	if len(os.Args) < 2 {
		log.Print(
			"Usage: autotiler -in <file_in> [-o <file_out>] [-p <padding>] [-e <export_type(16,28,48,all)>] [--missing-terrain-two]\n" +
				"       -in, -o and -e can be repeated\n")
		os.Exit(1)
	}
	args := make(map[string][]param)
	allTilesets := false
	for i := 1; i < len(os.Args); i += 2 {
		if strings.HasPrefix(os.Args[i], "--") || !strings.HasPrefix(os.Args[i], "-") {
			return nil, fmt.Errorf("Invalid argument: %s", os.Args[i])
		}
		key := strings.TrimPrefix(os.Args[i], "-")
		v, ok := args[key]
		value := param{
			value:   os.Args[i+1],
			options: []string{},
		}
		if key == exportKey && allTilesets {
			continue
		}
		if key == exportKey && strings.EqualFold(value.value, exportAll) {
			allTilesets = true
			v = []param{{
				value: "all",
			}}
			args[key] = v
			continue
		}
		if len(os.Args) > i+2 {
			if strings.HasPrefix(os.Args[i+2], "--") {
				value.options = append(value.options, os.Args[i+2])
				i++
			}
			for {
				if len(os.Args) <= i || !strings.HasPrefix(os.Args[i], "--") {
					break
				}
				if strings.HasPrefix(os.Args[i], "--") {
					value.options = append(value.options, os.Args[i+2])
					i++
				}
			}
		}
		if ok {
			v = append(v, value)
		} else {
			v = []param{value}
		}
		args[key] = v
	}
	return args, nil
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
func printArgs(args map[string][]param) {
	log.Printf("Args:\n")
	for key, arg := range args {
		for _, param := range arg {
			log.Printf("%s: %s\n", key, param)
		}
	}
}
