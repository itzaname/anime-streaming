package model

import (
	"encoding/json"
	"fmt"

	"github.com/itzaname/anime-streaming/site/app/lib/database"
)

type Media struct {
	ID    string
	Size  uint64
	Type  string // subtitle | audio | video
	Name  string
	Path  []string
	Flags int8
}

type Episode struct {
	ID       int
	Anime    int
	Episode  int
	Achieved bool // If this episode is in cold storage
	Type     int  // 1 Episodes 2 Extra 3 Op/ED
	Files    []Media
}

// EpisodeById returns an episode from its id
func EpisodeById(id int) (Episode, error) {
	var episode Episode
	data, err := database.DB.Get(fmt.Sprintf("episode:%d", id)).Result()
	if err != nil {
		return episode, err
	}
	err = json.Unmarshal([]byte(data), &episode)
	if err != nil {
		return episode, err
	}

	return episode, nil
}

// EpisodeAdd adds a new episode
func EpisodeAdd() (Episode, error) {
	var episode Episode
	id, err := database.DB.Incr("episode:total").Result()
	if err != nil {
		return episode, err
	}

	episode.ID = int(id)
	return episode, episode.Update()
}

// Key returns the database key used for object
func (a *Episode) Key() string {
	return fmt.Sprintf("episode:%d", a.ID)
}

// Update will update the current episode object in the redis store
func (a *Episode) Update() error {
	data, err := json.Marshal(a)
	if err != nil {
		return err
	}

	return database.DB.Set(a.Key(), data, -1).Err()
}
