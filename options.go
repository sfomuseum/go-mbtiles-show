package show

import (
	"context"
	"flag"
	"fmt"

	"github.com/sfomuseum/go-flags/flagset"
	www_show "github.com/sfomuseum/go-www-show/v2"
	"github.com/tilezen/go-tilepacks/tilepack"
)

type RunOptions struct {
	MapProvider    string
	BaseTileURI    string
	ProtomapsTheme string
	Port           int
	RasterCatalog  map[string]tilepack.MbtilesReader
	VectorCatalog  map[string]tilepack.MbtilesReader
	Browser        www_show.Browser
	Verbose        bool
}

func RunOptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	raster_catalog := make(map[string]tilepack.MbtilesReader)
	vector_catalog := make(map[string]tilepack.MbtilesReader)

	for _, kv := range raster_tiles {

		k := kv.Key()
		path := kv.Value().(string)

		r, err := tilepack.NewMbtilesReader(path)

		if err != nil {
			return nil, fmt.Errorf("Failed to create MBTiles reader for %s, %w", path, err)
		}

		raster_catalog[k] = r
	}

	for _, kv := range vector_tiles {

		k := kv.Key()
		path := kv.Value().(string)

		r, err := tilepack.NewMbtilesReader(path)

		if err != nil {
			return nil, fmt.Errorf("Failed to create MBTiles reader for %s, %w", path, err)
		}

		vector_catalog[k] = r
	}

	opts := &RunOptions{
		MapProvider:    map_provider,
		BaseTileURI:    base_tile_uri,
		ProtomapsTheme: protomaps_theme,
		RasterCatalog:  raster_catalog,
		VectorCatalog:  vector_catalog,
		Port:           port,
		Verbose:        verbose,
	}

	br, err := www_show.NewBrowser(ctx, browser_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new browser, %w", err)
	}

	opts.Browser = br

	return opts, nil
}
