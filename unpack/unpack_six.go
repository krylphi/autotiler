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
	quadMap := append(exp[:14], exp[15], exp[14]) // 15 - filled terrain 2, 14 - small terrain 1 patch on top

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
	quadMap := append(exp[14:], exp[1], exp[0]) // 0 - filled terrain 1, 1 - small terrain 2 patch on top

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

func (u *Unpacker) From6to48Terrain1() (*image.NRGBA, error) {
	quadMap := export6to28TileSet()
	quads := [16]quadTileData{
		quadMap[1],
		quadMap[4],
		&[4][2]int{
			{3, 1}, {2, 1},
			{3, 0}, {2, 0},
		},
		quadMap[2],
		quadMap[5],
		quadMap[9],
		&[4][2]int{
			{0, 0}, {0, 5},
			{3, 2}, {0, 2},
		},
		quadMap[13],
		&[4][2]int{
			{1, 5}, {2, 5},
			{3, 2}, {1, 1},
		},
		&[4][2]int{
			{1, 5}, {2, 5},
			{0, 1}, {0, 2},
		},
		quadMap[3],
		&[4][2]int{
			{3, 5}, {0, 5},
			{0, 1}, {1, 1},
		},
		quadMap[8],
		quadMap[0],
		quadMap[14],
		quadMap[12],
	}

	return u.from6to48Terrain(quads)
}

func (u *Unpacker) From6to48Terrain2() (*image.NRGBA, error) {
	exp := export6to28TileSet()
	quadMap := append(exp[14:], exp[1], exp[0]) // 0 - filled terrain 1, 1 - small terrain 2 patch on top
	quads := [16]quadTileData{
		quadMap[1],
		quadMap[4],
		&[4][2]int{
			{0, 2}, {3, 2},
			{0, 5}, {3, 5},
		},
		quadMap[2],
		quadMap[5],
		quadMap[9],
		// 2nd 16
		&[4][2]int{
			{1, 3}, {3, 0},
			{2, 1}, {3, 1},
		},
		quadMap[13],
		&[4][2]int{
			{1, 2}, {2, 2},
			{2, 1}, {2, 4},
		},
		&[4][2]int{
			{1, 2}, {2, 2},
			{1, 4}, {3, 1},
		},
		// 3rd 16
		quadMap[3],
		&[4][2]int{
			{2, 0}, {3, 0},
			{1, 4}, {2, 4},
		},
		quadMap[8],
		quadMap[0],
		quadMap[15],
		quadMap[12],
	}

	return u.from6to48Terrain(quads)
}

