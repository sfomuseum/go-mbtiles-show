package show

// mapConfig defines common configuration details for maps.
type mapConfig struct {
	// A valid map provider label.
	Provider string `json:"provider"`
	// A valid Leaflet tile layer URI.
	TileURL string `json:"tile_url"`
	// Optional Protomaps configuration details
	Protomaps       *protomapsConfig `json:"protomaps,omitempty"`
	MBTilesLayers map[string]string `json:"mbtiles_layers"`
}

// protomapsConfig defines configuration details for maps using Protomaps.
type protomapsConfig struct {
	// A valid Protomaps theme label
	Theme string `json:"theme"`
}
