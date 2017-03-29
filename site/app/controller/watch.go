package controller

import (
	"github.com/go-zoo/bone"
	"github.com/itzaname/anime-streaming/site/app/lib/view"
	"github.com/itzaname/anime-streaming/site/app/model"
	"net/http"
	"strconv"
)

func WatchVideo(w http.ResponseWriter, r *http.Request) {
	sid := bone.GetValue(r, "id")
	sep := bone.GetValue(r, "ep")

	id, err := strconv.ParseInt(sid, 10, 32)
	if err != nil {
		http.Error(w, "Internal server error.", 500)
		return
	}
	ep, err := strconv.ParseInt(sep, 10, 32)
	if err != nil {
		http.Error(w, "Internal server error.", 500)
		return
	}

	anime, err := model.AnimeByID(int(id))
	if err != nil {
		http.Error(w, "Anime not found.", 404)
		return
	}

	v := view.New(r)
	v.Name = "watch/download"
	v.Vars["episode"] = ep
	v.Vars["id"] = id
	v.Vars["anime"] = anime
	v.Render(w)
}
