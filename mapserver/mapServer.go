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
	return &MapServer{eventHub: newHub()}
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

	////////////////////////////
	//	go s.fakePump()
	////////////////////////////

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

/*
func (s *MapServer) fakePump() {
	timer := time.Tick(time.Millisecond * 250)
	var stuff = `{"tribe": "SomeDudes", "tribeID": 1234567893, "time": 123123123, "atlasTime": "Day 36, 15:44:35", "content": "Someone was promoted to a Company Admin by Someone!"}
{"tribe": "SomeDudes", "tribeID": 1234567893, "time": 123123123, "atlasTime": "Day 36, 20:10:02", "content": "Crew member Someone - Lvl 31 was killed!"}
{"tribe": "SomeDudes", "tribeID": 1234567893, "time": 123123123, "atlasTime": "Day 36, 20:13:01", "content": "Someone demolished a 'Wood Roof' at N6 [Long: 81.10 / Lat: 28.25]!"}
{"tribe": "SomeDudes", "tribeID": 1234567893, "time": 123123123, "atlasTime": "Day 36, 16:51:20", "content": "Crew member Someone - Lvl 1 was killed!"}
{"tribe": "SomeDudes", "tribeID": 1234567893, "time": 123123123, "atlasTime": "Day 36, 17:34:01", "content": "Crew member Someone - Lvl 27 was killed!"}
{"tribe": "SomeDudes", "tribeID": 1234567893, "time": 123123123, "atlasTime": "Day 36, 15:40:22", "content": "Crew member Someone - Lvl 30 was killed!"}
{"tribe": "SomeDudes", "tribeID": 1234567893, "time": 123123123, "atlasTime": "Day 36, 15:52:01", "content": "Crew member Someone - Lvl 33 was killed!"}
{"tribe": "SomeDudes", "tribeID": 1234567893, "time": 123123123, "atlasTime": "Day 36, 15:58:37", "content": "Crew member Someone - Lvl 19 was killed by a Rhino - Lvl 19!"}
{"tribe": "SomeDudes", "tribeID": 1234567893, "time": 123123123, "atlasTime": "Day 36, 16:00:13", "content": "Crew member Someone - Lvl 12 was killed by a Rhino - Lvl 19!"}
{"tribe": "SomeDudes", "tribeID": 1234567893, "time": 123123123, "atlasTime": "Day 36, 23:15:15", "content": "Bed 234234234234 was added to the Company!"}
{"tribe": "SomeDudes", "tribeID": 1234567893, "time": 123123123, "atlasTime": "Day 36, 23:16:48", "content": "Someone demolished a 'Bed' at N6 [Long: 81.51 / Lat: 28.11]!"}
{"tribe": "SomeDudes", "tribeID": 1234567893, "time": 123123123, "atlasTime": "Day 36, 23:16:55", "content": "Bed -234234234234 was removed from the Company!"}
{"tribe": "SomeDudes", "tribeID": 1234567893, "time": 123123123, "atlasTime": "Day 36, 23:20:29", "content": "Someone demolished a 'PillarTop (Bed)' at N6 [Long: 81.51 / Lat: 28.12]!"}
{"tribe": "SomeDudes", "tribeID": 1234567893, "time": 123123123, "atlasTime": "Day 36, 23:20:36", "content": "Bed 234234234234 was removed from the Company!"}
{"tribe": "SomeDudes", "tribeID": 1234567893, "time": 123123123, "atlasTime": "Day 36, 19:16:45", "content": "Someone demolished a 'Wood Doorway' at J5 [Long: 24.89 / Lat: 37.40]!"}
{"tribe": "SomeDudes", "tribeID": 1234567893, "time": 123123123, "atlasTime": "Day 36, 19:18:14", "content": "Someone demolished a 'Wood Wall' at J5 [Long: 24.89 / Lat: 37.40]!"}
{"tribe": "SomeDudes", "tribeID": 1234567893, "time": 123123123, "atlasTime": "Day 36, 19:21:38", "content": "Someone Tamed a Rhino - Lvl 27 (Rhino)!"}
{"tribe": "SomeDudes", "tribeID": 1234567893, "time": 123123123, "atlasTime": "Day 36, 20:19:46", "content": "Someone demolished a 'Wood Doorway' at J5 [Long: 24.89 / Lat: 37.42]!"}
{"tribe": "SomeDudes", "tribeID": 1234567893, "time": 123123123, "atlasTime": "Day 36, 20:21:51", "content": "Someone demolished a 'Wood Doorway' at J5 [Long: 24.89 / Lat: 37.41]!"}
{"tribe": "SomeDudes", "tribeID": 1234567893, "time": 123123123, "atlasTime": "Day 36, 20:28:16", "content": "Someone demolished a 'Wood Doorway' at J5 [Long: 24.89 / Lat: 37.41]!"}
{"tribe": "SomeDudes", "tribeID": 1234567893, "time": 123123123, "atlasTime": "Day 36, 20:30:07", "content": "Someone demolished a 'Wood Doorway' at J5 [Long: 24.89 / Lat: 37.42]!"}
{"tribe": "SomeDudes", "tribeID": 1234567893, "time": 123123123, "atlasTime": "Day 37, 01:39:50", "content": "Someone Recruited Bearded John the Wanderer - Lvl 1 (Crewmember)!"}
{"tribe": "SomeDudes", "tribeID": 1234567893, "time": 123123123, "atlasTime": "Day 37, 01:40:44", "content": "Someone Recruited Old Joe Silver - Lvl 1 (Crewmember)!"}
{"tribe": "SomeDudes", "tribeID": 1234567893, "time": 123123123, "atlasTime": "Day 37, 01:41:09", "content": "Someone Recruited Angry Alex the Wanderer - Lvl 1 (Crewmember)!"}
{"tribe": "SomeDudes", "tribeID": 1234567893, "time": 123123123, "atlasTime": "Day 37, 01:41:30", "content": "Someone Recruited Pretty Alex Senior - Lvl 1 (Crewmember)!"}
`
	messages := strings.Split(stuff, "\n")
	len := len(messages)
	for {
		s.eventHub.broadcast <- []byte(messages[rand.Intn(len)])
		<-timer
	}
}
*/
