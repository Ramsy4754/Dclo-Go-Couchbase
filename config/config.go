package config

import (
	"errors"
	"log"
	"os"
	"sync"
)

type Config struct {
	RunEnv
	*CouchbaseConfig
}

var (
	config *Config
	once   sync.Once
)

func getRunEnv() (RunEnv, error) {
	runEnv := os.Getenv("RUN_ENV")
	if runEnv == "" {
		return NoEnv, errors.New("RUN_ENV is not set")
	} else if StringToRunEnv(runEnv) == NoEnv {
		return NoEnv, errors.New("invalid RUN_ENV: " + runEnv)
	}
	return StringToRunEnv(runEnv), nil
}

func getCouchbaseConfig() *CouchbaseConfig {
	url := os.Getenv("DCLO_COUCHBASE_URL")
	user := os.Getenv("DCLO_COUCHBASE_USERNAME")
	password := os.Getenv("DCLO_COUCHBASE_PASSWORD")

	if url == "" || user == "" || password == "" {
		return nil
	}
	return &CouchbaseConfig{
		url,
		user,
		password,
	}
}

func GetConfig() *Config {
	once.Do(func() {
		runEnv, err := getRunEnv()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("RUN_ENV: %s\n", runEnv)

		couchCfg := getCouchbaseConfig()
		if couchCfg == nil {
			log.Fatal("Couchbase configuration is not set")
		}

		copiedCfg := CouchbaseConfig{
			couchCfg.URL,
			couchCfg.User,
			"******",
		}
		log.Printf("Couchbase configuration: %+v\n", copiedCfg)
		config = &Config{
			runEnv,
			couchCfg,
		}
	})
	return config
}
