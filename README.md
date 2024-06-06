# Autotiler
Unpacks tilesets from a compressed one

Pure go. No external dependencies

Tileset-Generator unpacks tileset like this:

![packed](./examples/2x3_packed.png)

To tilesets like this:

16x1 Terrain 1 to 2:

![16x1_T1](examples/output/tileset/16x1_terrain1_output.png)

16x1 Terrain 2 to 1:

![16x1_T1](examples/output/tileset/16x1_terrain2_output.png)

14x2:

![14x2](examples/output/tileset/14x2_output.png)


## How to use
* [get yourself Go](https://go.dev/doc/install) 
* clone this repository or download sources.
* put simple tileset image (for example 2x3_packed.png) to source folder
* run ```go run . [src image] [dst image]```
  e.g. ```go run . ./examples/2x3_packed.png output.local.png```
* grab complete tilesets from `out` directory
* enjoy
* alternatively you can build an application using `make build` command to use it as a standalone application without Go

## Features and plans
- [x] Unpack from 6 tiles to 16 tiles
- [x] Unpack from 6 tiles to 28 tiles
- [x] Unpack from 6 to 47 tiles
- [ ] Unpack from 6 or 16 tiles to 256 tiles
- [ ] Export to Tiled
- [ ] Export to Godot
- [ ] More build options (Win, Mac)
- [ ] Document and prettify code