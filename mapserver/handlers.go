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

type atlasWebhook struct {
	Content string `json:"content"`
}

func (s *MapServer) webhookHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t atlasWebhook
	err := decoder.Decode(&t)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//lines := strings.Split(t.Content, "\n")

	log.Println(t.Content)
}
