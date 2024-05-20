# Autotiler
Unpacks tilesets from a compressed one

Pure go. No external dependencies

Tileset-Generator convert tileset like this:

![packed](./examples/2x3_packed.png)

To tileset like this:

![complete](examples/output/tileset/output.15.png)

## How to use
* put simple tileset image (for example 2x3_packed.png) to source folder
* run
  ```go run . ./examples/2x3_packed.png output.local.png```
* grab complite tileset (for example output.local.png)
* enjoy

## Roadmap
- [x] Unpack from 6 tiles to 15 tiles
- [ ] Unpack from 6 or 15 tiles to 47 tiles
- [ ] Unpack from 6, 15 or 47 tiles to 256 tiles
- [ ] Pack 256 -> 47 -> 15 -> 6
- [ ] Export to Tiled
- [ ] Export to Godot