package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/itzaname/anime-streaming/site/app/lib/database"
	"github.com/itzaname/anime-streaming/site/app/lib/log"
	"github.com/itzaname/anime-streaming/site/app/lib/myanimelist"
	"github.com/itzaname/go-myanimelist/mal"
	"strings"
)

//var mutex sync.RWMutex

func updateData(user string, pass string, info mal.Anime) {
	exists, err := AnimeExists(info.SeriesAnimeDBID)
	if err != nil || exists {
		//log.Debugln("Existing anime Exiting...")
		return
	}

	anime, err := AnimeCreate(info.SeriesAnimeDBID)
	if err == nil {
		anime.Info.Episodes = info.SeriesEpisodes
		anime.Info.Image = info.SeriesImage
		anime.Info.Title = info.SeriesTitle
		anime.Info.Start = info.SeriesStart
		anime.Info.End = info.SeriesEnd
		anime.Info.Status = info.SeriesStatus
		anime.Info.Synonyms = info.SeriesSynonyms
		if len(anime.Info.Synonyms) > 0 && anime.Info.Synonyms[0] == ';' {
			anime.Info.Synonyms = strings.TrimSpace(anime.Info.Synonyms[1:])
		}
		switch info.SeriesStatus {
		case 1:
			anime.Info.StatusString = "Currently Airing"
			break
		case 2:
			anime.Info.StatusString = "Finished Airing"
			break
		case 3:
			anime.Info.StatusString = "Not yet aired"
			break
		default:
			anime.Info.StatusString = "Unkown"
			break
		}

		/*if anime.Info.Synopsis == "" {
			mutex.Lock()
			log.Debugln("Importing extra MAL data for anime:",info.SeriesTitle)
			malclient := myanimelist.NewClient(user, pass)
			data, err := malclient.GetAnimeInfo(info.SeriesTitle, info.SeriesAnimeDBID)
			if err == nil {
				anime.Info.Synopsis = data.Synopsis
			}
			mutex.Unlock()
		}*/
		anime.Update()
	}
}

// InvalidateListCache will invalidate the cached list for the user
func InvalidateListCache(user string) error {
	return database.DB.Del(fmt.Sprintf("mal:%s", user)).Err()
}

// GetAnimeList returns a the users mal list with a 60s cache
func GetAnimeList(user string, password string) (*mal.AnimeList, error) {
	result, err := database.DB.Get(fmt.Sprintf("mal:%s", user)).Result()
	if err != nil {
		malclient := myanimelist.NewClient(user, password)
		list, err := malclient.GetList()
		if err != nil {
			return &mal.AnimeList{}, err
		}

		data, err := json.Marshal(list)
		if err != nil {
			return &mal.AnimeList{}, err
		}

		for i := 0; i < len(list.Anime); i++ {
			go updateData(user, password, list.Anime[i])
		}

		database.DB.Set(fmt.Sprintf("mal:%s", user), data, time.Second*60)
		log.Debugf("GetAnimeList: %s's list was not cached. Caching for 60 seconds.\n", user)
		return list, nil
	}

	var list mal.AnimeList
	err = json.Unmarshal([]byte(result), &list)
	if err != nil {
		return &list, err
	}

	log.Debugf("GetAnimeList: %s's list was cached.\n", user)
	return &list, nil
}
