package show

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/sfomuseum/go-www-show/v2"
)

// Put these in go-www-show
const leaflet_osm_tile_url = "https://tile.openstreetmap.org/{z}/{x}/{y}.png"
const protomaps_api_tile_url string = "https://api.protomaps.com/tiles/v3/{z}/{x}/{y}.mvt?key={key}"

var port int
var verbose bool

var browser_uri string

var map_provider string
var base_tile_uri string

var protomaps_theme string

var raster_tiles multi.KeyValueString
var vector_tiles multi.KeyValueString

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("show")

	browser_schemes := show.BrowserSchemes()
	str_schemes := strings.Join(browser_schemes, ",")

	browser_desc := fmt.Sprintf("A valid sfomuseum/go-www-show/v2.Browser URI. Valid options are: %s", str_schemes)

	fs.StringVar(&browser_uri, "browser-uri", "web://", browser_desc)

	fs.StringVar(&map_provider, "map-provider", "maplibre", "The map provider to use for a base layer. Valid options are: leaflet, maplibre, protomaps")
	fs.StringVar(&base_tile_uri, "base-tile-uri", leaflet_osm_tile_url, "A valid raster tile layer or pmtiles:// URI.")

	fs.StringVar(&protomaps_theme, "protomaps-theme", "white", "A valid Protomaps theme label.")

	fs.IntVar(&port, "port", 0, "The port number to listen for requests on (on localhost). If 0 then a random port number will be chosen.")

	fs.Var(&raster_tiles, "raster", "Zero or more {LAYER_NAME}={PATH} pairs referencing MBTiles databases containing raster data.")
	fs.Var(&vector_tiles, "vector", "Zero or more {LAYER_NAME}={PATH} pairs referencing MBTiles databases containing vector (MVT) data.")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Command-line tool for serving MBTiles tiles from an on-demand web server.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}

	return fs
}
