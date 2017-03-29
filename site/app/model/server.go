package model

import (
	"encoding/json"
	"fmt"
	"github.com/itzaname/anime-streaming/site/app/lib/database"
)

type Server struct {
	ID       string
	Volume   string
	IP       string
	Address  string
	Encoding bool
}

// ServerData returns a servers data
func ServerData(encoder bool) (Server, error) {
	var server Server

	key := "server:downloader"
	if encoder {
		key = "server:encoder"
	}

	data, err := database.DB.Get(key).Result()
	if err != nil {
		return server, err
	}

	err = json.Unmarshal([]byte(data), &data)
	if err != nil {
		return server, err
	}

	return server, nil
}

// ServerDelete deleted a server
func ServerDelete(encoder bool) error {
	key := "server:downloader"
	if encoder {
		key = "server:encoder"
	}

	return database.DB.Del(key).Err()
}

// ServerExists if server exists
func ServerExists(encoder bool) (bool, error) {
	key := "server:downloader"
	if encoder {
		key = "server:encoder"
	}

	return database.DB.Exists(key).Result()
}

// ServerAdd add a new server
func ServerAdd(id, volume, ip, addr string, encoder bool) error {
	exists, err := ServerExists(encoder)
	if err != nil || exists {
		return fmt.Errorf("Exists %s", err)
	}

	key := "server:downloader"
	if encoder {
		key = "server:encoder"
	}

	server := Server{
		ID:       id,
		Volume:   volume,
		IP:       ip,
		Address:  addr,
		Encoding: encoder,
	}

	data, err := json.Marshal(server)
	if err != nil {
		return err
	}

	return database.DB.Set(key, data, -1).Err()
}
