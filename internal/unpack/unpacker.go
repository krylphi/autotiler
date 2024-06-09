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

package unpack

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"time"
)

var (
	errInvalidPackType = errors.New("invalid pack type")
)

// anchorSet represents a set of anchor points for a tile set.
// The anchor point is used to determine the starting point for drawing a tile on the canvas.
// The image will be inserted down and right from the anchor point.
type anchorSet [][]image.Point

// quadTileData represents a quad map for a 2x3 tile set.
// Basically, every tile of the original 2x3 tile set is split into 4 smaller tiles.
// Different combinations of said sub tiles will produce resulting tile patterns.
// Every sub tile of original tile set is represented by pair of their coordinates.
// So for example filled terrain 2 tile will be:
//
//	&[4][2]int{
//		{0, 0}, {0, 1},
//		{1, 0}, {1, 1},
//	}
//	 where {x, y} is a coordinate of a sub tile.
type quadTileData *[4][2]int

type Unpacker struct {
	anchors               anchorSet
	tileWidth, tileHeight int
	src                   image.Image
	xTiles                int
	yTiles                int
	padding               int
	options               Options
}

type Options struct {
	Padding           int
	MissingTerrainTwo bool
}

func NewUnpacker(src image.Image, xTiles, yTiles int, options Options) *Unpacker {
	// todo auto detect
	tileWidth := src.Bounds().Dx() / xTiles
	tileHeight := src.Bounds().Dy() / yTiles

	return &Unpacker{
		src:        src,
		tileWidth:  tileWidth,
		tileHeight: tileHeight,
		xTiles:     xTiles,
		yTiles:     yTiles,
		options:    options,
	}
}

// getAnchorPoint calculates the anchor point for a specific tile position and tile side segments.
// The anchor point is used to determine the starting point for drawing a tile on the canvas.
//
// Parameters:
// - x: The x-coordinate of the tile position.
// - y: The y-coordinate of the tile position.
// - tileSideSegments: The number of segments in a tile (e.g. packed tile was segmented. 2 for to 2x2, 3 for 3x3, etc.).
//
// Returns:
// - An image.Point representing the anchor point (top-left).
func (u *Unpacker) getAnchorPoint(x, y, tileSideSegments int) image.Point {
	anchor := image.Point{
		X: x * u.tileWidth / tileSideSegments,
		Y: y * u.tileHeight / tileSideSegments,
	}
	return anchor
}

func (u *Unpacker) Init(tileSideSegments int) error {
	xCnt := u.xTiles * tileSideSegments
	yCnt := u.yTiles * tileSideSegments
	anchors := make([][]image.Point, xCnt)
	for x := 0; x < xCnt; x++ {
		yAreas := make([]image.Point, yCnt)
		anchors[x] = yAreas
		for y := 0; y < yCnt; y++ {
			anchors[x][y] = u.getAnchorPoint(x, y, tileSideSegments)
		}
	}
	u.anchors = anchors
	return nil
}

// outXTiles is the number of output tiles in x direction
// tileSideSegments is the number of segments in a tile (e.g. packed tile was segmented. 2 for to 2x2, 3 for 3x3, etc.)
// drawFullTile draws a full tile on the canvas based on the provided data and index.
// It calculates the position of the tile on the canvas and draws the corresponding segments from the source image.
//
// Parameters:
// - canvas: The image.NRGBA on which the tile will be drawn.
// - data: The quadTileData containing the positions of the segments in the source image.
// - idx: The index of the tile to be drawn.
// - outXTiles: The number of output tiles in the x direction.
//
// Returns:
// - Nothing.
func (u *Unpacker) drawFullTile(canvas *image.NRGBA, data quadTileData, idx, outXTiles int) {
	if data == nil {
		return
	}
	for i, xy := range data {
		x := xy[0]
		y := xy[1]
		line := idx / outXTiles
		row := idx % outXTiles
		paddingY := u.options.Padding + line*2*u.options.Padding
		paddingX := u.options.Padding + row*2*u.options.Padding

		shiftX := i%2*u.tileWidth/2 + paddingX
		shiftY := i>>1*u.tileHeight/2 + paddingY
		canvasMin := image.Point{
			X: row*u.tileWidth + shiftX,
			Y: line*u.tileHeight + shiftY,
		}
		canvasArea := image.Rectangle{
			Min: canvasMin,
			Max: image.Point{
				X: canvasMin.X + u.tileWidth/2,
				Y: canvasMin.Y + u.tileHeight/2,
			},
		}
		point := u.anchors[x][y]
		draw.Draw(
			canvas,
			canvasArea,
			u.src,
			point,
			draw.Src,
		)
	}
}

// drawFullSingleTile used for rendering a single tile.
func (u *Unpacker) drawFullSingleTile(tile *image.NRGBA, data quadTileData) {
	u.drawFullTile(tile, data, 0, 1)
}

func (u *Unpacker) paddedTileWidth() int {
	return u.tileWidth + u.options.Padding*2
}

func (u *Unpacker) paddedTileHeight() int {
	return u.tileHeight + u.options.Padding*2
}

//nolint:unused //debug function
func debugSave(img *image.NRGBA) {
	filename := fmt.Sprintf("out/debug_%d.png", time.Now().UnixNano())
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
}
