package get_with_cache_go

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// GetDataFunc is a generic type for functions that return a value of type T and an error.
// T must be a type that can be marshaled and unmarshaled by the encoding/json package.
type GetDataFunc[T any] func() (T, error)

// FetchDataWithCache is a generic function that works with any data type T.
// It checks for cached data in a file named `<cacheKey>.json` within `cacheDir`.
// If the cache exists, it returns the cached data.
// If not, it calls `getDataFunc` to fetch the data, caches it, and returns the data.
func FetchDataWithCache[T any](getDataFunc GetDataFunc[T], cacheKey string, cacheDir string) (T, error) {
	var data T
	cacheFilePath := filepath.Join(cacheDir, cacheKey+".json")

	// Check if the cache file exists
	if _, err := os.Stat(cacheFilePath); err == nil {
		// Cache file exists, read and unmarshal it
		fileData, err := os.ReadFile(cacheFilePath)
		if err != nil {
			return data, fmt.Errorf("error reading cache file: %w", err)
		}

		if err := json.Unmarshal(fileData, &data); err != nil {
			return data, fmt.Errorf("error parsing cache file JSON: %w", err)
		}

		return data, nil
	}

	// Cache file does not exist, call getDataFunc to get the data
	data, err := getDataFunc()
	if err != nil {
		return data, fmt.Errorf("error fetching data: %w", err)
	}

	// Marshal the data and save it to cache
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return data, fmt.Errorf("error marshaling data to JSON: %w", err)
	}

	if err := os.WriteFile(cacheFilePath, dataBytes, 0644); err != nil {
		return data, fmt.Errorf("error writing cache file: %w", err)
	}

	return data, nil
}
