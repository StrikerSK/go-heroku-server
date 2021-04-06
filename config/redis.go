package config

import (
	"github.com/gomodule/redigo/redis"
	"os"
)

var Cache redis.Conn

func InitializeRedis() {
	conn, err := redis.DialURL(os.Getenv("REDIS_URL"), redis.DialTLSSkipVerify(true))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Assign the connection to the package level `cache` variable
	Cache = conn
}
