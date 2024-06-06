package unpack

import (
	"image"
	"image/draw"
)

type TileSet struct {
	canvas                *image.NRGBA
	tileWidth, tileHeight int
}

func NewTileSet(canvas *image.NRGBA, tileWidth, tileHeight int) *TileSet {
	return &TileSet{
		canvas:     canvas,
		tileWidth:  tileWidth,
		tileHeight: tileHeight,
	}
}

func (t *TileSet) GetTile(x, y int) *image.NRGBA {
	return t.canvas.SubImage(image.Rect(x*t.tileWidth, y*t.tileHeight, (x+1)*t.tileWidth, (y+1)*t.tileHeight)).(*image.NRGBA)
}

func (t *TileSet) SetTile(x, y int, tile *image.NRGBA) {
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

func (t *TileSet) SetTileWithRotationLeft(x, y int, tile *image.NRGBA) *image.NRGBA {
	rotatedTile := RotateLeft90(tile)
	t.SetTile(x, y, rotatedTile)
	return rotatedTile
}

func (t *TileSet) GetCanvas() *image.NRGBA {
	return t.canvas
}
