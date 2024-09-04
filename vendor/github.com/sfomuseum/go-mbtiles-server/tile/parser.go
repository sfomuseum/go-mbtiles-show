package tile

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/aaronland/go-mimetypes"
	"github.com/paulmach/orb/maptile"
)

// TileParser is an interface for parsing URI paths in to *TileRequest instances.
type TileParser interface {
	// Parse takes a URI path as its input and returns a *TileRequest instance.
	Parse(string) (*TileRequest, error)
}

// SimpleTileParse is an implementation of the TileParser interface. It uses a simple
// positional regular expression to match tile requests and determine which to return
// from a given MBTiles database:
//
//	`\/([^\/]+)\/(\d+)\/(\d+)\/(\d+)\.(\w+)$`
//
// Where:
// $1 is the name of the (MBTiles database) layer to read data from
// $2 is the Z value of the tile request
// $3 is the X value of the tile request
// $4 is the Y value of the tile request
// $5 is the mime-type extension for the tile request. Fully-qualified content-type values
// for file extensions are mapped using the aaronland/go-mimetypes package.
type SimpleTileParser struct {
	TileParser
	re *regexp.Regexp
}

// Return a new SimpleTileParser instance.
func NewSimpleTileParser() (TileParser, error) {

	re, err := regexp.Compile(`\/([^\/]+)\/(\d+)\/(\d+)\/(\d+)\.(\w+)$`)

	if err != nil {
		return nil, err
	}

	p := &SimpleTileParser{
		re: re,
	}

	return p, nil
}

// Parse a URI path in to a *TileRequest.
func (p *SimpleTileParser) Parse(path string) (*TileRequest, error) {

	match := p.re.FindStringSubmatch(path)

	if match == nil {
		return nil, fmt.Errorf("invalid tile path")
	}

	layer := match[1]
	ext := match[5]

	var content_type string

	switch ext {
	case "mvt":
		content_type = "application/vnd.mapbox-vector-tile"
	case "jpg", "jpeg":
		content_type = "image/jpeg"
	case "png":
		content_type = "image/png"
	default:

		t := mimetypes.TypesByExtension(ext)

		if len(t) == 0 {
			return nil, fmt.Errorf("Unsupported extension '%s'", ext)
		}

		content_type = t[0]
	}

	z, _ := strconv.ParseUint(match[2], 10, 32)
	x, _ := strconv.ParseUint(match[3], 10, 32)
	y, _ := strconv.ParseUint(match[4], 10, 32)

	// Just always invert the y coordinate because MBTiles
	// https://gist.github.com/tmcw/4954720
	// https://stackoverflow.com/questions/46822094/incorrect-coordinates-in-mbtiles-generated-with-tippecanoe

	inverted_y := (1 << z) - 1 - y

	tile := maptile.New(uint32(x), uint32(inverted_y), maptile.Zoom(z))

	tile_req := &TileRequest{
		Tile:        tile,
		Layer:       layer,
		ContentType: content_type,
	}

	return tile_req, nil
}
