package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

type atlasWebhook struct {
	Content string `json:"content"`
}

func main() {
	client := &http.Client{}
	throttle := time.Tick(time.Second)

	len := len(mock)
	for {
		<-throttle

		// Encapsulate a random message
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(
			atlasWebhook{
				mock[rand.Intn(len)],
			},
		)

		req, err := http.NewRequest("POST", getEnv("URL", "http://localhost:8000/eventwebhook"), buf)
		if err != nil {
			log.Println(err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			continue
		}

		b, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(b))
		log.Printf("%+v\n", resp)
		resp.Body.Close()
	}
}

var mock = []string{`**Some Tribe (1234567892)**
Day 150, 09:54:36: Someone demolished a 'Ballista Turret' at N6 [Long: 81.42 / Lat: 28.16]!
Day 150, 04:53:29: Crew member Someone - Lvl 44 was killed!
Day 150, 11:44:14: Someone demolished a 'Ballista Turret' at N6 [Long: 81.42 / Lat: 28.16]!
Day 150, 11:56:10: Someone Tamed a Bear - Lvl 25 (Bear)!
`, `**Some Tribe (1234567892)**
Day 150, 14:41:08: Crew member Someone - Lvl 34 was killed!
Day 150, 15:31:09: Someone demolished a 'Stone Wall' at N6 [Long: 81.42 / Lat: 28.16]!
Day 150, 16:02:07: Someone demolished a 'MAIN BASE (Bed)' at N6 [Long: 81.45 / Lat: 28.16]!
Day 150, 16:02:15: Bed 15934942 was removed from the Company!
Day 150, 16:03:15: Someone demolished a 'MAIN BASE (Bed)' at N6 [Long: 81.45 / Lat: 28.16]!
Day 150, 16:03:22: Bed 1710793044 was removed from the Company!
Day 150, 16:39:46: Bed 1918720058 was added to the Company!
Day 150, 16:44:42: Bed 902891152 was added to the Company!
Day 150, 16:48:48: Bed -614056899 was added to the Company!
Day 150, 16:57:40: Bed -798065826 was added to the Company!
Day 150, 17:00:21: Bed -703981167 was added to the Company!
Day 150, 17:04:52: Bed 1678019731 was added to the Company!
`, `**Some Tribe (1234567892)**
Day 150, 17:18:25: Bed -1833734353 was added to the Company!
Day 150, 17:21:20: Someone demolished a 'Stone Square Ceiling' at N6 [Long: 81.45 / Lat: 28.15]!
Day 150, 17:21:27: Bed -1833734353 was removed from the Company!
Day 150, 17:21:27: Bed 1678019731 was removed from the Company!
Day 150, 17:21:27: Bed -703981167 was removed from the Company!
Day 150, 17:21:27: Bed -798065826 was removed from the Company!
Day 150, 17:21:27: Bed -614056899 was removed from the Company!
Day 150, 17:21:27: Bed 902891152 was removed from the Company!
Day 150, 17:21:27: Bed 1918720058 was removed from the Company!
Day 150, 17:27:50: Someone demolished a 'Thatch Roof' at N6 [Long: 81.45 / Lat: 28.16]!
Day 150, 17:29:32: Someone demolished a 'Stone Pillar' at N6 [Long: 81.45 / Lat: 28.15]!
Day 150, 17:56:31: Bed -1988431771 was added to the Company!
Day 150, 17:57:46: Bed -1483624002 was added to the Company!
Day 150, 17:58:50: Bed -1898544056 was added to the Company!
Day 150, 18:04:43: Bed Bed was renamed to Crafting Main!
`, `**Some Tribe (1234567892)**
Day 150, 20:24:08: Crew member Someone - Lvl 51 was killed!
Day 150, 16:42:56: Your claim of Catzen Shizer has been interrupted! (K12)
Day 150, 21:32:54: Crew member Someone - Lvl 51 was killed!
Day 150, 17:03:24: Failed to steal Catzen Shizer! (K12)
Day 150, 21:54:22: Someone demolished a 'Territory Taxation Bank (Pin Coded)' at N6 [Long: 81.45 / Lat: 28.16]!
Day 150, 22:39:17: Ship 1627561374 was added to the Company!
`, `**Some Tribe (1234567892)**
Day 150, 23:00:20: Someone demolished a 'Stone Wall' at N6 [Long: 81.47 / Lat: 28.16]!
Day 150, 18:02:13: Bed -317003207 was added to the Company!
`, `**Some Tribe (1234567892)**
Day 151, 02:08:11: Someone demolished a 'Stone Square Ceiling' at N6 [Long: 81.43 / Lat: 28.15]!
`, `**Some Tribe (1234567892)**
Day 151, 02:14:13: Your 'Bed' was destroyed!
Day 151, 02:14:17: Bed 2051943427 was removed from the Company!
Day 151, 07:31:49: Someone demolished a 'Stone Wall' at N6 [Long: 81.43 / Lat: 28.16]!
Day 151, 07:39:42: Someone demolished a 'Stone Wall' at N6 [Long: 81.43 / Lat: 28.16]!
Day 151, 03:48:30: Your 'Large Storage Box (Pin Coded)' was destroyed!
`, `**Some Tribe (1234567892)**
Day 151, 09:33:40: Someone Tamed a Bear - Lvl 7 (Bear)!
Day 151, 10:29:26: Crew member Someone - Lvl 51 was killed by a Bear - Lvl 7!
Day 151, 10:46:00: Crew member Someone - Lvl 51 was killed!
Day 151, 11:46:05: Your 'Wood Half Wall' was destroyed!
`, `**Some Tribe (1234567892)**
Day 151, 13:04:24: Someone Recruited Two-Legged Bob Nine Toes - Lvl 1 (Crewmember)!
Day 151, 13:05:25: Someone Recruited Calico Mary the Mighty - Lvl 1 (Crewmember)!
Day 151, 13:17:48: Someone Recruited One-Legged Joan the Mighty - Lvl 4 (Crewmember)!
Day 151, 13:19:28: Someone Recruited Calico Joe Black - Lvl 4 (Crewmember)!
Day 151, 13:59:23: Someone Recruited Two-Legged Mary Senior - Lvl 7 (Crewmember)!
Day 151, 14:00:20: Someone Recruited Two-Legged Charlotte the Mighty - Lvl 7 (Crewmember)!
Day 151, 08:45:42: Crew member Someone - Lvl 49 was killed by a Cobra - Lvl 3!
`, `**Some Tribe (1234567892)**
Day 151, 15:48:37: Your Yoggie - Lvl 28 (Bear) was killed by a Bear - Lvl 29 at N6 [Long: 81.48 / Lat: 28.00]!
Day 151, 18:00:33: Someone Tamed a Bear - Lvl 1 (Bear)!
Day 151, 18:06:59: Your 'Wood Signpost' was destroyed!
Day 151, 18:06:59: Your 'Wood Signpost' was destroyed!
Day 151, 18:10:59: Your 'Wood Signpost' was destroyed!
`, `**Some Tribe (1234567892)**
Day 151, 21:34:26: Someone demolished a 'Stone Wall' at N6 [Long: 81.62 / Lat: 28.46]!
Day 151, 21:50:43: Someone demolished a 'Stone Wall' at N6 [Long: 81.61 / Lat: 28.46]!
Day 151, 21:51:51: Someone demolished a 'Stone Wall' at N6 [Long: 81.61 / Lat: 28.47]!
Day 151, 21:53:16: Someone demolished a 'Stone Wall' at N6 [Long: 81.61 / Lat: 28.47]!
Day 151, 21:54:31: Someone demolished a 'Stone Wall' at N6 [Long: 81.61 / Lat: 28.47]!
Day 151, 22:00:52: Someone demolished a 'Stone Wall' at N6 [Long: 81.60 / Lat: 28.47]!
Day 151, 22:01:57: Someone demolished a 'Stone Wall' at N6 [Long: 81.60 / Lat: 28.48]!
Day 151, 22:03:15: Someone demolished a 'Stone Wall' at N6 [Long: 81.60 / Lat: 28.47]!
`, `**Some Tribe (1234567892)**
Day 152, 00:58:34: Someone demolished a 'Stone Wall' at N6 [Long: 81.59 / Lat: 28.48]!
Day 152, 01:06:57: Someone demolished a 'Stone Doorway' at N6 [Long: 81.58 / Lat: 28.50]!
Day 152, 01:07:19: Someone Tamed a Bear - Lvl 43 (Bear)!
Day 152, 01:12:43: Someone demolished a 'Stone Square Ceiling' at N6 [Long: 81.58 / Lat: 28.50]!
Day 152, 02:37:15: Someone Recruited Calico Anne Junior - Lvl 6 (Crewmember)!
Day 152, 02:39:16: Someone Recruited Pretty Joan Nine Fingers - Lvl 6 (Crewmember)!
Day 152, 02:55:38: Someone Recruited Dread Joe the Strong - Lvl 6 (Crewmember)!
Day 152, 02:56:13: Someone Recruited Ugly Alex Nine Toes - Lvl 6 (Crewmember)!
Day 152, 02:57:07: Someone Recruited Old Charlotte the Strong - Lvl 6 (Crewmember)!
Day 152, 03:05:47: Someone Recruited Angry Sue the Mighty - Lvl 4 (Crewmember)!
Day 152, 03:07:27: Someone Recruited Calico John Silver - Lvl 15 (Crewmember)!
Day 152, 03:20:59: Someone demolished a 'Bookshelf' at N6 [Long: 81.44 / Lat: 28.15]!
`, `**Some Tribe (1234567892)**
Day 152, 03:38:20: Someone demolished a 'Bookshelf' at N6 [Long: 81.45 / Lat: 28.15]!
`, `**Some Tribe (1234567892)**
Day 152, 07:12:56: Crew member Someone - Lvl 39 was killed!
Day 152, 07:45:56: Crew member Someone - Lvl 51 was killed by a Giraffe - Lvl 16!
Day 152, 08:02:56: Crew member Someone - Lvl 49 was killed by a Pig - Lvl 1!
Day 152, 08:28:58: Someone demolished a 'Bookshelf' at N6 [Long: 81.44 / Lat: 28.15]!
`, `**Some Tribe (1234567892)**
Day 151, 14:40:08: Crew member Someone - Lvl 51 was killed!
Day 152, 11:44:42: Someone demolished a 'Wood Signpost' at N6 [Long: 81.41 / Lat: 28.01]!
Day 152, 11:46:25: Someone demolished a 'Wood Signpost' at N6 [Long: 81.42 / Lat: 28.02]!
Day 152, 11:48:16: Someone demolished a 'Wood Signpost' at N6 [Long: 81.41 / Lat: 28.02]!
Day 152, 11:50:00: Someone demolished a 'Wood Signpost' at N6 [Long: 81.41 / Lat: 28.02]!
Day 152, 11:52:40: Someone demolished a 'Wood Signpost' at N6 [Long: 81.41 / Lat: 28.01]!
Day 152, 11:55:24: Someone demolished a 'Wood Signpost' at N6 [Long: 81.41 / Lat: 28.02]!
Day 152, 12:04:08: Crew member Someone - Lvl 51 was killed by a Pig - Lvl 7!
Day 152, 12:39:54: Crew member Someone - Lvl 35 was killed!
`, `**Some Tribe (1234567892)**
Day 152, 14:31:23: Crew member Someone - Lvl 37 was killed!
Day 152, 15:16:45: 'Some Tribe' has left Alliance 'N6 Peeps'!
Day 152, 13:20:54: Crew member Someone - Lvl 51 was killed by an Elephant - Lvl 2!
Day 152, 15:35:08: Crew member Someone - Lvl 51 was killed by a Pig - Lvl 4!
`, `**Some Tribe (1234567892)**
Day 151, 20:32:09: Someone demolished a 'Campfire' at A6 - Northwest Tropical Freeport [Long: -89.64 / Lat: 27.50]!
Day 152, 14:40:38: Crew member Someone - Lvl 51 was killed by a Pig - Lvl 23!
Day 152, 18:48:42: Crew member Someone - Lvl 47 was killed by a Giraffe - Lvl 10!
Day 152, 13:42:49: Crew member Someone - Lvl 51 was killed!
`, `**Some Tribe (1234567892)**
Day 152, 00:02:14: Ship -2103950642 was added to the Company!
Day 152, 00:04:40: Ship Ramshackle Sloop was renamed to maps!
Day 152, 00:11:12: Crew member Someone - Lvl 49 was killed by a Bull - Lvl 2!
Day 152, 01:40:43: Bed -491537755 was added to the Company!
Day 152, 01:40:49: Crew member Someone - Lvl 49 was killed by a Shark - Lvl 8!
`, `**Some Tribe (1234567892)**
Day 153, 00:58:41: Crew member Someone - Lvl 51 was killed!
Day 153, 01:20:05: Someone demolished a 'Wood Signpost' at N6 [Long: 82.29 / Lat: 28.53]!
Day 153, 01:21:27: Someone demolished a 'Wood Signpost' at N6 [Long: 82.29 / Lat: 28.52]!
Day 153, 01:23:32: Someone demolished a 'Wood Signpost' at N6 [Long: 82.29 / Lat: 28.52]!
Day 152, 04:32:06: Ship 2051070174 was added to the Company!
Day 153, 01:25:51: Someone demolished a 'Wood Signpost' at N6 [Long: 82.29 / Lat: 28.53]!
Day 152, 04:35:50: Ship Raft was renamed to rafty mc raftface!
Day 152, 06:23:37: Crew member Someone - Lvl 49 was killed by a Bull - Lvl 15!
`, `**Some Tribe (1234567892)**
Day 153, 04:25:02: Crew member Someone - Lvl 51 was killed by a Seagull - Lvl 3!
Day 153, 04:48:18: Someone demolished a 'Wood Square Ceiling' at N6 [Long: 81.46 / Lat: 28.16]!
Day 153, 04:50:25: Someone demolished a 'Wood Square Ceiling' at N6 [Long: 81.46 / Lat: 28.16]!
`, `**Some Tribe (1234567892)**
Day 153, 07:44:58: Someone demolished a 'Ammo Container (Locked) ' at N6 [Long: 81.46 / Lat: 28.16]!
Day 153, 02:22:19: Crew member Someone - Lvl 51 was killed by an Elephant - Lvl 22!
Day 153, 09:46:58: Crew member Someone - Lvl 35 was killed!
`, `**Some Tribe (1234567892)**
Day 153, 10:46:09: Crew member LSG . - Lvl 39 was killed by a Rattlesnake - Lvl 1!
Day 153, 10:47:42: Crew member Someone - Lvl 38 was killed!
Day 153, 11:46:10: Crew member Someone - Lvl 38 was killed!
`, `**Some Tribe (1234567892)**
Day 152, 17:32:36: Someone demolished a 'Speed Ship Sail (Small)' at A6 - Northwest Tropical Freeport [Long: -89.48 / Lat: 27.63]!
Day 152, 17:37:35: Someone demolished a 'Speed Ship Sail (Small)' at A6 - Northwest Tropical Freeport [Long: -89.48 / Lat: 27.64]!
Day 153, 16:16:31: Someone Tamed a Bear - Lvl 20 (Bear)!
`, `**Some Tribe (1234567892)**
Day 153, 17:46:25: Incudem Malleus (Terminus - CO) - Lvl 13 has become a settler in your Settlement 'Home Place' (N6) [ Longitude: 80.97 / Latitude: 27.83 ]!
Day 153, 16:47:54: Crew member Someone - Lvl 51 was killed!
`, `**Some Tribe (1234567892)**
Day 153, 00:17:37: Bed -855120503 was added to the Company!
Day 153, 22:26:21: Someone demolished a 'Thatch Wall' at N6 [Long: 81.46 / Lat: 28.16]!
Day 153, 22:27:26: Someone demolished a 'Wood Half Wall' at N6 [Long: 81.46 / Lat: 28.16]!
Day 153, 20:33:15: Crew member Jata Bighead - Lvl 50 was killed by an Elephant - Lvl 9!
`, `**Some Tribe (1234567892)**
Day 153, 04:32:07: Crew member Someone - Lvl 51 was killed!
`, `**Some Tribe (1234567892)**
Day 154, 04:31:24: Crew member Someone - Lvl 51 was killed!
Day 154, 06:26:12: Someone Tamed a Bear - Lvl 14 (Bear)!
`, `**Some Tribe (1234567892)**
Day 153, 14:42:05: Crew member Someone - Lvl 49 was killed by a Shark - Lvl 3!
Day 153, 13:02:56: Crew member Someone - Lvl 47 was killed!
`, `**Some Tribe (1234567892)**
Day 154, 13:49:43: Crew member Someone - Lvl 38 was killed by a Shark - Lvl 5!
Day 154, 15:25:40: Someone demolished a 'Stone Square Ceiling' at N6 [Long: 81.62 / Lat: 28.46]!
Day 153, 20:40:25: Bed -1151434851 was added to the Company!
Day 154, 15:34:13: Someone demolished a 'Stone Wall' at N6 [Long: 81.62 / Lat: 28.46]!
`, `**Some Tribe (1234567892)**
Day 153, 21:59:02: Bed -63323165 was added to the Company!
Day 153, 22:08:33: Someone - Lvl 49 destroyed their 'Bed ()')!
Day 153, 23:12:08: Crew member Someone - Lvl 38 was killed by a Wolf - Lvl 3!
Day 154, 18:15:25: Crew member Someone - Lvl 37 was killed!
Day 154, 00:07:52: Bed 209614085 was added to the Company!
`, `**Some Tribe (1234567892)**
Day 154, 00:32:24: Bed 116370868 was added to the Company!
`, `**Some Tribe (1234567892)**
Day 154, 03:49:05: SavageNuke LaFever was removed from the Company by Someone!
Day 154, 04:31:27: Crew member Someone - Lvl 51 was killed by a Wolf - Lvl 29!
Day 154, 04:40:44: Crew member Someone - Lvl 47 was killed!
Day 154, 06:14:04: Bed -1565792037 was added to the Company!
`, `**Some Tribe (1234567892)**
Day 154, 06:38:25: Crew member Someone - Lvl 39 was killed by a Lion - Lvl 9!
Day 154, 07:43:46: Crew member Someone - Lvl 38 was killed by an Alpha Pig - Lvl 216!
`, `**Some Tribe (1234567892)**
Day 155, 04:03:53: Crew member Someone - Lvl 35 was killed!
Day 154, 09:44:02: Bed -1631962815 was added to the Company!
Day 154, 10:06:33: Crew member Someone - Lvl 51 was killed by a Bull - Lvl 4!
`, `**Some Tribe (1234567892)**
Day 154, 11:47:04: Your 'Bed' was destroyed!
Day 154, 11:47:07: Bed -1631962815 was removed from the Company!
Day 154, 12:00:46: Crew member Someone - Lvl 39 was killed by a Wolf - Lvl 20!
Day 155, 02:17:07: Crew member Someone - Lvl 35 was killed by a Pig - Lvl 10!
Day 154, 12:52:19: Crew member Someone - Lvl 38 was killed by Bearella - Lvl 45 (Bear) (DVS Gaming)!
`, `**Some Tribe (1234567892)**
Day 154, 14:39:56: Crew member Someone - Lvl 38 was killed by a Manta ray - Lvl 4!
Day 155, 11:31:48: Someone demolished a 'Stone Square Ceiling' at N6 [Long: 81.43 / Lat: 28.16]!
`, `**Some Tribe (1234567892)**
Day 154, 18:28:40: Your 'Small Wood Plank' was destroyed!
Day 154, 18:28:50: Your 'Small Wood Plank' was destroyed!
Day 154, 18:47:14: Your maps (Ramshackle Sloop) was destroyed by DVS DATBUTT at B4 [Long: -79.30 / Lat: 50.97]!
Day 154, 18:47:22: Bed -491537755 was removed from the Company!
Day 154, 18:48:44: Crew member Someone - Lvl 47 was killed by a Shark - Lvl 24!
Day 154, 19:29:32: Crew member Someone - Lvl 38 was killed by a Shark - Lvl 2!
Day 154, 19:37:20: Crew member Someone - Lvl 51 was killed by a Shark - Lvl 24!
Day 154, 19:55:06: Crew member Someone - Lvl 49 was killed by a Shark - Lvl 11!
Day 155, 15:03:31: Crew member Someone - Lvl 35 was killed!
Day 155, 15:35:35: Someone was added to the Company by Someone!
Day 155, 15:36:59: Crew member Someone - Lvl 52 was killed!
`, `**Some Tribe (1234567892)**
Day 155, 17:34:19: Someone Tamed a Monkey - Lvl 2 (Monkey)!
Day 155, 18:57:38: Crew member Someone - Lvl 52 was killed!
Day 155, 19:10:48: Someone set to Rank Group Leadership!
Day 155, 01:05:32: Someone claimed 'Filthy Milker - Lvl 35 (Cow)'!
`, `**Some Tribe (1234567892)**
Day 155, 01:24:58: Crew member Someone - Lvl 47 was killed by Tiffany - Lvl 39 (Bear) (DVS Gaming)!
Day 155, 01:41:56: Your Company killed Filthy Chick Two - Lvl 4 (Chicken) (The Filthy Few)!
Day 155, 01:43:00: Your Company killed Filthy Chick - Lvl 7 (Chicken) (The Filthy Few)!
Day 155, 22:00:16: Crew member Someone - Lvl 52 was killed by a Cobra - Lvl 4!
Day 155, 03:40:57: Someone was promoted to a Company Admin by Someone!
`, `**Some Tribe (1234567892)**
Day 155, 22:43:25: Crew member Someone - Lvl 35 was killed!
Day 155, 21:41:47: Crew member Someone - Lvl 39 was killed!
Day 156, 01:17:44: Crew member Someone - Lvl 44 was killed!
Day 155, 05:35:54: Crew member Someone - Lvl 51 was killed!
Day 155, 05:55:36: Your Filthy Milker - Lvl 36 (Cow) was killed by Sabel's Taming Bear - Lvl 66 (Bear) (DVS Gaming) at B4 [Long: -78.95 / Lat: 50.33]!
Day 156, 00:52:08: Crew member Someone - Lvl 52 was killed!
`, `**Some Tribe (1234567892)**
Day 156, 01:35:17: Someone Tamed a Bear - Lvl 40 (Bear)!
Day 155, 08:53:51: Crew member Someone - Lvl 51 was killed by Oponn Silverback - Lvl 32 (DVS Gaming - Citizen)!
`, `**Some Tribe (1234567892)**
Day 155, 10:56:35: Your 'Bed' was destroyed!
Day 155, 10:56:38: Bed -1565792037 was removed from the Company!
Day 155, 11:20:32: Crew member Someone - Lvl 51 was killed!
Day 156, 06:52:54: Crew member Someone - Lvl 36 was killed!
Day 156, 05:55:33: Crew member Jata Bighead - Lvl 50 was killed!
Day 156, 07:12:29: Crew member Someone - Lvl 44 was killed!
Day 156, 07:15:16: Crew member Someone - Lvl 36 was killed!
Day 155, 10:26:59: Crew member Someone - Lvl 52 was killed by a Shark - Lvl 10!
Day 155, 13:06:49: Crew member Someone - Lvl 51 was killed!
`, `**Some Tribe (1234567892)**
Day 155, 14:04:00: Crew member Someone - Lvl 44 was killed by a Lion - Lvl 3!
Day 156, 10:58:12: Someone demolished a 'Large Storage Box' at O6 - Lawless Region [Long: 93.34 / Lat: 29.81]!
Day 156, 11:02:54: Someone demolished a 'Wooden Chair' at O6 - Lawless Region [Long: 93.34 / Lat: 29.81]!
Day 156, 09:13:52: Crew member Someone - Lvl 44 was killed by a Pig - Lvl 16!
Day 156, 11:16:55: Someone demolished a 'Stone Square Ceiling' at O6 - Lawless Region [Long: 93.34 / Lat: 29.81]!
Day 156, 11:28:06: Someone demolished a 'Tannery' at O6 - Lawless Region [Long: 93.34 / Lat: 29.81]!
Day 156, 11:31:40: Someone demolished a 'Stone Square Ceiling' at O6 - Lawless Region [Long: 93.34 / Lat: 29.81]!
Day 155, 15:04:57: Bed -889053690 was added to the Company!
Day 155, 15:19:03: Crew member Someone - Lvl 47 was killed by a Bear - Lvl 1!
Day 156, 12:12:32: Someone demolished a 'Stone Water Reservoir' at O6 - Lawless Region [Long: 93.34 / Lat: 29.81]!
Day 156, 12:15:49: Someone demolished a 'Stone Square Ceiling' at O6 - Lawless Region [Long: 93.34 / Lat: 29.81]!
Day 156, 12:17:04: Someone demolished a 'Stone Square Ceiling' at O6 - Lawless Region [Long: 93.34 / Lat: 29.81]!
Day 156, 12:19:12: Someone demolished a 'Bookshelf' at O6 - Lawless Region [Long: 93.34 / Lat: 29.81]!
Day 156, 12:26:06: Someone demolished a 'Large Crop Plot' at O6 - Lawless Region [Long: 93.34 / Lat: 29.81]!
Day 156, 12:31:41: Someone demolished a 'Large Crop Plot' at O6 - Lawless Region [Long: 93.35 / Lat: 29.81]!
`, `**Some Tribe (1234567892)**
Day 155, 16:59:18: Crew member Someone - Lvl 36 was killed!
Day 155, 16:00:26: Crew member Someone - Lvl 39 was killed by a Shark - Lvl 4!
Day 156, 12:33:35: Crew member Someone - Lvl 47 was killed!
Day 155, 19:06:46: Crew member Someone - Lvl 51 was killed by a Wolf - Lvl 19!
`, `**Some Tribe (1234567892)**
Day 156, 15:28:08: Crew member Someone - Lvl 47 was killed!
Day 155, 18:57:57: Crew member Someone - Lvl 36 was killed by a Pig - Lvl 2!
Day 155, 21:04:05: Crew member Someone - Lvl 36 was killed by a Cobra - Lvl 9!
Day 156, 16:13:26: Your Diddy Kong - Lvl 16 (Monkey) was killed at N6 [Long: 81.51 / Lat: 28.44]!
Day 155, 22:56:54: Crew member Someone - Lvl 44 was killed by a Shark - Lvl 4!
`, `**Some Tribe (1234567892)**
Day 156, 21:54:17: Crew member Someone - Lvl 44 was killed!
Day 156, 20:12:13: Someone (Sometribe - Member) - Lvl 52 has become a settler in your Settlement 'Home Place' (N6) [ Longitude: 81.50 / Latitude: 28.37 ]!
Day 156, 20:12:13: Someone (Sometribe - Member) - Lvl 52 has become a settler in your Settlement 'Home Place' (N6) [ Longitude: 81.50 / Latitude: 28.37 ]!
Day 156, 20:12:13: Someone (Sometribe - Member) - Lvl 52 has become a settler in your Settlement 'Home Place' (N6) [ Longitude: 81.50 / Latitude: 28.37 ]!
Day 156, 20:12:13: Someone (Sometribe - Member) - Lvl 52 has become a settler in your Settlement 'Home Place' (N6) [ Longitude: 81.50 / Latitude: 28.37 ]!
Day 156, 20:12:13: Someone (Sometribe - Member) - Lvl 52 has become a settler in your Settlement 'Home Place' (N6) [ Longitude: 81.50 / Latitude: 28.37 ]!
Day 156, 20:12:13: Someone (Sometribe - Member) - Lvl 52 has become a settler in your Settlement 'Home Place' (N6) [ Longitude: 81.50 / Latitude: 28.37 ]!
Day 156, 20:12:13: Someone (Sometribe - Member) - Lvl 52 has become a settler in your Settlement 'Home Place' (N6) [ Longitude: 81.50 / Latitude: 28.37 ]!
Day 156, 20:12:13: Someone (Sometribe - Member) - Lvl 52 has become a settler in your Settlement 'Home Place' (N6) [ Longitude: 81.50 / Latitude: 28.37 ]!
Day 156, 20:12:13: Someone (Sometribe - Member) - Lvl 52 has become a settler in your Settlement 'Home Place' (N6) [ Longitude: 81.50 / Latitude: 28.37 ]!
Day 156, 20:12:13: Someone (Sometribe - Member) - Lvl 52 has become a settler in your Settlement 'Home Place' (N6) [ Longitude: 81.50 / Latitude: 28.37 ]!
Day 156, 20:12:13: Someone (Sometribe - Member) - Lvl 52 has become a settler in your Settlement 'Home Place' (N6) [ Longitude: 81.50 / Latitude: 28.37 ]!
Day 156, 20:12:13: Someone (Sometribe - Member) - Lvl 52 has become a settler in your Settlement 'Home Place' (N6) [ Longitude: 81.50 / Latitude: 28.37 ]!
Day 156, 20:12:13: Someone (Sometribe - Member) - Lvl 52 has become a settler in your Settlement 'Home Place' (N6) [ Longitude: 81.50 / Latitude: 28.37 ]!
Day 156, 20:12:13: Someone (Sometribe - Member) - Lvl 52 has become a settler in your Settlement 'Home Place' (N6) [ Longitude: 81.50 / Latitude: 28.37 ]!
Day 156, 20:12:13: Someone (Sometribe - Member) - Lvl 52 has become a settler in your Settlement 'Home Place' (N6) [ Longitude: 81.50 / Latitude: 28.37 ]!
Day 156, 20:12:13: Someone (Sometribe - Member) - Lvl 52 has become a settler in your Settlement 'Home Place' (N6) [ Longitude: 81.50 / Latitude: 28.37 ]!
Day 156, 20:12:13: Someone (Sometribe - Member) - Lvl 52 has become a settler in your Settlement 'Home Place' (N6) [ Longitude: 81.50 / Latitude: 28.37 ]!
Day 156, 20:12:13: Someone (Sometribe - Member) - Lvl 52 has become a settler in your Settlement 'Home Place' (N6) [ Longitude: 81.50 / Latitude: 28.37 ]!
Day 156, 20:12:13: Someone (Sometribe - Member) - Lvl 52 has become a settler in your Settlement 'Home Place' (N6) [ Longitude: 81.50 / Latitude: 28.37 ]!
Day 156, 20:12:13: Someone (Sometribe - Member) - Lvl 52 has become a settler in your Settlement 'Home Place' (N6) [ Longitude: 81.50 / Latitude: 28.37 ]!
Day 156, 20:12:13: Someone (Sometribe - Member) - Lvl 52 has become a settler in your Settlement 'Home Place' (N6) [ Longitude: 81.50 / Latitude: 28.37 ]!
Day 156, 20:12:13: Someone (Sometribe - Member) - Lvl 52 has become a settler in your Settlement 'Home Place' (N6) [ Longitude: 81.50 / Latitude: 28.37 ]!
Day 156, 20:12:13: Someone (Sometribe - Member) - Lvl 52 has become a settler in your Settlement 'Home Place' (N6) [ Longitude: 81.50 / Latitude: 28.37 ]!
Day 156, 20:12:13: Someone (Sometribe - Member) - Lvl 52 has become a settler in your Settlement 'Home Place' (N6) [ Longitude: 81.50 / Latitude: 28.37 ]!
Day 156, 01:29:28: Someone - Lvl 52 destroyed their 'Bed ()')!
Day 156, 01:32:26: Someone - Lvl 52 destroyed their 'Bed ()')!
Day 155, 23:49:09: Crew member Someone - Lvl 36 was killed!
Day 156, 02:53:04: Crew member Someone - Lvl 52 was killed by a Wolf - Lvl 19!
`, `**Some Tribe (1234567892)**
Day 156, 16:57:35: Your 'Bed' was destroyed!
Day 156, 16:57:38: Bed 1753491359 was removed from the Company!
Day 156, 04:37:03: Crew member Someone - Lvl 36 was killed by a Pig - Lvl 4!
Day 156, 06:08:50: Crew member Someone - Lvl 44 was killed by a Wolf - Lvl 4!
Day 156, 23:41:27: Crew member Someone - Lvl 44 was killed!
`, `**Some Tribe (1234567892)**
Day 156, 14:15:02: Crew member Someone - Lvl 36 was killed by a Cobra - Lvl 2!
Day 156, 05:22:00: Crew member Someone - Lvl 44 was killed!
Day 156, 05:33:58: Crew member Someone - Lvl 51 was killed!
Day 156, 18:31:53: Crew member Someone - Lvl 36 was killed by a Lion - Lvl 19!
Day 156, 22:31:12: Crew member Someone - Lvl 36 was killed by a Cobra - Lvl 9!
`, `**Some Tribe (1234567892)**
Day 156, 09:19:49: Crew member Someone - Lvl 52 was killed by a Lion - Lvl 5!
Day 157, 01:08:51: Crew member Someone - Lvl 36 was killed by a Pig - Lvl 8!
Day 157, 06:02:22: Crew member Someone - Lvl 44 was killed by a Shark - Lvl 2!
`, `**Some Tribe (1234567892)**
Day 156, 10:40:46: Crew member Someone - Lvl 38 was killed by a Pig - Lvl 13!
Day 157, 01:27:19: Crew member Someone - Lvl 37 was killed by a Crocodile - Lvl 15!
Day 157, 02:25:48: Crew member Someone - Lvl 44 was killed by a Cobra - Lvl 30!
Day 157, 02:27:38: Crew member Someone - Lvl 37 was killed by a Rhino - Lvl 9!
Day 157, 04:38:11: Crew member Someone - Lvl 37 was killed by a Cobra - Lvl 6!`}
