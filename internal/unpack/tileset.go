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

// setTile draws the given tile onto the canvas at the specified coordinates (x, y).
// The tile is drawn at its original orientation.
//
// Parameters:
// - x: The x-coordinate on the canvas where the tile will be placed.
// - y: The y-coordinate on the canvas where the tile will be placed.
// - tile: The image.NRGBA tile to be drawn onto the canvas.
func (t *tileSet) setTile(x, y int, tile *image.NRGBA) {
	// Calculate the minimum point of the area on the canvas where the tile will be drawn.
	canvasMin := image.Point{
		X: x * t.tileWidth,
		Y: y * t.tileHeight,
	}

	// Define the area on the canvas where the tile will be drawn.
	canvasArea := image.Rectangle{
		Min: canvasMin,
		Max: image.Point{
			X: canvasMin.X + t.tileWidth,
			Y: canvasMin.Y + t.tileHeight,
		},
	}

	// Draw the tile onto the canvas at the specified area.
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

// setTileWithRotationLeft draws the given tile onto the canvas at the specified coordinates (x, y)
// and rotates the tile 90 degrees to the left before drawing it.
// The tile is drawn at its rotated orientation.
//
// Parameters:
// - x: The x-coordinate on the canvas where the tile will be placed.
// - y: The y-coordinate on the canvas where the tile will be placed.
// - tile: The image.NRGBA tile to be drawn onto the canvas.
//
// Returns:
// - A pointer to the rotated image.NRGBA tile.
func (t *tileSet) setTileWithRotationLeft(x, y int, tile *image.NRGBA) *image.NRGBA {
	// Rotate the tile 90 degrees to the left.
	rotatedTile := rotateLeft90(tile)

	// Draw the rotated tile onto the canvas at the specified coordinates.
	t.setTile(x, y, rotatedTile)

	// Return the rotated tile.
	return rotatedTile
}

func (t *tileSet) getCanvas() *image.NRGBA {
	return t.canvas
}
