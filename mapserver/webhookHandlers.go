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
	reCoords := regexp.MustCompile("Long: ([0-9.]+) / Lat: ([0-9.]+)")

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
		event = strings.Replace(event, "itude:", ":", -1)

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

		// Split coords
		coordsMatch := reCoords.FindStringSubmatch(event)
		if len(coordsMatch) == 3 {
			n.Long = coordsMatch[1]
			n.Lat = coordsMatch[2]
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
	"territorychange",
}

func categorizeEvent(e string) string {
	switch {

	case regexp.MustCompile(`has become a settler in your Settlement`).MatchString(e):
		return "settler"
	case regexp.MustCompile(`demolished a .* at [A-Z]{1}[0-9]{1} `).MatchString(e):
		return "demolish"
	case regexp.MustCompile(`Successfully unclaimed `).MatchString(e):
		return "claim"
	case regexp.MustCompile(`Successfully claimed territory at `).MatchString(e):
		return "claim"
	case regexp.MustCompile(`claim .+? lost`).MatchString(e):
		return "claim"
	case regexp.MustCompile(`failed to steal territory`).MatchString(e):
		return "claim"
	case regexp.MustCompile(` is stealing your territory at`).MatchString(e):
		return "claim"
	case regexp.MustCompile(`Crew member .+? was killed by .+ \\(.+\\)!`).MatchString(e):
		return "playerkill"
	case regexp.MustCompile(`\\(Crewmember\\) mutinied`).MatchString(e):
		return "npc"
	case regexp.MustCompile(` is stealing your `).MatchString(e):
		return "stealing"
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
	Content   string          `json:"content,omitempty"`
	Username  string          `json:"username,omitempty"`
	AvatarURL string          `json:"avatar_url,omitempty"`
	TTS       bool            `json:"tts,omitempty"`
	File      string          `json:"file,omitempty"`
	Embeds    []*MessageEmbed `json:"embeds,omitempty"`
}

// MessageEmbedFooter is a part of a MessageEmbed struct.
type MessageEmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// MessageEmbedImage is a part of a MessageEmbed struct.
type MessageEmbedImage struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

// MessageEmbedThumbnail is a part of a MessageEmbed struct.
type MessageEmbedThumbnail struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

// MessageEmbedVideo is a part of a MessageEmbed struct.
type MessageEmbedVideo struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

// MessageEmbedProvider is a part of a MessageEmbed struct.
type MessageEmbedProvider struct {
	URL  string `json:"url,omitempty"`
	Name string `json:"name,omitempty"`
}

// MessageEmbedAuthor is a part of a MessageEmbed struct.
type MessageEmbedAuthor struct {
	URL          string `json:"url,omitempty"`
	Name         string `json:"name,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// MessageEmbedField is a part of a MessageEmbed struct.
type MessageEmbedField struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

// An MessageEmbed stores data for message embeds.
type MessageEmbed struct {
	URL         string                 `json:"url,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
	Timestamp   string                 `json:"timestamp,omitempty"`
	Color       int                    `json:"color,omitempty"`
	Footer      *MessageEmbedFooter    `json:"footer,omitempty"`
	Image       *MessageEmbedImage     `json:"image,omitempty"`
	Thumbnail   *MessageEmbedThumbnail `json:"thumbnail,omitempty"`
	Video       *MessageEmbedVideo     `json:"video,omitempty"`
	Provider    *MessageEmbedProvider  `json:"provider,omitempty"`
	Author      *MessageEmbedAuthor    `json:"author,omitempty"`
	Fields      []*MessageEmbedField   `json:"fields,omitempty"`
}

func (s *MapServer) passthroughWebhook(n *notification) {
	s.passthroughLock.RLock()
	defer s.passthroughLock.RUnlock()
	for _, pass := range s.passthrough {
		if pass[n.Event] == "1" {
			url := s.ourURL

			h := discordWebhooks{
				Embeds:   []*MessageEmbed{},
				Username: n.Tribe,
			}
			h.Embeds = append(h.Embeds, &MessageEmbed{
				Title: n.Content,
				Author: &MessageEmbedAuthor{
					Name: n.Tribe,
				},
			})
			if n.Lat != "" {
				url += "?lat=" + n.Lat + "&long=" + n.Long
				h.Embeds[0].URL = url
			}
			b := new(bytes.Buffer)

			json.NewEncoder(b).Encode(h)
			http.Post(pass["url"], "application/json", b)
		}
	}
}
