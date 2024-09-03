window.addEventListener("load", function load(event){

    // Null Island
    // var lat = 0.0;
    // var lon = 0.0;
    // var zoom = 3;
    
    // SFO 
    var lat = 37.621131;
    var lon = -122.384292;
    var zoom = 14;
    
    var init = function(cfg){
	init_maplibre(cfg);
    }

    var init_maplibre = function(cfg){

	var base_souce = {};
	var base_layer = {};

	switch (cfg.provider) {
	    case "leaflet":
		
		base_source = {
		    type: 'raster',
		    tiles: [
			cfg.tile_url,
		    ],
		    'tileSize': 256,
		};

		base_layer = {
		    'id': 'base',
		    'type': 'raster',
		    'source': 'base',
		};
		
		break;
		
	    case "protomaps":

		// add the PMTiles plugin to the maplibregl global.
		// https://maplibre.org/maplibre-gl-js/docs/examples/pmtiles/
		// https://github.com/protomaps/PMTiles/blob/main/js/examples/maplibre.html
		// https://unpkg.com/pmtiles@3.0.7/dist/pmtiles.js
		
		const protocol = new pmtiles.Protocol();
		maplibregl.addProtocol('pmtiles', protocol.tile);

		const p = new pmtiles.PMTiles(cfg.tile_url);
		protocol.add(p);
		
		base_source = {
		    type: "vector",
		    url: "pmtiles://" + cfg.tile_url,
		};

		base_layer = {
		    id: 'base',
		    source: 'base',
		    'source-layer': 'base',
		    type: "line",
                    paint: {
                        "line-color": "#fc8d62",
                    }
		};
		
		break;
		
	    default:
		console.error("Unsupported provider", cfg.provider)
		return;
	}

	// console.log("BASE", base_source, base_layer);
	
	var map_args = {
            container: 'map',
	    center: [ lon, lat ],
	    zoom: zoom,
	    style: {
		version: 8,
		sources: {
		    'base': base_source,
		},
		layers: [
		    base_layer,
		]
	    }
	};

	var legend = {
	    'base': [ 'base' ],
	};
	
	var map = new maplibregl.Map(map_args);
	
	map.on('load', () => {
	    console.log("Map done loading");
	    
	    if (cfg.raster_layers){

		// Basically inverted-y coordinates ({-y}) are not supported in maplibre-gl.js
		// https://maplibre.org/maplibre-style-spec/sources/#raster
		// https://maplibre.org/maplibre-gl-js/docs/API/type-aliases/CanvasSourceSpecification/
		// Despite seemingly being supported in the "native" builds...
		// https://docs.mapbox.com/ios/maps/api/6.4.1/tile-url-templates.html
		// https://maplibre.org/maplibre-native/docs/book/design/coordinate-system.html
		    
		for (k in cfg.raster_layers){

		    var tile_url = "http://" + location.host + cfg.raster_layers[k];
		    console.log("Add raster layer", k, tile_url);
		    
		    map.addSource(k, {
			type: 'raster',
			tiles: [
			   tile_url,
			],
			tileSize: 256,
		    });
		    
		    map.addLayer({
			'id': k,
			'type': 'raster',
			'source': k,
			'source-layer': k,
			'layout': {
			    'visibility': 'none'
			},
		    });

		    legend[k] = [ k ];		    
		}
		
	    }
	    
	    
	    if (cfg.vector_layers){
		
		for (k in cfg.vector_layers){

		    var tile_url = "http://" + location.host + cfg.vector_layers[k];
		    console.log("ADD", k, tile_url);
		    
		    map.addSource(k, {
			type: 'vector',
			tiles: [
			   tile_url,
			],
		    });
		    
		    map.addLayer({
			'id': k,
			'type': 'line',
			'source': k,
			'source-layer': k,
			'layout': {
			    'visibility': 'none'
			},			
			paint: {
			    'line-color':'#000000',
			    'line-width': 1,
			    'line-opacity': 1
			}
		    });

		    legend[k] = [ k ];
		}
	    }

	    // Create control
	    let lc = new LayersControl(legend);
	    map.addControl(lc);
	});
	
    };
        
    fetch("/map.json")
	.then((rsp) => rsp.json())
	.then((cfg) => {	    
	    init(cfg);
	}).catch((err) => {
	    console.error("Failed to retrieve map config", err);
	});
        
    
});
