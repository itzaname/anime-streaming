package model

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/itzaname/anime-streaming/site/app/lib/database"
)

// A download is a per torrent job that may contain one or more episodes

// Download contains download data
type Download struct {
	ID       int
	Anime    int
	Episodes []int
	Count    int // Ammount of files MINIMUM for us to grab
	Name     string
	Search   string
	Nyaa     string
	Magnet   string
	Complete bool
}

// DownloadByID retrieves an dl by id
func DownloadByID(id int) (Download, error) {
	var download Download
	data, err := database.DB.Get(fmt.Sprintf("download:%d", id)).Result()
	if err != nil {
		return download, err
	}
	err = json.Unmarshal([]byte(data), &download)
	if err != nil {
		return download, err
	}

	return download, nil
}

// DownloadByAnime retrieves an dl by anime
func DownloadByAnime(id int, active bool) (Download, error) {
	downloads, err := DownloadList()
	if err != nil {
		return Download{}, err
	}

	for _, dl := range downloads {
		if dl.Anime == id {
			if active && !dl.Complete {
				return dl, err
			} else if !active {
				return dl, err
			}
		}
	}

	return Download{}, fmt.Errorf("Not found")
}

// DownloadExistsByAnime retrieves an active state by anime
func DownloadExistsByAnime(id int, active bool) (bool, error) {
	downloads, err := DownloadList()
	if err != nil {
		return false, err
	}

	for _, dl := range downloads {
		if dl.Anime == id {
			if active && !dl.Complete {
				return true, err
			} else if !active {
				return true, err
			}
		}
	}

	return false, fmt.Errorf("Not found")
}

// DownloadAdd adds new download
func DownloadAdd(anime int) (Download, error) {
	var download Download
	id, err := database.DB.Incr("download:total").Result()
	if err != nil {
		return download, err
	}
	download.ID = int(id)
	download.Anime = anime
	err = download.Update()
	if err != nil {
		return download, err
	}
	return download, database.DB.LPush("download:id", download.ID).Err()
}

// DownloadIDList returns a list of all download ids
func DownloadIDList() ([]int, error) {
	var list []int
	data, err := database.DB.LRange("download:id", 0, -1).Result()
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

// DownloadList returns all the downloads
func DownloadList() ([]Download, error) {
	var download []Download
	list, err := DownloadIDList()
	if err != nil {
		return download, err
	}

	for i := 0; i < len(list); i++ {
		an, err := DownloadByID(list[i])
		if err != nil {
			return download, err
		}
		download = append(download, an)
	}
	return download, nil
}

// DownloadListByAnime returns all the downloads for an anime
func DownloadListByAnime(id int) ([]Download, error) {
	var download []Download
	list, err := DownloadIDList()
	if err != nil {
		return download, err
	}

	for i := 0; i < len(list); i++ {
		an, err := DownloadByID(list[i])
		if err != nil {
			return download, err
		}
		if an.Anime == id {
			download = append(download, an)
		}
	}
	return download, nil
}

// Key returns the database key used for object
func (d *Download) Key() string {
	return fmt.Sprintf("download:%d", d.ID)
}

// Update will update the current user object in the redis store
func (d *Download) Update() error {
	data, err := json.Marshal(d)
	if err != nil {
		return err
	}

	return database.DB.Set(d.Key(), data, -1).Err()
}

// GetEpisodes returns all episodes under the download
func (d *Download) GetEpisodes() ([]Episode, error) {
	episodes := []Episode{}
	for _, id := range d.Episodes {
		episode, err := EpisodeById(id)
		if err != nil {
			return episodes, err
		}

		episodes = append(episodes, episode)
	}
	return episodes, nil
}
