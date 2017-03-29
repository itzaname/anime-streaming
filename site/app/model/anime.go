package model

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/itzaname/anime-streaming/site/app/lib/database"
)

const (
	StateReady       = 0
	StateInOperation = 0
)

// Info contains data
type Info struct {
	Title        string
	Episodes     int
	Image        string
	Type         int
	Status       int
	StatusString string
	Start        string
	End          string
	Synonyms     string
}

// Anime contain anime data
type Anime struct {
	ID         int
	State      int
	Downloaded int
	Achieved   bool // If we are hosting part of this in cold storage
	Info       Info
	Episodes   []int
	Complete   bool
	Manual     bool // When marked manual we will not automatically check for torrents
}

// AnimeByID retrieves an anime by its MAL id
func AnimeByID(id int) (Anime, error) {
	var anime Anime
	data, err := database.DB.Get(fmt.Sprintf("anime:%d", id)).Result()
	if err != nil {
		return anime, err
	}
	err = json.Unmarshal([]byte(data), &anime)
	if err != nil {
		return anime, err
	}

	return anime, nil
}

// AddAnime adds new anime
func AddAnime(id int) (Anime, error) {
	var anime Anime
	anime.ID = id
	err := anime.Update()
	if err != nil {
		return anime, err
	}
	return anime, database.DB.LPush("anime:id", anime.ID).Err()
}

// AnimeExists returns if anime is in db
func AnimeExists(id int) (bool, error) {
	return database.DB.Exists(fmt.Sprintf("anime:%d", id)).Result()
}

// AnimeCreate creates an anime
func AnimeCreate(id int) (Anime, error) {
	if database.DB.Exists(fmt.Sprintf("anime:%d", id)).Val() {
		return AnimeByID(id)
	}
	return AddAnime(id)
}

// AllAnimeID returns a list of all anime ids
func AllAnimeID() ([]int, error) {
	var list []int
	data, err := database.DB.LRange("anime:id", 0, -1).Result()
	if err != nil {
		return list, err
	}

	for i := 0; i < len(data); i++ {
		parse, err := strconv.ParseInt(data[i], 10, 32)
		if err != nil {
			return list, err
		}
		list = append(list, int(parse))
	}

	return list, nil
}

// AnimeList returns all the anime in the db
func AnimeList() ([]Anime, error) {
	var anime []Anime
	list, err := AllAnimeID()
	if err != nil {
		return anime, err
	}

	for i := 0; i < len(list); i++ {
		an, err := AnimeByID(list[i])
		if err != nil {
			return anime, err
		}
		anime = append(anime, an)
	}
	return anime, nil
}

// Key returns the database key used for object
func (a *Anime) Key() string {
	return fmt.Sprintf("anime:%d", a.ID)
}

// Update will update the current user object in the redis store
func (a *Anime) Update() error {
	data, err := json.Marshal(a)
	if err != nil {
		return err
	}

	return database.DB.Set(a.Key(), data, -1).Err()
}
