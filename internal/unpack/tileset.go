package unpack

import (
	"image"
	"image/draw"
)

type tileSet struct {
	canvas                *image.NRGBA
	tileWidth, tileHeight int
}

func newTileSet(canvas *image.NRGBA, tileWidth, tileHeight int) *tileSet {
	return &tileSet{
		canvas:     canvas,
		tileWidth:  tileWidth,
		tileHeight: tileHeight,
	}
}

func (t *tileSet) setTile(x, y int, tile *image.NRGBA) {
	canvasMin := image.Point{
		X: x * t.tileWidth,
		Y: y * t.tileHeight,
	}
	canvasArea := image.Rectangle{
		Min: canvasMin,
		Max: image.Point{
			X: canvasMin.X + t.tileWidth,
			Y: canvasMin.Y + t.tileHeight,
		},
	}
	draw.Draw(
		t.canvas,
		canvasArea,
		tile,
		image.Point{
			X: 0,
			Y: 0,
		},
		draw.Src,
	)
}

func (t *tileSet) setTileWithRotationLeft(x, y int, tile *image.NRGBA) *image.NRGBA {
	rotatedTile := rotateLeft90(tile)
	t.setTile(x, y, rotatedTile)
	return rotatedTile
}

func (t *tileSet) getCanvas() *image.NRGBA {
	return t.canvas
}
