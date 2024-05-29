package tiled

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type Template struct {
	TiledVersion string
	TileWidth    int
	TileHeight   int
	TileCount    int
	Columns      int
	ImageSource  string // e.g. ../14x2_output.local.png
	ImageWidth   int
	ImageHeight  int
	Terrain1Name string
	Terrain2Name string
	TilesetName  string
}

type Exporter struct {
	tmpl *Template
}

var (
	errNilTemplate         = errors.New("template is nil")
	errInvalidTemplate     = errors.New("invalid template")
	errInvalidWidth        = errors.New("invalid width")
	errInvalidHeight       = errors.New("invalid height")
	errInvalidTileCount    = errors.New("invalid tile count")
	errInvalidColumns      = errors.New("invalid columns")
	errInvalidTerrain1Name = errors.New("invalid terrain1 name")
	errInvalidTerrain2Name = errors.New("invalid terrain2 name")
	errInvalidImageSource  = errors.New("invalid image source")
	errInvalidImageWidth   = errors.New("invalid image width")
	errInvalidImageHeight  = errors.New("invalid image height")
	errInvalidTileWidth    = errors.New("invalid tile width")
	errInvalidTileHeight   = errors.New("invalid tile height")
	errInvalidTiledVersion = errors.New("invalid tiled version. expected 1.10.x")
	errInvalidTilesetName  = errors.New("invalid tileset name")
)

func NewExporter(template *Template) *Exporter {
	return &Exporter{
		tmpl: template,
	}
}

func NewTemplate(
	tilesetName string,
	tileWidth int,
	tileHeight int,
	tileCount int,
	columns int,
	imageSource string,
	imageWidth int,
	imageHeight int,
	terrain1Name string,
	terrain2Name string,
) *Template {
	return &Template{
		TiledVersion: "1.10.2",
		TilesetName:  tilesetName,
		TileWidth:    tileWidth,
		TileHeight:   tileHeight,
		TileCount:    tileCount,
		Columns:      columns,
		ImageSource:  imageSource,
		ImageWidth:   imageWidth,
		ImageHeight:  imageHeight,
		Terrain1Name: terrain1Name,
		Terrain2Name: terrain2Name,
	}
}

func (e *Exporter) Export() error {
	err := e.tmpl.Validate()
	if err != nil {
		return err
	}
	var tmplFile = "templates/tiled/v1_10/14x2_tiled_1.10.x.tsx.tmpl"
	tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)
	if err != nil {
		return err
	}

	outFile, err := os.Create(fmt.Sprintf("out/%s.tsx", e.tmpl.TilesetName))
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	err = tmpl.Execute(outFile, e.tmpl)
	if err != nil {
		return err
	}
	return nil
}

func (t *Template) Validate() error {
	errs := make([]error, 0)

	if t == nil {
		return errNilTemplate
	}
	if t.TilesetName == "" {
		errs = append(errs, errInvalidTilesetName)
	}
	if t.TileWidth == 0 {
		errs = append(errs, errInvalidTileWidth)
	}
	if t.TileHeight == 0 {
		errs = append(errs, errInvalidTileHeight)
	}
	if t.TileCount == 0 {
		errs = append(errs, errInvalidTileCount)
	}
	if t.Columns == 0 {
		errs = append(errs, errInvalidColumns)
	}
	if t.Terrain1Name == "" {
		errs = append(errs, errInvalidTerrain1Name)
	}
	if t.Terrain2Name == "" {
		errs = append(errs, errInvalidTerrain2Name)
	}
	if t.ImageSource == "" {
		errs = append(errs, errInvalidImageSource)
	}
	if t.ImageWidth == 0 {
		errs = append(errs, errInvalidWidth)
	}
	if t.ImageHeight == 0 {
		errs = append(errs, errInvalidHeight)
	}
	if t.TiledVersion == "" || !strings.HasPrefix(t.TiledVersion, "1.10") {
		errs = append(errs, errInvalidTiledVersion)
	}
	if t.ImageHeight == 0 {
		errs = append(errs, errInvalidImageHeight)
	}
	if t.ImageWidth == 0 {
		errs = append(errs, errInvalidImageWidth)
	}

	if len(errs) > 0 {
		errs = append(errs, errInvalidTemplate)
		return errors.Join(errs...)
	}

	return nil
}
