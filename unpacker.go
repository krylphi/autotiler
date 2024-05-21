package main

import (
	"image"
	"image/draw"
)

const (
	xTiles = 2
	yTiles = 3
)

type anchorSet [][]image.Point

type Unpacker struct {
	anchors               anchorSet
	tileWidth, tileHeight int
	src                   image.Image
	outXTiles             int
}

func NewUnpacker(src image.Image, outXTiles int) *Unpacker {

	// todo auto detect

	tileWidth := src.Bounds().Dx() / xTiles
	tileHeight := src.Bounds().Dy() / yTiles

	return &Unpacker{
		src:        src,
		tileWidth:  tileWidth,
		tileHeight: tileHeight,
		outXTiles:  outXTiles,
	}
}

func (u *Unpacker) getAnchorPoint(x, y int) image.Point {
	anchor := image.Point{
		X: x * u.tileWidth / 2,
		Y: y * u.tileHeight / 2,
	}
	return anchor
}

func (u *Unpacker) Init() error {
	xCnt := xTiles * 2
	yCnt := yTiles * 2
	anchors := make([][]image.Point, xCnt)
	for x := 0; x < xCnt; x++ {
		yAreas := make([]image.Point, yCnt)
		anchors[x] = yAreas
		for y := 0; y < yCnt; y++ {
			anchors[x][y] = u.getAnchorPoint(x, y)
		}
	}
	u.anchors = anchors
	return nil
}

func (u *Unpacker) From6to16() image.Image {
	canvas := image.NewNRGBA(image.Rect(0, 0, u.tileWidth*4, u.tileHeight*4))
	// todo optimize to generate automatically and consider scaling for 47 and 255 tilesets
	quadMap := [16][4][2]int{
		{
			{1, 3}, {2, 3}, // tile 1
			{1, 4}, {2, 4},
		},
		{
			{3, 0}, {2, 0}, // tile 2
			{1, 4}, {2, 4},
		},
		{
			{3, 3}, {0, 5}, // tile 3
			{3, 0}, {2, 2},
		},
		{
			{3, 3}, {0, 0}, // tile 4
			{3, 0}, {2, 2},
		},
		{
			{3, 3}, {0, 3}, // tile 5
			{3, 4}, {0, 4},
		},
		{
			{3, 3}, {0, 5}, // tile 6
			{3, 4}, {0, 2},
		},
		{
			{3, 3}, {0, 0}, // tile 7
			{3, 4}, {0, 2},
		},
		{
			{3, 3}, {0, 5}, // tile 8
			{3, 4}, {1, 1},
		},
		{
			{3, 3}, {0, 0}, // tile 9
			{3, 4}, {1, 1},
		},
		{
			{3, 5}, {0, 5}, // tile 10
			{3, 2}, {0, 2},
		},
		{
			{3, 5}, {1, 0}, // tile 11
			{3, 2}, {0, 2},
		},
		{
			{3, 5}, {1, 0}, // tile 12
			{3, 2}, {1, 1},
		},
		{
			{3, 5}, {1, 0}, // tile 13
			{0, 1}, {0, 2},
		},
		{
			{3, 5}, {1, 0}, // tile 14
			{0, 1}, {1, 1},
		},
		{
			{0, 0}, {1, 0}, // tile 15
			{0, 2}, {3, 2},
		},
		{
			{0, 0}, {0, 1}, // tile 16
			{0, 1}, {1, 1},
		},
	}

	for idx := 0; idx < 16; idx++ {
		u.drawFullTile(canvas, quadMap[idx], idx)
	}
	return canvas
}

func (u *Unpacker) drawFullTile(canvas *image.NRGBA, data [4][2]int, idx int) {
	for i, xy := range data {
		x := xy[0]
		y := xy[1]
		line := idx / u.outXTiles
		row := idx % u.outXTiles
		shiftX := i % 2 * u.tileWidth / 2
		shiftY := i >> 1 * u.tileHeight / 2
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
