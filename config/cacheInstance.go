package config

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"os"
	"sync"
)

var cacheLock = &sync.Mutex{}
var cacheConnection redis.Conn

func GetCacheInstance() redis.Conn {
	//To prevent expensive lock operations
	//This means that the cacheConnection field is already populated
	if cacheConnection == nil {
		cacheLock.Lock()
		defer cacheLock.Unlock()

		//Only one goroutine can create the singleton instance.
		if cacheConnection == nil {
			log.Println("Creating Cache instance")

			conn, err := redis.DialURL(os.Getenv("REDIS_URL"), redis.DialTLSSkipVerify(true))
			if err != nil {
				log.Printf("Cache initialization: %s\n", err.Error())
				os.Exit(1)
			}

			cacheConnection = conn
		} else {
			log.Println("Application Cache instance already created!")
		}
	} else {
		//log.Println("Application Cache instance already created!")
	}

	return cacheConnection
}
