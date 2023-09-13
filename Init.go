package lamb

import (
	"errors"
	"fmt"
	"os"
	"time"
)

// Init initializes the lamb module.
func Init(config map[interface{}]interface{}) error {
	// validate the config
	if config["lamb"] == nil {
		return errors.New("lamb: missing config")
	}

	lambConfig, ok := config["lamb"].(map[interface{}]interface{})

	if !ok {
		return fmt.Errorf("lamb: config must be a map[interface{}]interface{} but got %T", config["lamb"])
	}

	dir, exists := lambConfig["dir"]

	if !exists {
		return errors.New("lamb: missing config: dir")
	}

	if _, ok := dir.(string); !ok {
		return errors.New("lamb: dir must be a string")
	}

	// validate the cache
	cache, exists := lambConfig["cache"]

	if exists {
		cacheMap, ok := cache.(map[interface{}]interface{})

		if !ok {
			return errors.New("lamb: cache must be a map[interface{}]interface{}")
		}

		// get the dir (optional) and cache time (required)
		dir, dirExists := cacheMap["dir"]
		cacheTime, timeExists := cacheMap["time"]

		if _, ok := dir.(string); !ok {
			return errors.New("lamb: cache: dir must be a string")
		}

		if !dirExists {
			// default to .cache
			dir = ".cache"
		}

		if !timeExists {
			return errors.New("lamb: cache: missing config: time")
		}

		if _, ok := cacheTime.(string); !ok {
			return errors.New("lamb: cache: time must be a string")
		}

		cacheTimeDuration, err := time.ParseDuration(cacheTime.(string))

		if err != nil {
			return errors.New("lamb: cache: time must be a valid duration")
		}

		os.Setenv("GOVEL_LAMB_CACHE_DIR", dir.(string))
		os.Setenv("GOVEL_LAMB_CACHE_TIME", cacheTimeDuration.String())
	}

	// set var in the environment
	os.Setenv("GOVEL_LAMB_BASE_DIR", dir.(string))

	return nil
}
