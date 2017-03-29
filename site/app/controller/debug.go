package controller

import (
	"github.com/itzaname/anime-streaming/site/app/model"
	"net/http"
)

func DebugKey(w http.ResponseWriter, r *http.Request) {
	invite, err := model.InviteCreate()
	if err != nil {
		http.Error(w, "Internal server error.", 500)
		return
	}

	w.Write([]byte(invite))
}
