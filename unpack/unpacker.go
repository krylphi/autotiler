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

func newTileData(
	top,
	bottom [2][2]int) quadTileData {
	return &[4][2]int{
		top[0], top[1],
		bottom[0], bottom[1],
	}
}

func (u *Unpacker) getAnchorPoint(x, y int, tileSideSegments int) image.Point {
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
func (u *Unpacker) drawFullTile(canvas *image.NRGBA, data quadTileData, idx int, outXTiles int, tileSideSegments int) {
	if data == nil {
		return
	}
	for i, xy := range data {
		x := xy[0]
		y := xy[1]
		line := idx / outXTiles
		row := idx % outXTiles
		shiftX := i % tileSideSegments * u.tileWidth / tileSideSegments
		shiftY := i >> 1 * u.tileHeight / tileSideSegments
		canvasMin := image.Point{
			X: row*u.tileWidth + shiftX,
			Y: line*u.tileHeight + shiftY,
		}
		canvasArea := image.Rectangle{
			Min: canvasMin,
			Max: image.Point{
				X: canvasMin.X + u.tileWidth/tileSideSegments,
				Y: canvasMin.Y + u.tileHeight/tileSideSegments,
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
	u.drawFullTile(tile, data, 0, 1, 2)
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
		tile = RotateLeft90(tile)
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

func (u *Unpacker) isVertical() bool {
	return u.yTiles > u.xTiles
}

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
