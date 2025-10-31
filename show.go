package show

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aaronland/go-http-maps/v2"
	mbtiles_http "github.com/sfomuseum/go-mbtiles-server/http"
	"github.com/sfomuseum/go-mbtiles-show/static/www"
	www_show "github.com/sfomuseum/go-www-show/v2"
	"github.com/tilezen/go-tilepacks/tilepack"
)

func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := RunOptionsFromFlagSet(ctx, fs)

	if err != nil {
		return err
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	if opts.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	mux := http.NewServeMux()

	www_fs := http.FS(www.FS)
	mux.Handle("/", http.FileServer(www_fs))

	//

	tiles_catalog := make(map[string]tilepack.MbtilesReader)

	for k, v := range opts.RasterCatalog {
		tiles_catalog[k] = v
	}

	for k, v := range opts.VectorCatalog {
		tiles_catalog[k] = v
	}

	mbtiles_handler, err := mbtiles_http.MBTilesHandler(tiles_catalog)

	if err != nil {
		return err
	}

	mux.Handle("/tiles/", mbtiles_handler)

	map_opts := &maps.AssignMapConfigHandlerOptions{
		MapProvider:          opts.MapProvider,
		MapTileURI:           opts.BaseTileURI,
		ProtomapsTheme:       opts.ProtomapsTheme,
		ProtomapsMaxDataZoom: opts.ProtomapsMaxDataZoom,
	}

	err = maps.AssignMapConfigHandler(map_opts, mux, "/map.json")

	if err != nil {
		return fmt.Errorf("Failed to assign map config handler, %w", err)
	}

	raster_layers := make(map[string]string, 0)

	for k, _ := range opts.RasterCatalog {
		layer_url := fmt.Sprintf("/tiles/%s/{z}/{x}/{y}.png", k)
		raster_layers[k] = layer_url
	}

	vector_layers := make(map[string]string, 0)

	for k, _ := range opts.VectorCatalog {
		layer_url := fmt.Sprintf("/tiles/%s/{z}/{x}/{y}.mvt", k)
		vector_layers[k] = layer_url
	}

	if opts.MapProvider == "leaflet" && len(vector_layers) > 0 {
		slog.Warn("Leaflet map provider does not support rendering vector layers yet.")
	}

	local_cfg := &LocalConfig{
		RasterLayers: raster_layers,
		VectorLayers: vector_layers,
	}

	local_cfg_handler := LocalConfigHandler(local_cfg)

	mux.Handle("/config.json", local_cfg_handler)

	www_show_opts := &www_show.RunOptions{
		Port:    opts.Port,
		Browser: opts.Browser,
		Mux:     mux,
	}

	return www_show.RunWithOptions(ctx, www_show_opts)
}
