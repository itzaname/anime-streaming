package controller

import (
	"encoding/json"
	"fmt"
	"github.com/itzaname/anime-streaming/site/app/lib/session"
	"github.com/itzaname/anime-streaming/site/app/model"
	"net/http"
	"strings"
)

func APIAnimeList(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)

	username := sess.Values["maluser"].(string)
	password := sess.Values["malpass"].(string)

	list, err := model.GetAnimeList(username, password)
	if err != nil {
		http.Error(w, "Internal server error.", 500)
		return
	}

	data := []map[string]string{}

	getLabel := func(status int) string {
		switch status {
		case 1:
			return "Watching"
		case 2:
			return "Completed"
		case 3:
			return "Hold"
		case 4:
			return "Dropped"
		case 6:
			return "Planned"
		}
		return "Error"
	}

	for _, anime := range list.Anime {
		entry := map[string]string{}
		entry["image"] = anime.SeriesImage
		entry["episodes"] = fmt.Sprintf("%d", anime.SeriesEpisodes)
		entry["status"] = getLabel(anime.MyStatus)
		entry["class"] = strings.ToLower(getLabel(anime.MyStatus))
		entry["title"] = anime.SeriesTitle
		entry["id"] = fmt.Sprintf("%d", anime.SeriesAnimeDBID)
		entry["episode"] = fmt.Sprintf("%d", anime.MyWatchedEpisodes)
		entry["url"] = fmt.Sprintf("/watch/%d/%d", anime.SeriesAnimeDBID, anime.MyWatchedEpisodes+1)
		entry["titles"] = strings.Replace(anime.SeriesTitle+anime.SeriesSynonyms, ";", " ", -1)
		if anime.MyWatchedEpisodes >= anime.SeriesEpisodes {
			entry["url"] = fmt.Sprintf("/watch/%d/%d", anime.SeriesAnimeDBID, anime.MyWatchedEpisodes)
		}

		data = append(data, entry)
	}

	serialized, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Internal server error.", 500)
		return
	}

	w.Write(serialized)
}
