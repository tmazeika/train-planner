package fetch

import (
	"encoding/gob"
	"errors"
	"os"
	"time"
)

var cacheColdMissErr = errors.New("cache cold miss")

func getCached(q Query) ([]Train, error) {
	cache, err := readCached()
	if err != nil {
		return nil, err
	}

	var trains []Train
	for _, t := range cache {
		if t.FromStation == q.FromStation &&
			t.ToStation == q.ToStation &&
			t.SameDay(q.When) {
			trains = append(trains, t)
		}
	}
	// If we don't find any trains for this day that match our query, then we
	// have a cold miss. If we find at least one train, then we have to have seen
	// all trains for that day, so we should be good to go.
	if len(trains) == 0 {
		return nil, cacheColdMissErr
	}
	return trains, nil
}

func saveCached(newTrains []Train) error {
	now := time.Now()
	// Read cache.
	cache, err := readCached()
	if err != nil && err != cacheColdMissErr {
		return err
	}
	// Clean cache.
	for key, t := range cache {
		if t.FromTime.Before(now.Truncate(24 * time.Hour)) {
			delete(cache, key)
		}
	}
	// Merge new trains into cache.
	for _, t := range newTrains {
		cache[t.hash()] = t
	}

	file, err := os.OpenFile(".trains.cache", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return nil
	}
	defer file.Close()

	// Encode cache.
	return gob.NewEncoder(file).Encode(cache)
}

func readCached() (map[string]Train, error) {
	file, err := os.OpenFile(".trains.cache", os.O_RDONLY, 0755)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]Train), cacheColdMissErr
		} else {
			return nil, err
		}
	}
	defer file.Close()

	// Decode cache.
	var cache map[string]Train
	if err := gob.NewDecoder(file).Decode(&cache); err != nil {
		return nil, err
	}
	return cache, nil
}
