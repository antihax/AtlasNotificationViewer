package atlasadminserver

import (
	"fmt"
	"log"

	"github.com/antihax/AtlasMap/internal/atlasdb"
	"github.com/gorilla/sessions"

	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

// AtlasAdminServer provides administrative services to an Atlas Cluster over http
type AtlasAdminServer struct {
	redisClient  *redis.Client
	gameData     map[string]EntityInfo
	gameDataLock sync.RWMutex

	tribeData     map[string]map[string]string
	tribeDataLock sync.RWMutex

	playerData     map[string]map[string]string
	playerDataLock sync.RWMutex

	steamData     map[string]string
	steamDataLock sync.RWMutex

	config *Configuration

	db *atlasdb.AtlasDB

	store *sessions.FilesystemStore
}

// NewAtlasAdminServer creates a new server
func NewAtlasAdminServer() *AtlasAdminServer {
	return &AtlasAdminServer{
		gameData:   make(map[string]EntityInfo),
		tribeData:  make(map[string]map[string]string),
		playerData: make(map[string]map[string]string),
		steamData:  make(map[string]string),
	}
}

// EntityInfo record in redis.
type EntityInfo struct {
	EntityID                string
	ParentEntityID          string
	EntityType              string
	EntitySubType           string
	EntityName              string
	TribeID                 string
	ServerXRelativeLocation float64
	ServerYRelativeLocation float64
	ServerID                [2]uint16
	LastUpdatedDBAt         uint64
	NextAllowedUseTime      uint64
}

// ServerLocation relative percentage to a specific server's origin.
type ServerLocation struct {
	ServerID                [2]uint16
	ServerXRelativeLocation float64
	ServerYRelativeLocation float64
}

// Run starts the server processing
func (s *AtlasAdminServer) Run() error {

	// Load configuration from environment
	if err := s.loadConfig(); err != nil {
		return err
	}

	s.store = sessions.NewFilesystemStore(s.config.SessionStore, []byte(s.config.SessionKey))
	/*
		// Setup redis connection
		s.redisClient = redis.NewClient(&redis.Options{
			Addr:     s.config.RedisAddress,
			Password: s.config.RedisPassword,
			DB:       s.config.RedisDB,
		})

		// Test connection
		_, err := s.redisClient.Ping().Result()
		if err != nil {
			return err
		}*/

	// Setup our DB pool
	db, err := atlasdb.NewAtlasDB(
		s.config.AtlasRedisAddress,
		s.config.AtlasRedisPassword,
		s.config.AtlasRedisDB,
	)
	if err != nil {
		return err
	}
	s.db = db

	// Poll the database for data
	go s.fetch()

	// Setup handlers
	http.Handle("/gettribes", ourMiddleware(http.HandlerFunc(s.getTribes)))
	http.Handle("/getdata", ourMiddleware(http.HandlerFunc(s.getData)))
	http.Handle("/travels", ourMiddleware(http.HandlerFunc(s.getPathTravelled)))
	http.Handle("/territoryURL", ourMiddleware(http.HandlerFunc(s.getTerritoryURL)))
	http.Handle("/", http.FileServer(http.Dir(s.config.StaticDir)))

	http.Handle("/login", ourMiddleware(http.HandlerFunc(s.loginHandler)))
	http.Handle("/logout", ourMiddleware(http.HandlerFunc(s.logoutHandler)))
	http.Handle("/account", ourMiddleware(http.HandlerFunc(s.accountHandler)))

	// Don't serve command handler if disabled
	if !s.config.DisableCommands {
		http.HandleFunc("/command", s.sendCommand)
	}

	endpoint := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	log.Println("Listening on ", endpoint)
	return http.ListenAndServe(endpoint, nil)
}

func (s *AtlasAdminServer) fetch() {
	kidsWithBadParents := make(map[string]bool)
	throttle := time.NewTicker(time.Duration(s.config.FetchRateInSeconds) * time.Second)

	for {
		entities := make(map[string]EntityInfo)

		// Get all tribe information
		tribes, err := s.db.GetAllTribes()
		if err != nil {
			log.Println(err)
			continue
		}

		s.tribeDataLock.Lock()
		s.tribeData = tribes
		s.tribeDataLock.Unlock()

		// Get the steamID=>playerID lookup
		steamData, err := s.db.GetReverseSteamIDMap()
		if err != nil {
			log.Println(err)
			continue
		}
		s.steamDataLock.Lock()
		s.steamData = steamData
		s.steamDataLock.Unlock()

		// Get player data
		playerData, err := s.db.GetAllPlayers()
		if err != nil {
			log.Println(err)
			continue
		}
		s.playerDataLock.Lock()
		s.playerData = playerData
		s.playerDataLock.Unlock()

		// Get all beds and ships
		e, err := s.db.GetAllEntities()
		if err != nil {
			log.Println(err)
			continue
		}
		for _, record := range e {
			info := newEntityInfo(record)
			entities[info.EntityID] = *info
		}

		// sanity check entity data, e.g. any missing parent ids?
		for k, v := range entities {
			if v.ParentEntityID != "0" {
				if _, parentFound := entities[v.ParentEntityID]; !parentFound {
					if _, dontSpamLog := kidsWithBadParents[k]; !dontSpamLog {
						kidsWithBadParents[k] = true
					}
					delete(entities, k)
				}
			}
		}
		s.gameDataLock.Lock()
		s.gameData = entities
		s.gameDataLock.Unlock()

		<-throttle.C
	}
}

// newEntityInfo transforms a Key-Value record into a new EntityInfo object.
func newEntityInfo(record map[string]string) *EntityInfo {
	var info EntityInfo
	info.EntityID = record["EntityID"]
	info.ParentEntityID = record["ParentEntityID"]
	info.EntityName = record["EntityName"]
	info.EntityType = record["EntityType"]
	info.ServerXRelativeLocation, _ = strconv.ParseFloat(record["ServerXRelativeLocation"], 64)
	info.ServerYRelativeLocation, _ = strconv.ParseFloat(record["ServerYRelativeLocation"], 64)
	info.LastUpdatedDBAt, _ = strconv.ParseUint(record["LastUpdatedDBAt"], 10, 64)
	info.NextAllowedUseTime, _ = strconv.ParseUint(record["NextAllowedUseTime"], 10, 64)

	var ok bool
	var tmpTribeID string
	if tmpTribeID, ok = record["TribeID"]; !ok {
		tmpTribeID = record["TribeId"]
	}
	info.TribeID = tmpTribeID

	var tmpServerID string
	if tmpServerID, ok = record["ServerID"]; !ok {
		tmpServerID = record["ServerId"]
	}
	info.ServerID, _ = atlasdb.ServerID(tmpServerID)

	// convert entity class to a subtype
	var tmpEntityClass string
	if tmpEntityClass, ok = record["EntityClass"]; !ok {
		tmpEntityClass = "none"
	}
	tmpEntityClass = strings.ToLower(tmpEntityClass)
	if strings.Contains(tmpEntityClass, "none") {
		info.EntitySubType = "None"
	} else if strings.Contains(tmpEntityClass, "brigantine") {
		info.EntitySubType = "Brigantine"
	} else if strings.Contains(tmpEntityClass, "dinghy") {
		info.EntitySubType = "Dingy"
	} else if strings.Contains(tmpEntityClass, "raft") {
		info.EntitySubType = "Raft"
	} else if strings.Contains(tmpEntityClass, "sloop") {
		info.EntitySubType = "Sloop"
	} else if strings.Contains(tmpEntityClass, "schooner") {
		info.EntitySubType = "Schooner"
	} else if strings.Contains(tmpEntityClass, "galleon") {
		info.EntitySubType = "Galleon"
	} else {
		info.EntitySubType = "None"
	}

	return &info
}
