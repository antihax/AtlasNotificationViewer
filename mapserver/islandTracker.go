package mapserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type AtlasIsland struct {
	IslandPoints                             int      `json:"islandPoints"`
	IslandTreasureBottleSupplyCrateOverrides string   `json:"islandTreasureBottleSupplyCrateOverrides"`
	Discoveries                              []string `json:"discoveries"`
	UseLevelBoundsForTreasures               bool     `json:"useLevelBoundsForTreasures"`
	Rotation                                 float64  `json:"rotation"`
	UseNpcVolumesForTreasures                bool     `json:"useNpcVolumesForTreasures"`
	Sublevels                                []string `json:"sublevels"`
	WorldY                                   int      `json:"worldY"`
	MaxTreasureQuality                       int      `json:"maxTreasureQuality"`
	IslandHeight                             int      `json:"islandHeight"`
	Grid                                     string   `json:"grid"`
	SpawnerOverrides                         struct {
	} `json:"spawnerOverrides"`
	WorldX                        int      `json:"worldX"`
	PrioritizeVolumesForTreasures bool     `json:"prioritizeVolumesForTreasures"`
	Overrides                     []string `json:"overrides"`
	MinTreasureQuality            int      `json:"minTreasureQuality"`
	ID                            string   `json:"id"`
	IslandWidth                   int      `json:"islandWidth"`
	Resources                     []string `json:"resources"`
	Name                          string   `json:"name"`
}

type Island struct {
	IslandID             int       `json:"IslandID"`
	X                    float64   `json:"X"`
	Y                    float64   `json:"Y"`
	Size                 float64   `json:"Size"`
	TribeID              int       `json:"TribeId"`
	Color                string    `json:"Color"`
	IslandPoints         int       `json:"IslandPoints"`
	SettlementName       string    `json:"SettlementName"`
	TaxRate              float64   `json:"TaxRate"`
	CombatPhaseStartTime int       `json:"CombatPhaseStartTime"`
	WarringTribeID       int       `json:"WarringTribeID"`
	WarStartUTC          int       `json:"WarStartUTC"`
	WarEndUTC            int       `json:"WarEndUTC"`
	NumSettlers          int       `json:"NumSettlers"`
	LastUpdate           time.Time `json:"LastUpdate"`
}

type Company struct {
	TribeID   int         `json:"TribeId"`
	TribeName string      `json:"TribeName"`
	FlagURL   interface{} `json:"FlagURL"`
}

type IslandPackage struct {
	Version   int       `json:"version"`
	Islands   []Island  `json:"Islands"`
	Companies []Company `json:"Companies"`
}

func (s *MapServer) trackIslandData() {
	ticker := time.Tick(time.Minute * 5)

	url := getEnv("ISLANDS_URL", "")
	if url == "" {
		log.Println("No island url set; not tracking.")
		return
	}
	for {
		// Update the island package
		err := getJSON(url, &s.islandPackage)
		if err != nil {
			log.Println(err)
			<-ticker
			continue
		}

		for _, corp := range s.islandPackage.Companies {
			s.companies[corp.TribeID] = corp
		}

		// Look for ownership changes
		for _, island := range s.islandPackage.Islands {
			i, ok := s.islands[island.IslandID]
			if ok {
				if i.TribeID != island.TribeID {
					go s.passthroughWebhook(&notification{
						Event: "territorychange",
						Content: fmt.Sprintf("%s claimed %s in %s",
							s.companies[i.TribeID].TribeName,
							i.SettlementName,
							s.atlasIslands[i.IslandID].Grid,
						),
					})
				}
			}

			// Update to latest version
			island.LastUpdate = time.Now()
			s.islands[island.IslandID] = island
		}

		for k, island := range s.islands {
			if island.LastUpdate.Before(time.Now().Add(-time.Minute)) {
				go s.passthroughWebhook(&notification{
					Event: "territorychange",
					Content: fmt.Sprintf("%s lost %s in %s",
						s.companies[island.TribeID].TribeName,
						island.SettlementName,
						s.atlasIslands[island.IslandID].Grid,
					),
				})
				delete(s.islands, k)
			}
		}

		<-ticker
	}
}

func (s *MapServer) getIslands(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(s.islandPackage); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getJSON(url string, v interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return json.Unmarshal(b, v)
}

func (s *MapServer) loadIslandData() error {
	data, err := ioutil.ReadFile(getEnv("STATIC_DIR", "./www") + "/json/islands.json")
	if err != nil {
		return err
	}

	// Unpack island data
	var unpacked map[string]*json.RawMessage
	if err := json.Unmarshal(data, &unpacked); err != nil {
		return err
	}

	for key, islandData := range unpacked {
		islandID, err := strconv.Atoi(key)
		if err != nil {
			return err
		}
		island := &AtlasIsland{}
		json.Unmarshal(*islandData, island)
		s.atlasIslands[islandID] = island
	}
	return nil
}
