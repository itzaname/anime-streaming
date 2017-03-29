package controller

import (
	"net/http"

	"fmt"
	"strings"

	"github.com/itzaname/anime-streaming/site/app/lib/session"
	"github.com/itzaname/anime-streaming/site/app/lib/view"
	"github.com/itzaname/anime-streaming/site/app/model"
)

func StreamHome(w http.ResponseWriter, r *http.Request) {
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

	trimTitle := func(title string) string {
		if len(title) > 50 {
			return strings.TrimSpace(title[:50]) + "..."
		}
		return title
	}

	for _, anime := range list.Anime {
		if anime.MyStatus != 1 && anime.MyRewatching == "0" {
			continue
		}
		entry := map[string]string{}
		entry["image"] = anime.SeriesImage
		entry["episodes"] = fmt.Sprintf("%d", anime.SeriesEpisodes)
		entry["status"] = getLabel(anime.MyStatus)
		entry["class"] = strings.ToLower(getLabel(anime.MyStatus))
		entry["title"] = trimTitle(anime.SeriesTitle)
		entry["id"] = fmt.Sprintf("%d", anime.SeriesAnimeDBID)
		entry["episode"] = fmt.Sprintf("%d", anime.MyWatchedEpisodes)
		entry["url"] = fmt.Sprintf("/watch/%d/%d", anime.SeriesAnimeDBID, anime.MyWatchedEpisodes+1)
		entry["titles"] = strings.Replace(anime.SeriesTitle+anime.SeriesSynonyms, ";", " ", -1)
		if anime.MyWatchedEpisodes >= anime.SeriesEpisodes {
			entry["url"] = fmt.Sprintf("/watch/%d/%d", anime.SeriesAnimeDBID, anime.MyWatchedEpisodes)
		}

		data = append(data, entry)
	}

	v := view.New(r)
	v.Name = "browse/home"
	v.Vars["anime"] = data
	v.Render(w)
}

func StreamSearch(w http.ResponseWriter, r *http.Request) {
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

	trimTitle := func(title string) string {
		if len(title) > 35 {
			return strings.TrimSpace(title[:35]) + "..."
		}
		return title
	}

	for _, anime := range list.Anime {
		entry := map[string]string{}
		entry["image"] = anime.SeriesImage
		entry["episodes"] = fmt.Sprintf("%d", anime.SeriesEpisodes)
		entry["status"] = getLabel(anime.MyStatus)
		entry["class"] = strings.ToLower(getLabel(anime.MyStatus))
		entry["title"] = trimTitle(anime.SeriesTitle)
		entry["id"] = fmt.Sprintf("%d", anime.SeriesAnimeDBID)
		entry["episode"] = fmt.Sprintf("%d", anime.MyWatchedEpisodes)
		entry["url"] = fmt.Sprintf("/watch/%d/%d", anime.SeriesAnimeDBID, anime.MyWatchedEpisodes+1)
		entry["titles"] = strings.Replace(anime.SeriesTitle+anime.SeriesSynonyms, ";", " ", -1)
		if anime.MyWatchedEpisodes >= anime.SeriesEpisodes {
			entry["url"] = fmt.Sprintf("/watch/%d/%d", anime.SeriesAnimeDBID, anime.MyWatchedEpisodes)
		}

		data = append(data, entry)
	}

	v := view.New(r)
	v.Name = "browse/search"
	v.Vars["anime"] = data
	v.Render(w)
}
