package tile

import (
	"fmt"

	"github.com/aaronland/go-mimetypes"
	"github.com/paulmach/orb/maptile"
)

// TileRequest provides a simple struct encapsulating the specific details of a tile request.
type TileRequest struct {
	// Tile is the underlying *tilepack.Tile instance where tile data is stored
	Tile maptile.Tile
	// Layer is the name of the MBTiles to read from
	Layer string
	// ContentType is the mime-type of the data to return
	ContentType string
}

// String returns a relative URI path for the TileRequest instance.
func (t *TileRequest) String() string {

	ext := mimetypes.ExtensionsByType(t.ContentType)

	if len(ext) == 0 {
		ext = []string{".???"}
	}

	return fmt.Sprintf("%s/%d/%d/%d.%s", t.Layer, t.Tile.Z, t.Tile.X, t.Tile.Y, ext[0])
}
