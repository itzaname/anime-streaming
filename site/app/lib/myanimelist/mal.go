package myanimelist

import (
	"fmt"

	"github.com/itzaname/go-myanimelist/mal"
)

type Client mal.Client

// Search check MAL with string
func (C *Client) Search(query string) (*mal.AnimeResult, error) {
	result, _, err := C.Anime.Search(query)
	return result, err
}

// Add adds anime to list
func (C *Client) Add(id int, entry mal.AnimeEntry) error {
	_, err := C.Anime.Add(id, entry)
	return err
}

// SetEpisode sets the current watching episode
func (C *Client) SetEpisode(ep int, id int, status int) error {
	var entry mal.AnimeEntry
	entry.Episode = ep
	entry.Status = status
	_, err := C.Anime.Update(id, entry)
	return err
}

// GetList returns MAL anime list
func (C *Client) GetList() (*mal.AnimeList, error) {
	list, _, err := C.Anime.List(C.Username)
	return list, err
}

// GetAnimeInfo gets MAL anime info
func (C *Client) GetAnimeInfo(name string, id int) (mal.AnimeRow, error) {
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

// NewClient returns a new client for a user
func NewClient(user string, pass string) *Client {
	c := mal.NewClient(nil)
	c.SetCredentials(user, pass)
	client := Client(*c)
	return &client
}
