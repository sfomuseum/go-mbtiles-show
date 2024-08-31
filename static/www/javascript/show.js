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

	var base_maps = {};
	var overlays = {};

	if (cfg.raster_layers){
	    
	    for (k in cfg.raster_layers){
		var l = L.tileLayer(cfg.raster_layers[k])
		overlays[k] = l;
	    }
	}

	if (cfg.vector_layers){
	    console.log("Vector layers not supported yet.")
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
