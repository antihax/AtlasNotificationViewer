package mapserver

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
)

func (s *MapServer) listUsers(w http.ResponseWriter, r *http.Request) {
	v, err := s.scanHash("user:*", 2000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (s *MapServer) changeUser(w http.ResponseWriter, r *http.Request) {

	steamID := r.FormValue("steamID")
	change := r.FormValue("change")
	value := r.FormValue("value")

	testNumeric := regexp.MustCompile("^[0-9]+$")

	// Validate user input
	switch change {
	case "allowed":
	case "admin":
	default:
		http.Error(w, "bad input", http.StatusBadRequest)
		return
	}

	switch value {
	case "0":
	case "1":
	default:
		http.Error(w, "bad input", http.StatusBadRequest)
		return
	}

	if !testNumeric.MatchString(steamID) { // Also tests empty
		http.Error(w, "bad input", http.StatusBadRequest)
		return
	}

	_, err := s.redisClient.HSet(context.Background(), "user:"+steamID, change, value).Result()
	if err != nil {
		http.Error(w, "bad input", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(1)
}
