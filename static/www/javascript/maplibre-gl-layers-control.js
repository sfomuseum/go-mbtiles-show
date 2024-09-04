// https://blog.wxm.be/2024/01/24/maplibre-layer-visibility-control.html

class LayersControl {
    constructor(ctrls) {
    // This div will hold all the checkboxes and their labels
    this._container = document.createElement("div");
    this._container.classList.add(
      // Built-in classes for consistency
      "maplibregl-ctrl",
      "maplibregl-ctrl-group",
      // Custom class, see later
      "layers-control",
    );
    // Might be cleaner to deep copy these instead
    this._ctrls = ctrls;
    // Direct access to the input elements so I can decide which should be
    // checked when adding the control to the map.
    this._inputs = [];
    // Create the checkboxes and add them to the container
    for (const key of Object.keys(this._ctrls)) {
	let labeled_checkbox = this._createLabeledCheckbox(key);
      this._container.appendChild(labeled_checkbox);
    }
  }

  // Creates one checkbox and its label
  _createLabeledCheckbox(key) {
    let label = document.createElement("label");
    label.classList.add("layer-control");
    let text = document.createTextNode(key);
    let input = document.createElement("input");
    this._inputs.push(input);
    input.type = "checkbox";
    input.id = key;
    // `=>` function syntax keeps `this` to the LayersControl object
    // When changed, toggle all the layers associated with the checkbox via
      // `this._ctrls`.

    input.addEventListener("change", () => {
      let visibility = input.checked ? "visible" : "none";
	for (const layer of this._ctrls[input.id]) {
            this._map.setLayoutProperty(layer, "visibility", visibility);
	}
    });
    label.appendChild(input);
    label.appendChild(text);
    return label;
  }

  onAdd(map) {
    this._map = map;
    // For every checkbox, find out if all its associated layers are visible.
      // Check the box if so.
      for (const input of this._inputs) {
      // List of all layer ids associated with this checkbox
	  let layers = this._ctrls[input.id];
      // Check whether every layer is currently visible
      let is_visible = true;
	  for (const layername of layers) {
        is_visible =
          is_visible &&
          this._map.getLayoutProperty(layername, "visibility") !== "none";
      }
      input.checked = is_visible;
    }
    return this._container;
  }

  onRemove(map) {
    // Not sure why we have to do this ourselves since we are not the ones
    // adding us to the map.
    // Copied from their example so keeping it in.
    this._container.parentNode.removeChild(this._container);
    // This might be to help garbage collection? Also from their example.
    // Or perhaps to ensure calls to this object do not change the map still
    // after removal.
    this._map = undefined;
  }
}
