package mapserver

import (
	"encoding/json"
	"log"
	"net/http"
)

func (s *MapServer) getTerritoryURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	wrapper := make(map[string]string)
	wrapper["url"] = s.territoryURL

	js, err := json.Marshal(wrapper)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(js)
}
