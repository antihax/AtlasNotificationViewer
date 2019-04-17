package mapserver

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"net/http"

	"github.com/go-redis/redis"
	gsr "gopkg.in/boj/redistore.v1"
)

// MapServer provides administrative services to an Atlas Cluster over http
type MapServer struct {
	redisClient  *redis.Client
	territoryURL string
	globalAdmin  string
	store        *gsr.RediStore

	eventHub *Hub
}

// NewMapServer creates a new server
func NewMapServer() *MapServer {
	return &MapServer{
		eventHub: newHub(),
	}
}

func (s *MapServer) addHandler(uri string, handlerFunc http.HandlerFunc) {
	http.Handle(uri, s.logMiddleware(handlerFunc))
}

func (s *MapServer) addAuthHandler(uri string, handlerFunc http.HandlerFunc) {
	http.Handle(uri, s.authenticatedEndpoint(s.logMiddleware(handlerFunc)))
}

func (s *MapServer) addAdminHandler(uri string, handlerFunc http.HandlerFunc) {
	http.Handle(uri, s.adminEndpoint(s.logMiddleware(handlerFunc)))
}

// Run starts the server processing
func (s *MapServer) Run() error {
	testNumeric := regexp.MustCompile("^[0-9]+$")
	testAlphaNumeric := regexp.MustCompile("^[0-9a-zA-Z._-]+$")

	// Setup redis connection
	s.redisClient = redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_ADDRESS", "localhost:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
	})

	s.globalAdmin = getEnv("GLOBAL_ADMIN_STEAMID", "")
	if !testNumeric.MatchString(s.globalAdmin) {
		log.Fatal("Must have a valid SteamID set in GLOBAL_ADMIN_STEAMID environment variable")
	}

	// Create a redis session store.
	store, err := gsr.NewRediStore(10, "tcp", getEnv("REDIS_ADDRESS", "localhost:6379"), getEnv("REDIS_PASSWORD", ""), []byte(os.Getenv("COOKIE_SECRET")))
	if err != nil {
		log.Fatalf("Cannot build redis store: %v", err)
	}
	// Set options for the store
	store.SetMaxLength(1024 * 100)
	s.store = store

	// Test connection
	_, err = s.redisClient.Ping().Result()
	if err != nil {
		return err
	}

	// Setup handlers
	http.Handle("/", http.FileServer(http.Dir(getEnv("STATIC_DIR", "./www"))))

	s.addHandler("/login", s.loginHandler)
	s.addHandler("/logout", s.logoutHandler)

	s.addAuthHandler("/account", s.accountHandler)
	s.addAuthHandler("/events", s.streamEvents)

	s.addAdminHandler("/listUsers", s.listUsers)
	s.addAdminHandler("/changeUser", s.changeUser)

	webhook := getEnv("WEBHOOK_KEY", "")
	if !testAlphaNumeric.MatchString(webhook) {
		log.Fatal("Webhook key is not valid. Must be alpha numeric")
	}
	s.addHandler("/"+webhook, http.HandlerFunc(s.webhookHandler))

	go s.eventHub.run()
	endpoint := fmt.Sprintf("%s:%s", getEnv("HOST", "localhost"), getEnv("PORT", "8000"))
	log.Println("Listening on ", endpoint)
	return http.ListenAndServe(endpoint, nil)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
