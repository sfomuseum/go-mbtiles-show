window.addEventListener("load", function load(event){

    // Null Island
    // var lat = 0.0;
    // var lon = 0.0;
    // var zoom = 3;
    
    // SFO 
    var lat = 37.621131;
    var lon = -122.384292;
    var zoom = 14;
    
    const init = function(local_cfg, map_cfg){

	switch (map_cfg.provider) {
	    case "leaflet":
		init_leaflet(local_cfg, map_cfg);
		break;
	    default:		
		init_maplibre(local_cfg, map_cfg);
		break;
	}
    }

    const init_maplibre = function(local_cfg, map_cfg){

	console.debug("Initialize maplibre w/ base tile url", map_cfg.tile_url);
	console.debug("Initialize maplibre w/ local config", local_cfg);	
	
	var base_souce = {};
	var base_layer = {};

	if (map_cfg.protomaps && map_cfg.protomaps.use_pmtiles) {

	    // add the PMTiles plugin to the maplibregl global.
	    // https://maplibre.org/maplibre-gl-js/docs/examples/pmtiles/
	    // https://github.com/protomaps/PMTiles/blob/main/js/examples/maplibre.html
	    // https://unpkg.com/pmtiles@3.0.7/dist/pmtiles.js

	    if (! map_cfg.tile_url.startsWith("http")){
		map_cfg.tile_url = "http://" + location.host + map_cfg.tile_url;
	    }
	    
	    const protocol = new pmtiles.Protocol();
	    maplibregl.addProtocol('pmtiles', protocol.tile);
	    
	    const p = new pmtiles.PMTiles(map_cfg.tile_url);
	    protocol.add(p);
	    
	    base_source = {
		type: "vector",
		url: "pmtiles://" + map_cfg.tile_url,
	    };
	    
	    base_layer = {
		'id': 'base',
		'source': 'base',
		// I wish there were a way to specify "all the layers" ...
		'source-layer': 'roads',
		'type': "line",
                'paint': {
                    "line-color": "#fc8d62",
                }
	    };
	    
	} else {

	    
	    base_source = {
		type: 'raster',
		tiles: [
		    map_cfg.tile_url,
		],
		'tileSize': 256,
	    };
	    
	    base_layer = {
		'id': 'base',
		'type': 'raster',
		'source': 'base',
	    };
	    
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
	    
	    if (local_cfg.raster_layers){
		
		// Basically inverted-y coordinates ({-y}) are not supported in maplibre-gl.js
		// https://maplibre.org/maplibre-style-spec/sources/#raster
		// https://maplibre.org/maplibre-gl-js/docs/API/type-aliases/CanvasSourceSpecification/
		// Despite seemingly being supported in the "native" builds...
		// https://docs.mapbox.com/ios/maps/api/6.4.1/tile-url-templates.html
		// https://maplibre.org/maplibre-native/docs/book/design/coordinate-system.html
		    
		for (k in local_cfg.raster_layers){
		    
		    var tile_url = "http://" + location.host + local_cfg.raster_layers[k];
		    console.debug("Add raster layer", k, tile_url);
		    
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
	    
	    if (local_cfg.vector_layers){
		
		for (k in local_cfg.vector_layers){

		    var tile_url = "http://" + location.host + local_cfg.vector_layers[k];
		    console.debug("Add vector layer", k, tile_url);
		    
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

    const init_leaflet = function(local_cfg, map_cfg){

	var map = L.map('map');
	map.setView([lat, lon], zoom);
    
	var base_maps = {};
	var overlays = {};

	if (local_cfg.raster_layers){
	    
	    for (k in local_cfg.raster_layers){
		var l = L.tileLayer(local_cfg.raster_layers[k])
		overlays[k] = l;
	    }
	}

	if (local_cfg.vector_layers){
	    console.warn("Vector layers not supported yet.")
	}

	if (map_cfg.protomaps && map_cfg.protomaps.use_pmtiles) {

	    var tile_url = map_cfg.tile_url;

	    var pm_args = {
		url: tile_url,
		theme: map_cfg.protomaps.theme,
		flavor: map_cfg.protomaps.theme,
	    };
	    
	    if ("max_data_zoom" in map_cfg.protomaps){
		pm_args.maxDataZoom = map_cfg.protomaps.max_data_zoom;
	    }
	    
	    var tile_layer = protomapsL.leafletLayer(pm_args);
	    
	    tile_layer.addTo(map);
	    base_maps["protomaps"] = tile_layer;
	    
	} else {
		
	    var tile_url = map_cfg.tile_url;
		
	    var tile_layer = L.tileLayer(tile_url);
	    tile_layer.addTo(map);
	    
	    base_maps["leaflet"] = tile_layer;		
	}
		
	var layerControl = L.control.layers(base_maps, overlays);
	layerControl.addTo(map);	
    };

    console.debug("Fetch local config");
    
    fetch("/config.json")
	.then(rsp => rsp.json())
	.then((local_cfg) => {

	    console.debug("Fetch map config");
	    
	    fetch("/map.json")
		.then((rsp) => rsp.json())
		.then((map_cfg) => {	    
		    init(local_cfg, map_cfg);
		}).catch((err) => {
		    console.error("Failed to retrieve map config", err);
		});
        }).catch((err) => {
	    console.error("Failed to retrieve local config", err);
	});
    
});
