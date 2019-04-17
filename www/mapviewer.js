if (!window) {
  function require() { }
  const
    React = require("react"),
    ReactDOM = require("react-dom"),
    L = require("leaflet")
}

const colors = (
  "red green yellow blue orange purple cyan magenta lime pink teal lavender brown beige maroon olive coral navy"
).split(" ")

class WorldMap extends React.Component {
  constructor(props) {
    super(props)
    this.forceTileReload = this.forceTileReload.bind(this)
  }

  forceTileReload() {
    fetch("territoryURL")
      .then(res => res.json())
      .then(config => {
        if (config.url) {
          if (this.territoryLayer) {
            this.worldMap.removeLayer(this.territoryLayer)
            delete this.territoryLayer
          }

          this.territoryLayer = L.tileLayer(config.url + "{z}/{x}/{y}.png?t={cachebuster}", {
            maxZoom: 9,
            maxNativeZoom: 6,
            minZoom: 1,
            bounds: L.latLngBounds([0, 0], [-256, 256]),
            noWrap: true,
            cachebuster: function () { return Math.random(); }
          })

          this.territoryLayer.addTo(this.worldMap)
        } else {
          console.error("Did not receive territory URL")
        }
      })
      .catch((err) => {
        this.setState({
          notification: {
            type: "error",
            msg: "Failed to get territory URL from server",
          }
        })
      })
  }

  componentWillUnmount() {
    clearInterval(this.timer);
  }

