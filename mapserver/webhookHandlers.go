package mapserver

import (
	"bytes"
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
	Event     string `json:"event,omitempty"`
	Long      string `json:"long,omitempty"`
	Lat       string `json:"lat,omitempty"`
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

		n := notification{
			Tribe:     tribe,
			TribeID:   tribeID,
			Time:      time.Now().Unix(),
			AtlasTime: eventMatch[1],
			Content:   eventMatch[2],
			Event:     categorizeEvent(eventMatch[2]),
		}

		// add event to the queue
		s.eventQueue <- n
	}
}

var eventTypes = []string{
	"stealing",
	"settler",
	"demolish",
	"claim",
	"tame",
	"playerkill",
	"npc",
	"info",
}

func categorizeEvent(e string) string {
	switch {
	case regexp.MustCompile(`Company of .* is stealing your `).MatchString(e):
		return "stealing"
	case regexp.MustCompile(`has become a settler in your Settlement`).MatchString(e):
		return "settler"
	case regexp.MustCompile(`demolished a .* at [A-Z]{1}[0-9]{1} `).MatchString(e):
		return "demolish"
	case regexp.MustCompile(`Successfully unclaimed `).MatchString(e):
		return "claim"
	case regexp.MustCompile(`Claiming territory at `).MatchString(e):
		return "claim"
	case regexp.MustCompile(`Successfully claimed territory at `).MatchString(e):
		return "claim"
	case regexp.MustCompile(`Successfully claimed territory at `).MatchString(e):
		return "tame"
	case regexp.MustCompile(`Crew member .+? was killed by .+ \\(.+\\)!`).MatchString(e):
		return "playerkill"
	case regexp.MustCompile(`\\(Crewmember\\) mutinied`).MatchString(e):
		return "npc"
	}
	return "info"
}

func (s *MapServer) eventPump() {
	for {
		n := <-s.eventQueue

		go s.passthroughWebhook(&n)
		// Encode to JSON
		json, err := json.Marshal(n)
		if err != nil {
			log.Println(err)
			continue
		}

		// Send to valid clients
		s.eventHub.broadcast <- append(json, byte('\n'))
	}
}

type discordWebhooks struct {
	Content   string `json:"content,omitempty"`
	Username  string `json:"username,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	TTS       bool   `json:"tts,omitempty"`
	File      string `json:"file,omitempty"`
}

func (s *MapServer) passthroughWebhook(n *notification) {
	s.passthroughLock.RLock()
	defer s.passthroughLock.RUnlock()
	for _, pass := range s.passthrough {
		if pass[n.Event] == "1" {
			b := new(bytes.Buffer)
			json.NewEncoder(b).Encode(discordWebhooks{
				Content:  n.Content,
				Username: n.Tribe,
			})
			http.Post(pass["url"], "application/json", b)
		}
	}
}
