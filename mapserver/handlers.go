package atlasadminserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

func (s *AtlasAdminServer) getTerritoryURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	wrapper := make(map[string]string)
	wrapper["url"] = s.config.TerritoryURL

	js, err := json.Marshal(wrapper)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(js)
}

func (s *AtlasAdminServer) getTribes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	s.tribeDataLock.RLock()
	js, err := json.Marshal(s.tribeData)
	s.tribeDataLock.RUnlock()

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(js)
}

// getData returns the latest game data pulled from the backend.
func (s *AtlasAdminServer) getData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	s.gameDataLock.RLock()
	js, err := json.Marshal(s.gameData)
	s.gameDataLock.RUnlock()

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(js)
}

// getPathTravelled returns a fake path for the specified ship.
func (s *AtlasAdminServer) getPathTravelled(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	shipID := r.URL.Query().Get("id")

	s.gameDataLock.RLock()
	info, ok := s.gameData[shipID]
	s.gameDataLock.RUnlock()

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var path []ServerLocation
	path = append(path, ServerLocation{
		info.ServerID,
		info.ServerXRelativeLocation,
		info.ServerYRelativeLocation,
	})

	for i, end := 0, rand.Intn(10); i <= end; i++ {
		path = append(path, ServerLocation{
			info.ServerID,
			path[len(path)-1].ServerXRelativeLocation * rand.Float64(),
			path[len(path)-1].ServerYRelativeLocation * rand.Float64(),
		})
	}

	js, err := json.Marshal(path)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// sendCommand publishes an event to the GeneralNotifications:GlobalCommands
// redis PubSub channel. To send a server command, prepend "ID::X,Y::" where
// ID is the packed server ID; X and Y are the relative lng and lat locations.
func (s *AtlasAdminServer) sendCommand(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" || s.config.DisableCommands {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

	encoded := fmt.Sprintf("%s", body)

	log.Println("publish:", encoded)
	result, err := s.redisClient.Publish("GeneralNotifications:GlobalCommands", encoded).Result()
	if err != nil {
		log.Println("redis error for: ", encoded, "; ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(strconv.FormatInt(result, 10)))
}
