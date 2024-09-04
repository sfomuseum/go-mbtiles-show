package show

// mapConfig defines common configuration details for maps.
type mapConfig struct {
	// A valid map provider label.
	Provider string `json:"provider"`
	// A valid Leaflet tile layer URI.
	BaseTileURL string `json:"base_tile_url"`
	// Optional Protomaps configuration details
	Protomaps    *protomapsConfig  `json:"protomaps,omitempty"`
	RasterLayers map[string]string `json:"raster_layers,omitempty"`
	VectorLayers map[string]string `json:"vector_layers,omitempty"`
}

// protomapsConfig defines configuration details for maps using Protomaps.
type protomapsConfig struct {
	UsePMTiles bool `json:"use_pmtiles"`
	// A valid Protomaps theme label
	Theme string `json:"theme"`
}
