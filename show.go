package show

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/sfomuseum/go-http-protomaps"
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

	map_cfg := &mapConfig{
		Provider:     opts.MapProvider,
		BaseTileURL:  opts.BaseTileURI,
		RasterLayers: raster_layers,
		VectorLayers: vector_layers,
	}

	u, err := url.Parse(opts.BaseTileURI)

	if err != nil {
		return fmt.Errorf("Failed to parse Protomaps tile URL, %w", err)
	}

	switch u.Scheme {
	case "pmtiles":

		switch u.Host {
		case "api":
			if opts.MapProvider == "maplibre" {
				slog.Warn("Remote PMTiles endpoints don't seem to work yet.")
			}

			q := u.Query()
			key := q.Get("key")
			map_cfg.BaseTileURL = strings.Replace(protomaps_api_tile_url, "{key}", key, 1)

		case "":
			mux_url, mux_handler, err := protomaps.FileHandlerFromPath(u.Path, "")

			if err != nil {
				return fmt.Errorf("Failed to determine absolute path for '%s', %v", opts.BaseTileURI, err)
			}

			mux.Handle(mux_url, mux_handler)
			map_cfg.BaseTileURL = mux_url

		default:
			return fmt.Errorf("Unsupported host (%s) for pmtiles base URI", u.Host)
		}

		map_cfg.Protomaps = &protomapsConfig{
			UsePMTiles: true,
			Theme:      opts.ProtomapsTheme,
		}

	default:
		map_cfg.BaseTileURL = opts.BaseTileURI
	}

	map_cfg_handler := mapConfigHandler(map_cfg)

	mux.Handle("/map.json", map_cfg_handler)

	www_show_opts := &www_show.RunOptions{
		Port:    opts.Port,
		Browser: opts.Browser,
		Mux:     mux,
	}

	return www_show.RunWithOptions(ctx, www_show_opts)
}

func mapConfigHandler(cfg *mapConfig) http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		rsp.Header().Set("Content-type", "application/json")

		enc := json.NewEncoder(rsp)
		err := enc.Encode(cfg)

		if err != nil {
			slog.Error("Failed to encode map config", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
		}

		return
	}

	return http.HandlerFunc(fn)
}
