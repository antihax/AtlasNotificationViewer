package atlasadminserver

import "errors"

func (s *AtlasAdminServer) getPlayerIDFromSteamID(steamID string) (string, error) {
	s.steamDataLock.RLock()
	playerID, ok := s.steamData[steamID]
	s.steamDataLock.RUnlock()
	if !ok {
		return "", errors.New("No pathfinder found for SteamID")
	}
	return playerID, nil
}

func (s *AtlasAdminServer) getPlayerDataFromSteamID(steamID string) (map[string]string, error) {
	playerID, err := s.getPlayerIDFromSteamID(steamID)
	if err != nil {
		return nil, err
	}
	s.playerDataLock.RLock()
	defer s.playerDataLock.RUnlock()
	return s.playerData[playerID], nil
}

func (s *AtlasAdminServer) getTribeDataFromSteamID(steamID string) (map[string]string, error) {
	playerData, err := s.getPlayerDataFromSteamID(steamID)
	if err != nil {
		return nil, err
	}

	if _, ok := playerData["tribeID"]; ok {
		s.tribeDataLock.RLock()
		defer s.tribeDataLock.RUnlock()
		return s.tribeData[playerData["tribeID"]], nil
	}
	return nil, errors.New("No tribe information pathfinder")
}
