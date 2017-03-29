package middleware

import (
	"net/http"

	"github.com/itzaname/anime-streaming/site/app/controller"
	"github.com/itzaname/anime-streaming/site/app/lib/session"
	"github.com/itzaname/anime-streaming/site/app/lib/webserver"
)

func RequireAuth(w http.ResponseWriter, r *http.Request, next webserver.Next) {
	sess := session.Instance(r)

	if sess.Values["invite"] == true {
		controller.LoginInvite(w, r)
		return
	}

	if sess.Values["loggedin"] == true {
		next()
		return
	}

	controller.LoginAuth(w, r)
}

func None(w http.ResponseWriter, r *http.Request, next webserver.Next) {
	next()
}
