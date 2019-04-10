package mapserver

import (
	"log"
	"net/http"
)

func (s *MapServer) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func (s *MapServer) authenticatedEndpoint(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.isAuthorized(r) {
			http.Error(w, "Not authorized", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *MapServer) adminEndpoint(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.isAdmin(r) {
			http.Error(w, "Not an admin", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)

	})
}
