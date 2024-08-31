window.addEventListener("load", function load(event){

    // Null Island
    // var lat = 0.0;
    // var lon = 0.0;
    // var zoom = 3;
    
    // SFO 
    var lat = 37.621131;
    var lon = -122.384292;
    var zoom = 12;
    
    var map = L.map('map').setView([lat, lon], zoom);
    
    var init = function(cfg){

	console.log(cfg);

	var base_maps = {};
	var overlays = {};

	for (k in cfg.mbtiles_layers){
	    var l = L.tileLayer(cfg.mbtiles_layers[k])
	    overlays[k] = l;
	}
	
	switch (cfg.provider) {
	    case "leaflet":
		
		var tile_url = cfg.tile_url;
		
		var tile_layer = L.tileLayer(tile_url, {
		    maxZoom: 19,
		});

		base_maps["leaflet"] = tile_layer;
		tile_layer.addTo(map);
		break;
		
	    case "protomaps":		    
		
		var tile_url = cfg.tile_url;
		
		var tile_layer = protomapsL.leafletLayer({
		    url: tile_url,
		    theme: cfg.protomaps.theme,
		})

		base_maps["protomaps"] = tile_layer;
		tile_layer.addTo(map);
		break;
		
	    default:
		console.error("Uknown or unsupported map provider");
		return;
	}

	console.log("BASE", base_maps);
	console.log("OVERLAYS", overlays);
	
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
