package controller

import (
	"encoding/json"
	"net/http"

	"github.com/itzaname/anime-streaming/site/app/model"
	"github.com/itzaname/anime-streaming/site/app/lib/log"
)

func ServerAPIDownloadList(w http.ResponseWriter, r *http.Request) {
	downloads, err := model.DownloadList()
	if err != nil {
		log.Errorln(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	data := []map[string]interface{}{}

	log.Debugln(len(downloads))

	for _, dl := range downloads {
		anime, err := model.AnimeByID(dl.Anime)
		if err != nil {
			log.Errorln("error")
			http.Error(w, "Internal Server Error", 500)
			return
		}

		info := map[string]interface{}{}
		info["torrent"] = dl.Nyaa
		info["title"] = anime.Info.Title
		info["episodes"] = dl.Count
		info["download"] = dl.ID
		info["search"] = dl.Search
		info["magnet"] = dl.Magnet

		data = append(data, info)
	}

	serialized, err := json.Marshal(data)
	if err != nil {
		log.Errorln(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(serialized)
}