  componentDidMount() {

    this.forceTileReload()
    this.timer = setInterval(this.forceTileReload, 15000)

    const layerOpts = {
      maxZoom: 9,
      maxNativeZoom: 6,
      minZoom: 1,
      bounds: L.latLngBounds([0, 0], [-256, 256]),
      noWrap: true,
    };

    const baseLayer = L.tileLayer("tiles/background/{z}/{x}/{y}.png", layerOpts)

    const map = this.worldMap = L.map("worldmap", {
      crs: L.CRS.Simple,
      layers: [baseLayer],
      zoomControl: false,
    })

    // Remove features after 5 minutes
    setInterval(function () {
      map.eachLayer(function (layer) {
        if (layer.feature && new Date - layer.feature.properties.created > 5 * 60 * 1000) {
          map.removeLayer(layer)
        }
      });
    }, 1000);

    // Add zoom control
    L.control.zoom({
      position: 'topright'
    }).addTo(map);

    // Add Layer Control
    L.control.layers({}, {
      Islands: L.tileLayer("tiles/islands/{z}/{x}/{y}.png", layerOpts).addTo(map),
      Discoveries: L.tileLayer("tiles/disco/{z}/{x}/{y}.png", layerOpts),
      Names: L.tileLayer("tiles/names/{z}/{x}/{y}.png", layerOpts),
      Grid: L.tileLayer("tiles/grid/{z}/{x}/{y}.png", layerOpts).addTo(map)
    }, { position: 'topright' }).addTo(map)

    fetch("/account")
      .then(response => {
        if (response.ok) {
          return response.json();
        } else {
          L.easyButton('<img src="/sits_02.png">', function (btn, map) {
            window.location.href = '/login';
          }).addTo(map);
        }
      }).then(d => {
        if (d) {
          L.easyButton('fa-sign-out', function (btn, map) {
            window.location.href = '/logout';
          }).addTo(map);
          if (d.allowed === "1") {
            connect();
          }
          if (d.admin === "1") {
            L.easyButton('fa-globe', function (btn, map) {
              window.location.href = '/admin';
            }).addTo(map);
          }
        }
      })
      .catch(error => { console.log(error) });

    map.setView([-128, 128], 2)

    if (this.props.onContextMenu)
      map.on("contextmenu.show", this.props.onContextMenu)
    if (this.props.onContextMenuClose)
      map.on("contextmenu.hide", this.props.onContextMenuClose)

    function connect() {
      var ws = new WebSocket(getLocalURI() + '/events');
      ws.onopen = function () {

      };
      ws.onclose = function () {
        setTimeout(function () { connect(); }, 1000);
      };
      ws.onmessage = function (event) {
        event.data.split(/\n/).forEach(function (d) {
          if (d == "") return;

          var f = JSON.parse(d);
          let coordReg = /Long: ([0-9.]+) \/ Lat: ([0-9.]+)/;
          let coords = f.content.match(coordReg);
          if (coords) {
            let n = [scaleAtlasToLeaflet(parseFloat(coords[1])), -scaleAtlasToLeaflet(-parseFloat(coords[2]))];
            console.log(f.content)
            gj.addData({
              "type": "Feature",
              "properties": {
                "name": f.tribe,
                "popupContent": f.content,
                "created": new Date()
              },
              "geometry": {
                "type": "Point",
                "coordinates": n
              }
            }).addTo(map)
          }
        });
      };
    }

    var gj = L.geoJson(false, {
      onEachFeature: function (feature, layer) {
        layer.bindPopup(feature.properties.popupContent);
      },
      pointToLayer: function (feature, latlng) {
        return L.circleMarker(latlng, {
          radius: 8,
          fillColor: "#ff4400",
          color: "#000",
          weight: 1,
          opacity: 1,
          fillOpacity: 0.8
        });
      }
    }).addTo(map);

    L.Control.MousePosition = L.Control.extend({
      options: {
        position: 'bottomleft',
        separator: ' : ',
        emptyString: 'Unavailable',
        lngFirst: false,
        numDigits: 5,
        lngFormatter: undefined,
        latFormatter: undefined,
        prefix: ""
      },

      onAdd: function (map) {
        this._container = L.DomUtil.create('div', 'leaflet-control-mouseposition');
        L.DomEvent.disableClickPropagation(this._container);
        map.on('mousemove', this._onMouseMove, this);
        this._container.innerHTML = this.options.emptyString;
        return this._container;
      },

      onRemove: function (map) {
        map.off('mousemove', this._onMouseMove)
      },

      _onMouseMove: function (e) {
        var lng = this.options.lngFormatter ? this.options.lngFormatter(e.latlng.lng) : L.Util.formatNum(e.latlng.lng, this.options.numDigits);
        var lat = this.options.latFormatter ? this.options.latFormatter(e.latlng.lat) : L.Util.formatNum(e.latlng.lat, this.options.numDigits);
        var value = this.options.lngFirst ? lng + this.options.separator + lat : lat + this.options.separator + lng;
        var prefixAndValue = this.options.prefix + ' ' + value;
        this._container.innerHTML = prefixAndValue;
      }
    });

    L.Map.mergeOptions({
      positionControl: false
    });

    L.Map.addInitHook(function () {
      if (this.options.positionControl) {
        this.positionControl = new L.Control.MousePosition();
        this.addControl(this.positionControl);
      }
    });

    L.control.mousePosition = function (options) {
      return new L.Control.MousePosition(options);
    };
    L.control.mousePosition().addTo(map);
  }

  render() {
    return (
      <div id="worldmap"></div>
    )
  }
}


function isTribeID(tribeID) {
  return tribeID > 1000000000 + 50000
}

function getTribeColor(tribeID) {
  if (!tribeID)
    return "black"
  if (!isTribeID(tribeID))
    return "grey"

  var idx = tribeID % colors.length
  return colors[idx]
}

class App extends React.Component {
  constructor(props) {
    super(props)

    this.state = {
      notification: {},
      entities: {},
      tribes: {},
      activeTribeColor: "",
      sending: false,
    }
  }

  render() {
    const {
      activeTribeColor,
      notification
    } = this.state
    return (
      <div className="App">
        <WorldMap
          color={activeTribeColor}
        />
        <div className={"notification " + (notification.type || "hidden")}>
          {notification.msg}
          <button className="close" onClick={() => this.setState({ notification: {} })}>Dismiss</button>
        </div>
      </div>
    )
  }
}

function scaleAtlasToLeaflet(e) {
  return (e + 100) * (1.28);
}
// Get local URI for requests
function getLocalURI() {
  var loc = window.location, new_uri;
  if (loc.protocol === "https:") {
    new_uri = "wss:";
  } else {
    new_uri = "ws:";
  }
  new_uri += "//" + loc.host;
  return new_uri;
}

ReactDOM.render(
  <App refresh={5 * 1000 /* 5 seconds */} />,
  document.getElementById("app")
)