func (u *Unpacker) from6to48Terrain(quadMap [16]quadTileData) (*image.NRGBA, error) {
	if u.xTiles*u.yTiles != sixPackType {
		return nil, errInvalidPackType
	}

	// todo pass list of tile patterns so it would be reusing code for terrain 2
	canvas := image.NewNRGBA(image.Rect(0, 0, u.tileWidth*12, u.tileHeight*4))
	tileset := NewTileSet(canvas, u.tileWidth, u.tileHeight)

	tilePattern := quadMap[0] // 1 - small terrain 2 patch on top
	tile := image.NewNRGBA(image.Rect(0, 0, u.tileWidth, u.tileHeight))
	// fully drawn
	u.drawFullSingleTile(tile, tilePattern)
	tileset.SetTile(0, 2, tile)
	tile = tileset.SetTileWithRotationLeft(3, 3, tile)
	tile = tileset.SetTileWithRotationLeft(0, 0, tile)
	tile = tileset.SetTileWithRotationLeft(1, 3, tile)

	tilePattern = quadMap[1]
	u.drawFullSingleTile(tile, tilePattern)
	tileset.SetTile(0, 1, tile)
	tileset.SetTileWithRotationLeft(2, 3, tile)

	tilePattern = quadMap[2]
	u.drawFullSingleTile(tile, tilePattern)
	tileset.SetTile(0, 3, tile)

	tilePattern = quadMap[3]
	u.drawFullSingleTile(tile, tilePattern)
	tileset.SetTile(1, 2, tile)
	tile = tileset.SetTileWithRotationLeft(3, 2, tile)
	tile = tileset.SetTileWithRotationLeft(3, 0, tile)
	tile = tileset.SetTileWithRotationLeft(1, 0, tile)

	tilePattern = quadMap[4]
	u.drawFullSingleTile(tile, tilePattern)
	tileset.SetTile(1, 1, tile)
	tile = tileset.SetTileWithRotationLeft(2, 2, tile)
	tile = tileset.SetTileWithRotationLeft(3, 1, tile)
	tile = tileset.SetTileWithRotationLeft(2, 0, tile)

	tilePattern = quadMap[5]
	u.drawFullSingleTile(tile, tilePattern)
	tileset.SetTile(2, 1, tile)

	// 2nd 16
	tilePattern = quadMap[6]
	u.drawFullSingleTile(tile, tilePattern)
	tileset.SetTile(4, 0, tile)
	tile = tileset.SetTileWithRotationLeft(4, 3, tile)
	tile = tileset.SetTileWithRotationLeft(7, 3, tile)
	tile = tileset.SetTileWithRotationLeft(7, 0, tile)

	tilePattern = quadMap[7]
	u.drawFullSingleTile(tile, tilePattern)
	tileset.SetTile(5, 1, tile)
	tile = tileset.SetTileWithRotationLeft(5, 2, tile)
	tile = tileset.SetTileWithRotationLeft(6, 2, tile)
	tile = tileset.SetTileWithRotationLeft(6, 1, tile)

	tilePattern = quadMap[8]
	u.drawFullSingleTile(tile, tilePattern)
	tileset.SetTile(5, 0, tile)
	tile = tileset.SetTileWithRotationLeft(4, 2, tile)
	tile = tileset.SetTileWithRotationLeft(6, 3, tile)
	tile = tileset.SetTileWithRotationLeft(7, 1, tile)

	tilePattern = quadMap[9]
	u.drawFullSingleTile(tile, tilePattern)
	tileset.SetTile(6, 0, tile)
	tile = tileset.SetTileWithRotationLeft(4, 1, tile)
	tile = tileset.SetTileWithRotationLeft(5, 3, tile)
	tile = tileset.SetTileWithRotationLeft(7, 2, tile)

	// 3rd 16
	tilePattern = quadMap[10]
	u.drawFullSingleTile(tile, tilePattern)
	tileset.SetTile(8, 3, tile)
	tile = tileset.SetTileWithRotationLeft(11, 3, tile)
	tile = tileset.SetTileWithRotationLeft(11, 0, tile)
	tile = tileset.SetTileWithRotationLeft(8, 0, tile)

	tilePattern = quadMap[11]
	u.drawFullSingleTile(tile, tilePattern)
	tileset.SetTile(9, 0, tile)
	tile = tileset.SetTileWithRotationLeft(8, 2, tile)
	tile = tileset.SetTileWithRotationLeft(10, 3, tile)
	tile = tileset.SetTileWithRotationLeft(11, 2, tile)

	tilePattern = quadMap[12]
	u.drawFullSingleTile(tile, tilePattern)
	tileset.SetTile(8, 1, tile)
	tile = tileset.SetTileWithRotationLeft(9, 3, tile)
	tile = tileset.SetTileWithRotationLeft(11, 1, tile)
	tile = tileset.SetTileWithRotationLeft(10, 0, tile)

	tilePattern = quadMap[13]
	u.drawFullSingleTile(tile, tilePattern)
	tileset.SetTile(10, 1, tile)

	tilePattern = quadMap[14]
	u.drawFullSingleTile(tile, tilePattern)
	tileset.SetTile(9, 2, tile)

	tilePattern = quadMap[15]
	u.drawFullSingleTile(tile, tilePattern)
	tileset.SetTile(9, 1, tile)
	tile = tileset.SetTileWithRotationLeft(10, 2, tile)

	return tileset.GetCanvas(), nil
}

func export6to28TileSet() [28]quadTileData {
	return [28]quadTileData{
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
