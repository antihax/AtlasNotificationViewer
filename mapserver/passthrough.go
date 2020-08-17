package mapserver

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func (s *MapServer) loadPassthrough() error {
	v, err := s.scanHash("passthrough:*", 2000)
	if err != nil {
		return err
	}

	s.passthroughLock.Lock()
	for id, keys := range v {
		for field, value := range keys {
			if s.passthrough[id] == nil {
				s.passthrough[id] = make(map[string]string)
			}
			s.passthrough[id][field] = value
		}
	}

	s.passthroughLock.Unlock()
	return nil
}

func (s *MapServer) listEventTypes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(eventTypes)
}

func (s *MapServer) listPassthrough(w http.ResponseWriter, r *http.Request) {
	v, err := s.scanHash("passthrough:*", 2000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (s *MapServer) changePassthroughMap(id string, field string, value string) {
	s.passthroughLock.Lock()
	if s.passthrough[id] == nil {
		s.passthrough[id] = make(map[string]string)
	}
	s.passthrough[id][field] = value
	s.passthroughLock.Unlock()
}

func (s *MapServer) changePassthrough(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	change := r.FormValue("change")
	value := r.FormValue("value")

	testNumeric := regexp.MustCompile("^[0-9]+$")

	// Validate user input
	found := false
	for _, cT := range eventTypes {
		if change == cT {
			found = true
			break
		}
	}
	if !found && change != "url" {
		http.Error(w, "bad input", http.StatusBadRequest)
		log.Printf("Type not found")
		return
	}

	if !testNumeric.MatchString(id) { // Also tests empty
		http.Error(w, "bad input", http.StatusBadRequest)
		return
	}

	_, err := s.redisClient.HSet(context.Background(), "passthrough:"+id, change, value).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.changePassthroughMap(id, change, value)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(1)
}

func (s *MapServer) addPassthrough(w http.ResponseWriter, r *http.Request) {
	hash, err := s.scanHash("passthrough:*", 1000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	max := 1
	for key := range hash {
		split := strings.Split(key, ":")
		id, err := strconv.Atoi(split[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if max < id {
			max = id
		}
	}

	id := strconv.Itoa(max)
	_, err = s.redisClient.HSet(context.Background(), "passthrough:"+id, "id", id).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.changePassthroughMap(id, "id", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(id)
}
