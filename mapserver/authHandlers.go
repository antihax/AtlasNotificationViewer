package mapserver

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/prometheus/common/log"
	"github.com/solovev/steam_go"
)

func (s *MapServer) logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := s.store.Get(r, "atlas-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}

	// Force deletion of the cookie and session
	session.Options.MaxAge = -1

	// Save changes
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}

	// Go back home
	http.Redirect(w, r, "/", 301)
}

func (s *MapServer) loginHandler(w http.ResponseWriter, r *http.Request) {
	opID := steam_go.NewOpenId(r)
	switch opID.Mode() {
	case "":
		http.Redirect(w, r, opID.AuthUrl(), 301)
	case "cancel":
		http.Redirect(w, r, "/", 301)
	default:
		// Validate we got an authenticated response and get the steamID
		steamID, err := opID.ValidateAndGetId()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		// Get our session
		session, err := s.store.Get(r, "atlas-session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}

		// Set the steamID in the session
		session.Values["steamID"] = steamID

		// Save the session values to db
		if err := session.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}

		// Set last logon
		_, err = s.redisClient.HSet(context.Background(), "user:"+steamID, "lastLogin", time.Now().Format(time.RFC822)).Result()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}

		// Go back home
		http.Redirect(w, r, "/", 301)
	}
}

// Determine if the request is an administrator
func (s *MapServer) isAdmin(r *http.Request) bool {
	session, err := s.store.Get(r, "atlas-session")
	if err != nil {
		return false
	}
	steamID, ok := session.Values["steamID"].(string)
	if ok {
		if steamID == s.globalAdmin {
			return true
		}
		if ok {
			v, err := s.redisClient.HGet(context.Background(), "user:"+steamID, "admin").Result()
			if err == nil && v == "1" {
				return true
			}
		}
	}
	return false
}

// Determine if the request is an Authorized to access the application
func (s *MapServer) isAuthorized(r *http.Request) bool {
	session, err := s.store.Get(r, "atlas-session")
	if err != nil {
		return false
	}
	steamID, ok := session.Values["steamID"].(string)
	if steamID == s.globalAdmin {
		return true
	}
	if ok {
		v, err := s.redisClient.HGet(context.Background(), "user:"+steamID, "allowed").Result()
		if err == nil && v == "1" {
			return true
		}
	}

	return false
}

func (s *MapServer) accountHandler(w http.ResponseWriter, r *http.Request) {
	session, err := s.store.Get(r, "atlas-session")
	steamID, ok := session.Values["steamID"].(string)
	if !ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}

	v, err := s.redisClient.HGetAll(context.Background(), "user:"+steamID).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}

	if steamID == s.globalAdmin {
		v["admin"] = "1"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
