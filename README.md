# go-mbtiles-show

Command-line tool for serving MBTiles databases from an on-demand web server.

## Motivation

It's basically a simpler and dumber version of [geojson.io](https://geojson.io/) but for [MBTiles databases](https://wiki.openstreetmap.org/wiki/MBTiles) that you can run locally from a single binary application.

## Documentation

Documentation (`godoc`) is incomplete at this time.

## Important

This is an early-stage project. It is incomplete and lacking features. Notably:

* It only supports MBTiles databases with raster data. It does not support vector tiles yet. That's not a feature but a reflection of the need to write this tool in a hurry to debug a problem with a raster tiles database.

* It assumes that raster tiles are encoded as PNG data when constructing URLs and when constructing `Content-Type` headers when returning tile data.

* The initial map zoom is the [San Francisco International Airport](https://spelunker.whosonfirst.org/id/102527513/). Code to query the extent of each MBTiles database being served needs to be written and then logic about which extent(s) should be the default view needs to be decided on.

* Minimum and maximum zoom levels for MBTiles database layers need to be derived in the Go code and passed down to the JavaScript/Leaflet code.

_I will get to these things eventually but patches are welcome._

## Tools

```
$> make cli
go build -mod vendor -ldflags="-s -w" -o bin/show cmd/show/main.go
```

### show

```
$> ./bin/show -h
Command-line tool for serving MBTiles tiles from an on-demand web server.
Usage:
	 ./bin/show [options]
Valid options are:
  -map-provider string
    	The map provider to use for a base layer. Valid options are: leaflet, protomaps (default "leaflet")
  -map-tile-uri string
    	A valid Leaflet tile layer URI. See documentation for special-case (interpolated tile) URIs. (default "https://tile.openstreetmap.org/{z}/{x}/{y}.png")
  -port int
    	The port number to listen for requests on (on localhost). If 0 then a random port number will be chosen.
  -protomaps-theme string
    	A valid Protomaps theme label. (default "white")
  -raster value
    	Zero or more {LAYER_NAME}={PATH} pairs referencing MBTiles databases containing raster data.
  -verbose
    	Enable verbose (debug) logging.
```

#### Example

![](docs/images/go-mbtiles-show-2023-24.png)

```
$> ./bin/show \
	-raster 2023=/usr/local/sfomuseum/tiles/sqlite/2023.db \
	-raster 2024=/usr/local/sfomuseum/tiles/sqlite/2024.db \
	-verbose
	
2024/08/30 18:16:14 DEBUG Verbose logging enabled
2024/08/30 18:16:14 DEBUG Start server
2024/08/30 18:16:14 DEBUG HEAD request succeeded url=http://localhost:64211
2024/08/30 18:16:14 INFO Server is ready and features are viewable url=http://localhost:64115
2024/08/30 18:16:22 DEBUG Tile data is nil path=/tiles/2023/12/654/2511.png layer=2023
2024/08/30 18:16:22 DEBUG Tile data is nil path=/tiles/2023/12/655/2511.png layer=2023
2024/08/30 18:16:22 DEBUG Tile data is nil path=/tiles/2023/12/656/2511.png layer=2023
2024/08/30 18:16:22 DEBUG Serve tile path=/tiles/2023/12/655/2510.png layer=2023 "content type"=image/png length=49883
2024/08/30 18:16:22 DEBUG Tile data is nil path=/tiles/2023/12/654/2510.png layer=2023
2024/08/30 18:16:22 DEBUG Tile data is nil path=/tiles/2023/12/656/2510.png layer=2023
2024/08/30 18:16:22 DEBUG Tile data is nil path=/tiles/2023/12/655/2512.png layer=2023
... and so on
```

## See also

* https://github.com/sfomuseum/go-mbtiles-server
* https://github.com/sfomuseum/go-geojson-show
* https://github.com/tilezen/go-tilepacks