package mapserver

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type notification struct {
	Tribe     string `json:"tribe"`
	TribeID   int64  `json:"tribeID"`
	Time      int64  `json:"time"`
	AtlasTime string `json:"atlasTine"`
	Content   string `json:"content"`
}

type atlasWebhook struct {
	Content string `json:"content"`
}

func (s *MapServer) webhookHandler(w http.ResponseWriter, r *http.Request) {
	reTribe := regexp.MustCompile("^\\*\\*(.+?) \\(([0-9]+)\\)\\*\\*")
	reEvent := regexp.MustCompile("^(Day .+?): (.+?)$")

	decoder := json.NewDecoder(r.Body)
	var t atlasWebhook
	err := decoder.Decode(&t)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Build list of events
	events := strings.Split(t.Content, "\n")

	// Get tribe data from first line
	tribeData := reTribe.FindStringSubmatch(events[0])
	if len(tribeData) != 3 {
		http.Error(w, "Invalid tribe string", http.StatusBadRequest)
		return
	}

	// Get Tribe Data
	tribe := tribeData[1]
	tribeID, err := strconv.ParseInt(tribeData[2], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Process all remaining events
	for _, event := range events[1:] {
		if event == "" {
			continue
		}

		// Split date from event
		eventMatch := reEvent.FindStringSubmatch(event)
		if len(eventMatch) != 3 {
			log.Printf("Did not match %#v\n", eventMatch)
			continue
		}

		// Encode to JSON
		json, err := json.Marshal(notification{
			Tribe:     tribe,
			TribeID:   tribeID,
			Time:      time.Now().Unix(),
			AtlasTime: eventMatch[1],
			Content:   eventMatch[2],
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			continue
		}

		// Send to valid clients
		s.eventHub.broadcast <- append(json, byte('\n'))
	}
}
