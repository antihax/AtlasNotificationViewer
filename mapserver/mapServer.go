package mapserver

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"

	"net/http"

	"github.com/go-redis/redis"
	gsr "gopkg.in/boj/redistore.v1"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// MapServer provides administrative services to an Atlas Cluster over http
type MapServer struct {
	redisClient  *redis.Client
	territoryURL string
	homeURL      string
	globalAdmin  string
	store        *gsr.RediStore

	passthrough     map[string]map[string]string
	passthroughLock sync.RWMutex

	eventQueue chan notification
	eventHub   *Hub
}

// NewMapServer creates a new server
func NewMapServer() *MapServer {
	return &MapServer{
		eventHub:    newHub(),
		eventQueue:  make(chan notification, 1000),
		passthrough: make(map[string]map[string]string),
	}
}

func (s *MapServer) addHandler(uri string, handlerFunc http.HandlerFunc) {
	http.Handle(uri, s.logMiddleware(handlerFunc))
}

func (s *MapServer) addAuthHandler(uri string, handlerFunc http.HandlerFunc) {
	http.Handle(uri, s.noCache(s.authenticatedEndpoint(s.logMiddleware(handlerFunc))))
}

func (s *MapServer) addNoCacheHandler(uri string, handlerFunc http.HandlerFunc) {
	http.Handle(uri, s.noCache(s.logMiddleware(handlerFunc)))
}

func (s *MapServer) addAdminHandler(uri string, handlerFunc http.HandlerFunc) {
	http.Handle(uri, s.noCache(s.adminEndpoint(s.logMiddleware(handlerFunc))))
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

	err := s.loadPassthrough()
	if err != nil {
		log.Fatalf("Cannot load passthroughs: %v", err)
	}

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

	http.Handle("/admin", s.adminEndpoint(
		http.FileServer(
			http.Dir(getEnv("STATIC_DIR", "./www/admin"))),
	),
	)

	s.addNoCacheHandler("/login", s.loginHandler)
	s.addNoCacheHandler("/logout", s.logoutHandler)

	s.addAuthHandler("/account", s.accountHandler)
	s.addAuthHandler("/events", s.streamEvents)

	s.addAdminHandler("/listUsers", s.listUsers)
	s.addAdminHandler("/changeUser", s.changeUser)

	s.addAdminHandler("/listPassthroughs", s.listPassthrough)
	s.addAdminHandler("/changePassthrough", s.changePassthrough)
	s.addAdminHandler("/addPassthrough", s.addPassthrough)
	s.addAdminHandler("/listEventTypes", s.listEventTypes)

	webhook := getEnv("WEBHOOK_KEY", "")
	if !testAlphaNumeric.MatchString(webhook) {
		log.Fatal("Webhook key is not valid. Must be alpha numeric")
	}
	s.addNoCacheHandler("/"+webhook, http.HandlerFunc(s.webhookHandler))

	// Run in the background
	go s.eventHub.run()
	go s.eventPump()

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
