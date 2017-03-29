package myanimelist

import (
	"fmt"

	"github.com/itzaname/go-myanimelist/mal"
)

// Configure sets login data
func Configure() {
	c := mal.NewClient(nil)
	C = c
}

// C MAL client
var C *mal.Client

// Search check MAL with string
func Search(user string, pass string, query string) (*mal.AnimeResult, error) {
	C.SetCredentials(user, pass)

	result, _, err := C.Anime.Search(query)
	return result, err
}

// Add adds anime to list
func Add(user string, pass string, id int, entry mal.AnimeEntry) error {
	C.SetCredentials(user, pass)

	_, err := C.Anime.Add(id, entry)
	return err
}

// SetEpisode sets the current watching episode
func SetEpisode(user string, pass string, ep int, id int, status int) error {
	C.SetCredentials(user, pass)

	var entry mal.AnimeEntry
	entry.Episode = ep
	entry.Status = status
	_, err := C.Anime.Update(id, entry)
	return err
}

// GetList returns MAL anime list
func GetList(user string, pass string) (*mal.AnimeList, error) {
	C.SetCredentials(user, pass)

	list, _, err := C.Anime.List(C.Username)
	return list, err
}

// GetAnimeInfo gets MAL anime info
func GetAnimeInfo(user string, pass string, name string, id int) (mal.AnimeRow, error) {
	C.SetCredentials(user, pass)

	result, _, err := C.Anime.Search(name)
	if err != nil {
		return mal.AnimeRow{}, err
	}
	for _, anime := range result.Rows {
		if anime.ID == id {
			return anime, nil
		}
	}
	return mal.AnimeRow{}, fmt.Errorf("Not found!")
}
