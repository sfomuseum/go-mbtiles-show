package show

import (
	"context"
	"flag"
	"fmt"

	"github.com/sfomuseum/go-flags/flagset"
	www_show "github.com/sfomuseum/go-www-show"
	"github.com/tilezen/go-tilepacks/tilepack"
)

type RunOptions struct {
	MapProvider     string
	MapTileURI      string
	ProtomapsTheme  string
	Port            int
	MBTilesCatalog map[string]tilepack.MbtilesReader
	Browser         www_show.Browser
}

func RunOptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	mbtiles_catalog := make(map[string]tilepack.MbtilesReader)

	for _, kv := range tiles_config {

		k := kv.Key()
		path := kv.Value().(string)

		r, err := tilepack.NewMbtilesReader(path)

		if err != nil {
			return nil, fmt.Errorf("Failed to create MBTiles reader for %s, %w", path, err)
		}

		mbtiles_catalog[k] = r
	}
	
	opts := &RunOptions{
		MapProvider:     map_provider,
		MapTileURI:      map_tile_uri,
		ProtomapsTheme:  protomaps_theme,
		MBTilesCatalog: mbtiles_catalog,
		Port:            port,
	}
	
	br, err := www_show.NewBrowser(ctx, "web://")

	if err != nil {
		return nil, fmt.Errorf("Failed to create new browser, %w", err)
	}

	opts.Browser = br

	return opts, nil
}
