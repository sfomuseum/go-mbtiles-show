package http

import (
	"github.com/sfomuseum/go-mbtiles-server/tile"
	"github.com/tilezen/go-tilepacks/tilepack"
	"log"
	gohttp "net/http"
	"strconv"
)

// MBTilesHandler will return a http.HandlerFunc handler for serving tile requests from 'catalog' using the
// tile.NewSimpleTileParser parser.
func MBTilesHandler(catalog map[string]tilepack.MbtilesReader) (gohttp.HandlerFunc, error) {

	p, err := tile.NewSimpleTileParser()

	if err != nil {
		return nil, err
	}

	return MBTilesHandlerWithParser(catalog, p)
}

// MBTilesHandler will return a http.HandlerFunc handler for serving tile requests from 'catalog' using a
// custom tile.TileParser instance.
func MBTilesHandlerWithParser(catalog map[string]tilepack.MbtilesReader, p tile.TileParser) (gohttp.HandlerFunc, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		path := req.URL.Path
		tile_req, err := p.Parse(path)

		if err != nil {
			gohttp.NotFound(rsp, req)
			return
		}

		reader, ok := catalog[tile_req.Layer]

		if !ok {
			gohttp.NotFound(rsp, req)
			return
		}

		result, err := reader.GetTile(tile_req.Tile)

		if err != nil {
			log.Printf("Error getting tile: %+v", err)
			gohttp.NotFound(rsp, req)
			return
		}

		if result.Data == nil {
			gohttp.NotFound(rsp, req)
			return
		}

		l := len(*result.Data)
		str_l := strconv.Itoa(l)
		
		rsp.Header().Set("Content-Type", tile_req.ContentType)
		rsp.Header().Set("Content-Length", str_l)
		rsp.Write(*result.Data)
	}

	return fn, nil
}
