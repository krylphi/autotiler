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

type anchorSet [][]image.Point

type quadTileData *[4][2]int

type Unpacker struct {
	anchors               anchorSet
	tileWidth, tileHeight int
	src                   image.Image
	xTiles                int
	yTiles                int
}

func NewUnpacker(src image.Image, xTiles, yTiles int) *Unpacker {
	// todo auto detect
	tileWidth := src.Bounds().Dx() / xTiles
	tileHeight := src.Bounds().Dy() / yTiles

	return &Unpacker{
		src:        src,
		tileWidth:  tileWidth,
		tileHeight: tileHeight,
		xTiles:     xTiles,
		yTiles:     yTiles,
	}
}

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
func (u *Unpacker) drawFullTile(canvas *image.NRGBA, data quadTileData, idx, outXTiles int) {
	if data == nil {
		return
	}
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

func (u *Unpacker) drawFullSingleTile(tile *image.NRGBA, data quadTileData) {
	u.drawFullTile(tile, data, 0, 1)
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
