package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"time"
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

func (u *Unpacker) From6to16Terrain1() *image.NRGBA {
	canvas := image.NewNRGBA(image.Rect(0, 0, u.tileWidth*16, u.tileHeight*1))
	// todo optimize to generate automatically and consider scaling for 47 and 255 tilesets
	exp := export28TileSet()
	quadMap := append(exp[:14], exp[15], exp[14]) // 15 - filled terrain 2, 14 - small terrain 2 patch on top

	for idx := 0; idx < 16; idx++ {
		u.drawFullTile(canvas, quadMap[idx], idx, 16)
	}
	return canvas
}

func (u *Unpacker) From6to16Terrain2() *image.NRGBA {
	canvas := image.NewNRGBA(image.Rect(0, 0, u.tileWidth*16, u.tileHeight*1))
	// todo optimize to generate automatically and consider scaling for 47 and 255 tilesets
	exp := export28TileSet()
	quadMap := append(exp[14:], exp[1], exp[0]) // 0 - filled terrain 1, 1 - small terrain 1 patch on top

	for idx := 0; idx < 16; idx++ {
		u.drawFullTile(canvas, quadMap[idx], idx, 16)
	}
	return canvas
}

func (u *Unpacker) From6to28() *image.NRGBA {
	canvas := image.NewNRGBA(image.Rect(0, 0, u.tileWidth*14, u.tileHeight*2))
	// todo optimize to generate automatically and consider scaling for 47 and 255 tilesets

	quadMap := export28TileSet()

	for idx := 0; idx < 28; idx++ {
		u.drawFullTile(canvas, quadMap[idx], idx, 14)
	}
	return canvas
}

func (u *Unpacker) From28To92(img *image.NRGBA) *image.NRGBA {
	//51
	//6 x 9 region (54 total spaces, 3 empty)
	canvas := image.NewNRGBA(image.Rect(0, 0, u.tileWidth*8, u.tileHeight*12))
	rotations := [28]int{
		// we don't rotate 1st tile, and also the one with all 4 corners filled.
		// The opposing tiles need only 1 rotation
		0, 3, 3, 3, 1, 3, 3, 3, 3, 0, 3, 3, 1, 3,
		0, 3, 3, 3, 1, 3, 3, 3, 3, 0, 3, 3, 1, 3,
	}
	idx := 0 // result tile index
	i := 0   // rotations map index
	for y := 0; y < 2; y++ {
		for x := 0; x < 14; x++ {
			idx = u.expandWithRotation(img, canvas, idx, rotations[i], x, y, 8)
			i++
		}
	}
	return canvas
}

func export28TileSet() [28][4][2]int {
	return [28][4][2]int{
		// terrain 1
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
		// end of terrain 1
		// terrain 2
		{
			{0, 0}, {0, 1}, // tile 15
			{0, 1}, {1, 1},
		},
		{
			{0, 5}, {3, 5}, // tile 16
			{0, 1}, {1, 1},
		},
		{
			{0, 4}, {3, 0}, // tile 17
			{0, 5}, {1, 5},
		},
		{
			{0, 4}, {1, 4}, // tile 18
			{0, 5}, {1, 5},
		},
		{
			{0, 3}, {3, 3}, // tile 19
			{0, 4}, {3, 4},
		},
		{
			{0, 3}, {3, 0}, // tile 20
			{0, 4}, {3, 1},
		},
		{
			{0, 3}, {2, 3}, // tile 21
			{0, 4}, {3, 1},
		},
		{
			{0, 3}, {3, 0}, // tile 22
			{0, 4}, {2, 4},
		},
		{
			{0, 3}, {2, 3}, // tile 23
			{0, 4}, {2, 4},
		},
		{
			{2, 0}, {3, 0}, // tile 24
			{2, 1}, {3, 1},
		},
		{
			{2, 0}, {2, 3}, // tile 25
			{2, 1}, {3, 1},
		},
		{
			{2, 0}, {2, 3}, // tile 26
			{2, 1}, {2, 4},
		},
		{
			{2, 0}, {2, 3}, // tile 27
			{1, 4}, {3, 1},
		},
		{
			{2, 0}, {2, 3}, // tile 28
			{1, 4}, {2, 4},
		},
	}
}

func (u *Unpacker) drawFullTile(canvas *image.NRGBA, data [4][2]int, idx int, outXTiles int) {
	for i, xy := range data {
		x := xy[0]
		y := xy[1]
		line := idx / outXTiles
		row := idx % outXTiles
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

func (u *Unpacker) expandWithRotation(src, canvas *image.NRGBA, idx, rotations, x, y, outXTiles int) int {
	index := idx
	line := idx / outXTiles
	row := idx % outXTiles
	canvasMin := image.Point{
		X: row * u.tileWidth,
		Y: line * u.tileHeight,
	}
	canvasArea := image.Rectangle{
		Min: canvasMin,
		Max: canvasMin.Add(image.Point{
			X: u.tileWidth,
			Y: u.tileHeight,
		}),
	}
	point := image.Point{
		X: x * u.tileWidth,
		Y: y * u.tileHeight,
	}
	tile := src.SubImage(image.Rectangle{
		Min: point,
		Max: point.Add(image.Point{
			X: u.tileWidth,
			Y: u.tileHeight,
		}),
	}).(*image.NRGBA)
	draw.Draw(
		canvas,
		canvasArea,
		src,
		point,
		draw.Src,
	)
	//debugSave(canvas)

	for i := 0; i < rotations; i++ { // rotate 90 degrees for [rotations] times
		index++
		line := index / outXTiles
		row := index % outXTiles
		canvasMin := image.Point{
			X: row * u.tileWidth,
			Y: line * u.tileHeight,
		}
		canvasArea := image.Rectangle{
			Min: canvasMin,
			Max: canvasMin.Add(image.Point{
				X: u.tileWidth,
				Y: u.tileHeight,
			}),
		}
		tile = Rotate90(tile)
		draw.Draw(
			canvas,
			canvasArea,
			tile,
			image.Point{
				X: 0,
				Y: 0,
			},
			draw.Src,
		)
		//debugSave(canvas)
	}
	return index + 1 // because we draw original tile first
}

func debugSave(img *image.NRGBA) {
	filename := fmt.Sprintf("out/%d.png", time.Now().UnixNano())
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
