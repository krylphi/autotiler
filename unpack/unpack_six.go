package unpack

import (
	"image"
)

const (
	sixPackType = 6
)

func (u *Unpacker) From6to16Terrain1() (*image.NRGBA, error) {
	if u.xTiles*u.yTiles != sixPackType {
		return nil, errInvalidPackType
	}
	canvas := image.NewNRGBA(image.Rect(0, 0, u.tileWidth*16, u.tileHeight*1))
	// todo optimize to generate automatically and consider scaling for 47 and 255 tilesets
	exp := export6to28TileSet()
	quadMap := append(exp[:14], exp[15], exp[14]) // 15 - filled terrain 2, 14 - small terrain 2 patch on top

	for idx := 0; idx < 16; idx++ {
		u.drawFullTile(canvas, quadMap[idx], idx, 16, 2)
	}
	return canvas, nil
}

func (u *Unpacker) From6to16Terrain2() (*image.NRGBA, error) {
	if u.xTiles*u.yTiles != sixPackType {
		return nil, errInvalidPackType
	}
	canvas := image.NewNRGBA(image.Rect(0, 0, u.tileWidth*16, u.tileHeight*1))
	// todo optimize to generate automatically and consider scaling for 47 and 255 tilesets
	exp := export6to28TileSet()
	quadMap := append(exp[14:], exp[1], exp[0]) // 0 - filled terrain 1, 1 - small terrain 1 patch on top

	for idx := 0; idx < 16; idx++ {
		u.drawFullTile(canvas, quadMap[idx], idx, 16, 2)
	}
	return canvas, nil
}

func (u *Unpacker) From6to28() (*image.NRGBA, error) {
	if u.xTiles*u.yTiles != sixPackType {
		return nil, errInvalidPackType
	}
	canvas := image.NewNRGBA(image.Rect(0, 0, u.tileWidth*14, u.tileHeight*2))
	// todo optimize to generate automatically and consider scaling for 47 and 255 tilesets

	quadMap := export6to28TileSet()

	for idx := 0; idx < 28; idx++ {
		u.drawFullTile(canvas, quadMap[idx], idx, 14, 2)
	}
	return canvas, nil
}

// Deprecated: inconsistent behaviour in Tiled
func (u *Unpacker) from28To92(img *image.NRGBA) *image.NRGBA {
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

func export6to28TileSet() [28]tileData {
	return [28]tileData{
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
			{1, 0}, {1, 1},
		},
		// end of terrain 1
		// terrain 2
		{
			{0, 0}, {1, 0}, // tile 15
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
