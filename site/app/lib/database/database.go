package database

import redis "gopkg.in/redis.v5"

// DB redis database instance
var DB *redis.Client

// Init starts database connection and sets package db variable
func Init(info Info) error {
	client := redis.NewClient(&redis.Options{
		Addr:     info.Address + ":" + info.Port,
		Password: info.Password,
		DB:       info.Database,
	})
	DB = client
	return client.Ping().Err()
}

// Info config data
type Info struct {
	Address  string
	Port     string
	Password string
	Database int
}
