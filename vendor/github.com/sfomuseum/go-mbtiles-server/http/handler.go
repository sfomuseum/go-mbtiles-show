package http

import (
	"log/slog"
	gohttp "net/http"
	"strconv"

	"github.com/sfomuseum/go-mbtiles-server/tile"
	"github.com/tilezen/go-tilepacks/tilepack"
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

		logger := slog.Default()
		logger = logger.With("path", req.URL.Path)

		path := req.URL.Path
		tile_req, err := p.Parse(path)

		if err != nil {
			logger.Error("Failed to parse tile path", "error", err)
			gohttp.NotFound(rsp, req)
			return
		}

		logger = logger.With("layer", tile_req.Layer)

		reader, ok := catalog[tile_req.Layer]

		if !ok {
			logger.Error("Tile layer not found")
			gohttp.NotFound(rsp, req)
			return
		}

		result, err := reader.GetTile(tile_req.Tile)

		if err != nil {
			logger.Error("Failed to retrieve tile", "error", err)
			gohttp.NotFound(rsp, req)
			return
		}

		if result.Data == nil {
			logger.Debug("Tile data is nil")
			gohttp.NotFound(rsp, req)
			return
		}

		l := len(*result.Data)
		str_l := strconv.Itoa(l)

		logger.Debug("Serve tile", "content type", tile_req.ContentType, "length", l)

		rsp.Header().Set("Content-Type", tile_req.ContentType)
		rsp.Header().Set("Content-Length", str_l)
		rsp.Write(*result.Data)
	}

	return fn, nil
}
