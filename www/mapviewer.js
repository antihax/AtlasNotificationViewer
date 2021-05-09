if (!window) {
  function require() { }
  const
    React = require("react"),
    ReactDOM = require("react-dom"),
    L = require("leaflet");
}

class WorldMap extends React.Component {
  constructor(props) {
    super(props)
  }


  componentWillUnmount() {
    clearInterval(this.timer);
  }

  componentDidMount() {
    const layerOpts = {
      maxZoom: 9,
      maxNativeZoom: 6,
      minZoom: 1,
      bounds: L.latLngBounds([0, 0], [-256, 256]),
      noWrap: true,
    };

    const baseLayer = L.tileLayer("tiles/tiles/{z}/{x}/{y}.png", layerOpts)

    const map = this.worldMap = L.map("worldmap", {
      crs: L.CRS.Simple,
      layers: [baseLayer],
      zoomControl: false,
      attributionControl: false,
    })
    map._originalBounds = layerOpts.bounds;
    // Remove features after 5 minutes
    setInterval(function () {
      map.eachLayer(function (layer) {
        if (layer.feature && new Date - layer.feature.properties.created > 15 * 60 * 1000) {
          map.removeLayer(layer)
        }
      });
    }, 1000);

    // Add zoom control
    L.control.zoom({
      position: 'topright'
    }).addTo(map);

    map.IslandTerritories = L.layerGroup(layerOpts);
    map.IslandResources = L.layerGroup(layerOpts);

    // Add Layer Control
    L.control.layers({}, {
      Grid: new L.AtlasGrid({
        xticks: config.ServersX,
        yticks: config.ServersY
      }).addTo(map),
      Territories: map.IslandTerritories.addTo(map),
    }, { position: 'topright' }).addTo(map);

    fetch("/account")
      .then(response => {
        if (response.ok) {
          return response.json();
        } else {
          L.easyButton('fab fa-steam-symbol fa-lg', function (btn, map) {
            window.location.href = '/login';
          }).addTo(map);
        }
      }).then(d => {
        if (d) {
          L.easyButton('fa fa-times fa-lg', function (btn, map) {
            window.location.href = '/logout';
          }).addTo(map);
          if (d.allowed === "1") {
            connect();
          }
          if (d.admin === "1") {
            L.easyButton('fa fa-globe fa-lg', function (btn, map) {
              window.location.href = '/admin/';
            }).addTo(map);
          }
        }
      })
      .catch(error => { console.log(error) });
    map.setView([-128, 128], 2)

    var createIslandLabel = function (island) {
      var label = "";
      label += '<div id="island_' + island.IslandID + '" class="islandlabel">';
      label += '<div class="islandlabel_icon"><img class="islandlabel_size" src="' + getIslandIcon(island) + '" width="32" height="32"/></div>';
      label += '</div>'
      return L.divIcon({
        className: "islandlabel",
        html: label
      })
    }
    var createLabelIcon = function (labelClass, labelText) {
      return L.divIcon({
        className: labelClass,
        html: labelText
      })
    }



    fetch('/islands', {
      dataType: 'json'
    })
      .then(res => res.json())
      .then(function (IslandDataJson) {
        var IslandEntries = IslandDataJson.Islands;
        var CompanyHashMap = IslandDataJson.Companies.reduce(function (map, obj) {
          map[obj.TribeId] = obj;
          return map;
        }, {});
        GlobalPriortyQueue.clear();
        var now = Math.floor(Date.now() / 1000);
        for (var j = 0; j < IslandEntries.length; j++) {
          var Island = IslandEntries[j];
          getWarState(Island);
          getPeaceState(Island);
          var nextUpdate = Island.CombatNextUpdateSec;
          if (Island.WarNextUpdateSec < nextUpdate)
            nextUpdate = Island.WarNextUpdateSec;
          GlobalPriortyQueue.enqueue(Island, now + nextUpdate + 1);

          var OwningTribe = CompanyHashMap[Island.TribeId];
          if (OwningTribe) {
            var circle = new IslandCircle([-256 * Island.Y, 256 * Island.X], {
              radius: Island.Size * 256,
              color: Island.Color,
              opacity: 0,
              fillColor: Island.Color,
              fillOpacity: 0.5
            });
            var PopupHTML = '';
            circle.Island = Island;
            if (OwningTribe.FlagURL) {
              PopupHTML = '<p><img border="0" alt="CompanyFlag" src="' + OwningTribe.FlagURL + '" width="100" height="100"></p>';
            }
            PopupHTML += '<strong>' + escapeHTML(Island.SettlementName) + '</strong> <sup>[' + Island.IslandPoints + ' pts]</sup>'
            PopupHTML += '<div style="width: 250px;" id="pop_up_war">---</div>';
            PopupHTML += '<div style="width: 250px;" id="pop_up_phase">---</div>';
            PopupHTML += 'Owner: ' + escapeHTML(OwningTribe.TribeName);
            if (Island.NumSettlers >= 0) {
              PopupHTML += '<br>Settlers: ' + Island.NumSettlers;
            }
            PopupHTML += '<br>Taxation: ' + Island.TaxRate.toFixed(1) + '%';
            circle.bindPopup(PopupHTML, {
              showOnMouseOver: true
            });
            map.IslandTerritories.addLayer(circle);
          }
        }
      });

    var urlParams = new URLSearchParams(window.location.search);
    if (urlParams.get('lat')) {
      let n = [-scaleAtlasToLeaflet(-parseFloat(urlParams.get('lat'))), scaleAtlasToLeaflet(parseFloat(urlParams.get('long')))];
      L.circleMarker(n, {
        "radius": 30,
        "color": '#FF0000',
      }).addTo(map);
      L.circleMarker(n, {
        "radius": 3,
        "color": '#000',
        "weight": 4,
        "fill": true,
        "fillColor": '#FFFFFF',
        "fillOpacity": 1
      }).addTo(map);
      map.setView(n, 8)
    }

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
          let coordReg = /Long: ([0-9.-]+) \/ Lat: ([0-9.-]+)/;
          let coords = f.content.match(coordReg);

          if (coords) {
            let n = [scaleAtlasToLeaflet(parseFloat(coords[1])), -scaleAtlasToLeaflet(-parseFloat(coords[2]))];
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

          new Noty({
            text: f.content,
            timeout: 10000,
            layout: 'bottomRight',
            theme: 'light'
          }).show();

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
        var lng = L.Util.formatNum(scaleLeafletToAtlas(e.latlng.lng), 2);
        var lat = -L.Util.formatNum(scaleLeafletToAtlas(-e.latlng.lat), 2);
        var value = lng + this.options.separator + lat + " / " + e.latlng.lng + this.options.separator + e.latlng.lat;
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

class App extends React.Component {
  constructor(props) {
    super(props)

    this.state = {
      notification: {},
      entities: {},
      tribes: {},
      sending: false,
    }
  }

  render() {
    const {
      notification
    } = this.state
    return (
      <div className="App">
        <WorldMap />
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

function scaleLeafletToAtlas(e) {
  return (e / 1.28);
}

function GPStoLeaflet(x, y) {
  var long = (y - 100) * 1.28,
    lat = (100 + x) * 1.28;

  return [long, lat];
}

function unrealToLeaflet(x, y) {
  const unreal = 15400000;
  var lat = ((x / unreal) * 256),
    long = -((y / unreal) * 256);
  return [long, lat];
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

class QElement {
  constructor(element, priority) {
    this.element = element;
    this.priority = priority;
  }
}
class PriorityQueue {
  constructor() {
    this.items = [];
  }
  enqueue(element, priority) {
    var qElement = new QElement(element, priority);
    var contain = false;
    for (var i = 0; i < this.items.length; i++) {
      if (this.items[i].priority > qElement.priority) {
        this.items.splice(i, 0, qElement);
        contain = true;
        break;
      }
    }
    if (!contain) {
      this.items.push(qElement);
    }
  }
  dequeue() {
    return this.items.shift();
  }
  front() {
    return this.items[0];
  }
  isEmpty() {
    return this.items.length == 0;
  }
  clear() {
    this.items = []
  }
}

function formatSeconds(InTime) {
  var Days = 0
  var Hours = Math.floor(InTime / 3600);
  var Minutes = Math.floor((InTime % 3600) / 60);
  var Seconds = Math.floor((InTime % 3600) % 60);
  if (Hours >= 24) {
    Days = Math.floor(Hours / 24);
    Hours = Hours - (Days * 24)
  }
  if (Days > 0)
    return Days + "d:" + Hours + "h:" + Minutes + "n:" + Seconds + "s";
  else if (Hours > 0)
    return Hours + "h:" + Minutes + "m:" + Seconds + "s";
  else if (Minutes > 0)
    return Minutes + "m:" + Seconds + "s";
  else
    return Seconds + "s";
}

function getWarState(Island) {
  var now = Math.floor(Date.now() / 1000)
  if (now >= Island.WarStartUTC && now < Island.WarEndUTC) {
    Island.bWar = true;
    Island.WarNextUpdateSec = Island.WarEndUTC - now;
    return "AT WAR! ENDS IN " + formatSeconds(Island.WarNextUpdateSec)
  } else if (now < Island.WarStartUTC) {
    Island.bWar = false;
    Island.WarNextUpdateSec = Island.WarStartUTC - now;
    return "WAR BEGINS IN " + formatSeconds(Island.WarNextUpdateSec)
  } else if (now < Island.WarEndUTC + 5 * 24 * 3600) {
    Island.bWar = false;
    Island.WarNextUpdateSec = Island.WarEndUTC + 5 * 24 * 3600 - now;
    return "CAN DECLARE WAR IN " + formatSeconds(Island.WarNextUpdateSec)
  } else {
    Island.bWar = false;
    Island.WarNextUpdateSec = Number.MAX_SAFE_INTEGER;
    return "War can be declared on this settlement."
  }
}

function getPeaceState(Island) {
  var now = new Date();
  var CombatStartSeconds = Island.CombatPhaseStartTime;
  var CombatEndSeconds = (CombatStartSeconds + 32400) % 86400;
  var CurrentDaySeconds = (3600 * now.getUTCHours()) + (60 * now.getUTCMinutes()) + now.getUTCSeconds();
  if (CombatEndSeconds > CombatStartSeconds) {
    if (CurrentDaySeconds < CombatStartSeconds) {
      Island.bCombat = false;
      Island.CombatNextUpdateSec = CombatStartSeconds - CurrentDaySeconds;
      return "In Peace Phase. " + formatSeconds(Island.CombatNextUpdateSec) + " remaining"
    } else if (CurrentDaySeconds >= CombatStartSeconds && CurrentDaySeconds < CombatEndSeconds) {
      Island.bCombat = true;
      Island.CombatNextUpdateSec = CombatEndSeconds - CurrentDaySeconds;
      return "In Combat Phase! " + formatSeconds(Island.CombatNextUpdateSec) + " remaining"
    } else {
      Island.bCombat = false;
      Island.CombatNextUpdateSec = 86400 - CurrentDaySeconds + CombatStartSeconds
      return "In Peace Phase." + formatSeconds(Island.CombatNextUpdateSec) + " remaining"
    }
  } else {
    if (CurrentDaySeconds >= CombatStartSeconds) {
      Island.bCombat = true;
      Island.CombatNextUpdateSec = 86400 - CurrentDaySeconds + CombatEndSeconds;
      return "In Combat Phase! " + formatSeconds(Island.CombatNextUpdateSec) + " remaining"
    } else if (CurrentDaySeconds < CombatEndSeconds) {
      Island.bCombat = true;
      Island.CombatNextUpdateSec = CombatEndSeconds - CurrentDaySeconds;
      return "In Combat Phase! " + formatSeconds(Island.CombatNextUpdateSec) + " remaining"
    } else {
      Island.bCombat = false;
      Island.CombatNextUpdateSec = CombatStartSeconds - CurrentDaySeconds;
      return "In Peace Phase. " + formatSeconds(Island.CombatNextUpdateSec) + " remaining"
    }
  }
}

function getIslandIcon(Island) {
  if (Island.bWar || Island.bCombat)
    return "HUD_War_Icon.png";
  else
    return "HUD_Peace_Icon.png";
}
var GlobalSelectedIsland = null;
var GlobalPriortyQueue = new PriorityQueue();
setInterval(updateIsland, 1000)

function updateIsland() {
  while (!GlobalPriortyQueue.isEmpty()) {
    var now = Math.floor(Date.now() / 1000);
    if (GlobalPriortyQueue.front().priority > now)
      break;
    var Island = GlobalPriortyQueue.dequeue().element;
    getWarState(Island);
    getPeaceState(Island);
    var el = document.getElementById("island_" + Island.IslandID);
    if (el != null) {
      var img = el.getElementsByClassName("islandlabel_size")[0];
      img.src = getIslandIcon(Island);
    }
    var nextUpdate = Island.CombatNextUpdateSec;
    if (Island.WarNextUpdateSec < nextUpdate)
      nextUpdate = Island.WarNextUpdateSec;
    GlobalPriortyQueue.enqueue(Island, now + nextUpdate + 1);
  }
  if (GlobalSelectedIsland != null) {
    var phase = document.getElementById("pop_up_phase")
    if (phase != null)
      phase.innerHTML = getPeaceState(GlobalSelectedIsland)
    var war = document.getElementById("pop_up_war")
    if (war != null)
      war.innerHTML = getWarState(GlobalSelectedIsland)
  }
}
class IslandCircle extends L.Circle {
  constructor(latlng, options) {
    super(latlng, options)
    this.Island = null
    this.bindPopup = this.bindPopup.bind(this)
    this._popupMouseOut = this._popupMouseOut.bind(this)
    this._getParent = this._getParent.bind(this)
  }
  bindPopup(htmlContent, options) {
    if (options && options.showOnMouseOver) {
      L.Marker.prototype.bindPopup.apply(this, [htmlContent, options]);
      this.off("click", this.openPopup, this);
      this.on("mouseover", function (e) {
        var target = e.originalEvent.fromElement || e.originalEvent.relatedTarget;
        var parent = this._getParent(target, "leaflet-popup");
        if (parent == this._popup._container)
          return true;
        GlobalSelectedIsland = this.Island
        this.openPopup();
      }, this);
      this.on("mouseout", function (e) {
        var target = e.originalEvent.toElement || e.originalEvent.relatedTarget;
        if (this._getParent(target, "leaflet-popup")) {
          L.DomEvent.on(this._popup._container, "mouseout", this._popupMouseOut, this);
          return true;
        }
        this.closePopup();
        GlobalSelectedIsland = null
      }, this);
    }
  }
  _popupMouseOut(e) {
    L.DomEvent.off(this._popup, "mouseout", this._popupMouseOut, this);
    var target = e.toElement || e.relatedTarget;
    if (this._getParent(target, "leaflet-popup"))
      return true;
    if (target == this._path)
      return true;
    this.closePopup();
    GlobalSelectedIsland = null;
  }
  _getParent(element, className) {
    if (element == null)
      return false;
    var parent = element.parentNode;
    while (parent != null) {
      if (parent.className && L.DomUtil.hasClass(parent, className))
        return parent;
      parent = parent.parentNode;
    }
    return false;
  }
}

function escapeHTML(unsafe_str) {
  return unsafe_str.replace(/&/g, '&').replace(/</g, '<').replace(/>/g, '>').replace(/\"/g, '"').replace(/\'/g, '\'');
}
