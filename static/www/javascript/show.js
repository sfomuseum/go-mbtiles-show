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
	// console.log(cfg);
	// init_leaflet(cfg);
	init_maplibre(cfg);
    }

    var init_maplibre = function(cfg){

	// If protomaps
	// https://maplibre.org/maplibre-gl-js/docs/examples/pmtiles/
	
	var map_args = {
            container: 'map',
	    center: [ lon, lat ],
	    zoom: zoom,
	    style: {
		version: 8,
		sources: {
		    'base': {
			type: 'raster',
			tiles: [
			    cfg.tile_url,
			],
			'tileSize': 256,
		    },
		},
		layers: [
		    {
			'id': 'base',
			'type': 'raster',
			'source': 'base',
		    }
		]
	    }
	};

	var map = new maplibregl.Map(map_args);
	
	map.on('load', () => {
	    console.log("Map done loading");

	    var legend = {};
	    
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
    
    var init_leaflet = function(cfg){

	var map = L.map('map');
	map.setView([lat, lon], zoom);
	
	var base_maps = {};
	var overlays = {};

	if (cfg.raster_layers){
	    
	    for (k in cfg.raster_layers){
		var layer_url = cfg.raster_layers[k];
		var l = L.tileLayer(layer_url)

		console.log("Add raster overlay", k, layer_url);		
		overlays[k] = l;
	    }
	}

	if (cfg.vector_layers){
	    
	    for (k in cfg.vector_layers){

		var layer_url = cfg.vector_layers[k]

		var layer_styles = {
		
		    k: function(properties, zoom) {
			return {
			    weight: 2,
			    color: 'red',
			    opacity: .5,
			    fillColor: 'yellow',
			    fill: true,
			    radius: 6,
			    fillOpacity: 0.1
			}
		    }
		};
		
		var layer_opts = {
		    rendererFactory: L.canvas.tile,
		    vectorTileLayerStyles: layer_styles,
		    interactive: true,
		};
	
		var layer = L.vectorGrid.protobuf(layer_url, layer_opts);
		console.log("Add vector overlay", k, layer_url);
		overlays[k] = layer;		
	    }
	    
	}

	switch (cfg.provider) {
	    case "leaflet":
		
		var tile_url = cfg.tile_url;
		
		var tile_layer = L.tileLayer(tile_url);
		tile_layer.addTo(map);

		base_maps["leaflet"] = tile_layer;		
		break;
		
	    case "protomaps":		    
		
		var tile_url = cfg.tile_url;
		
		var tile_layer = protomapsL.leafletLayer({
		    url: tile_url,
		    theme: cfg.protomaps.theme,
		})

		tile_layer.addTo(map);

		base_maps["protomaps"] = tile_layer;		
		break;
		
	    default:
		console.error("Uknown or unsupported map provider");
		return;
	}

	console.log("Add overlays", overlays);
	
	var layerControl = L.control.layers(base_maps, overlays);
	layerControl.addTo(map);
	
    };
    
    fetch("/map.json")
	.then((rsp) => rsp.json())
	.then((cfg) => {	    
	    init(cfg);
	}).catch((err) => {
	    console.error("Failed to retrieve map config", err);
	});
        
    
});
