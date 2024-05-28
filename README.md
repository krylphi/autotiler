# Autotiler
Unpacks tilesets from a compressed one

Pure go. No external dependencies

Tileset-Generator unpacks tileset like this:

![packed](./examples/2x3_packed.png)

To tilesets like this:

15x1 Terrain 1:

![15x1_T1](examples/output/tileset/15x1_terrain1_output.png)

15x1 Terrain 2:

![15x1_T1](examples/output/tileset/15x1_terrain2_output.png)

14x2:

![14x2](examples/output/tileset/14x2_output.png)

8x12:

![8x12](examples/output/tileset/8x12_output.png)



## How to use
* put simple tileset image (for example 2x3_packed.png) to source folder
* run
  ```go run . ./examples/2x3_packed.png output.local.png```
* grab complite tileset (for example output.local.png)
* enjoy

## Features and plans
- [x] Unpack from 6 tiles to 15 tiles
- [x] Unpack from 6 tiles to 28 tiles
- [x] Unpack from 6 or 28 tiles to 92 tiles
- [ ] Unpack from 6 or 15 tiles to 47 tiles
- [ ] Export to Tiled
- [ ] Export to Godot