package show

import (
	"encoding/json"
	"net/http"
)

// LocalConfig defines application-specific configuration details.
type LocalConfig struct {
	RasterLayers map[string]string `json:"raster_layers,omitempty"`
	VectorLayers map[string]string `json:"vector_layers,omitempty"`
}

// LocalConfigHandler returns an `http.Handler` instance that when called will return 'cfg' as a JSON-encoded string.
func LocalConfigHandler(cfg *LocalConfig) http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		rsp.Header().Set("Content-type", "application/json")

		enc := json.NewEncoder(rsp)
		err := enc.Encode(cfg)

		if err != nil {
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	return http.HandlerFunc(fn)
}
