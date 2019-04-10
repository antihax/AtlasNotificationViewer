package mapserver

import (
	"strings"

	"github.com/go-redis/redis"
)

func (s *MapServer) scanHash(pattern string, maxItems int64) (map[string]map[string]string, error) {
	records := make(map[string]map[string]string)

	// [GSG] Scan is slower than Keys but provides gaps for other things to execute
	var keys []string
	iter := s.redisClient.Scan(0, pattern, maxItems).Iterator()
	for iter.Next() {
		keys = append(keys, iter.Val())
	}

	// If the iterator has an error, do not continue to pull the hash
	if err := iter.Err(); err != nil {
		return nil, err
	}

	// Build a pipeline of requests for each key
	results := make(map[string]*redis.StringStringMapCmd)
	pipe := s.redisClient.Pipeline()
	for _, id := range keys {
		results[id] = pipe.HGetAll(id)
	}

	// Execute all requests
	if _, err := pipe.Exec(); err != nil {
		return nil, err
	}

	// Build a map of the results
	for _, id := range keys {
		key := id
		if strings.Contains(id, ":") {
			parts := strings.Split(id, ":")
			key = parts[1]
		}
		var err error
		records[key], err = results[id].Result()
		if err != nil {
			return nil, err
		}
	}

	return records, nil
}
