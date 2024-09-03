window.addEventListener("load", function load(event){

    // Null Island
    // var lat = 0.0;
    // var lon = 0.0;
    // var zoom = 3;
    
    // SFO 
    var lat = 37.621131;
    var lon = -122.384292;
    var zoom = 12;
    
    var map = L.map('map');
    map.setView([lat, lon], zoom);
    
    var init = function(cfg){

	console.log("INIT", cfg);
	
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

	    var tiles_styles = {
		
		all: function(properties, zoom) {
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
	    
	    for (k in cfg.vector_layers){

		var layer_url = cfg.vector_layers[k]

		var layer_opts = {
		    rendererFactory: L.canvas.tile,
		    vectorTileLayerStyles: tiles_styles,
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
