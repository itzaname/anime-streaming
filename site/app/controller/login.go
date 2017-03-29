package controller

import (
	"github.com/itzaname/anime-streaming/site/app/lib/log"
	"github.com/itzaname/anime-streaming/site/app/lib/myanimelist"
	"github.com/itzaname/anime-streaming/site/app/lib/session"
	"github.com/itzaname/anime-streaming/site/app/lib/view"
	"github.com/itzaname/anime-streaming/site/app/model"
	"net/http"
)

func LoginAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if r.ParseForm() != nil {
			v := view.New(r)
			v.Name = "auth"
			v.Vars["error"] = "An unexpected error occured."
			v.RenderSingle(w)
			return
		}

		if r.Form.Get("username") == "" || r.Form.Get("password") == "" {
			v := view.New(r)
			v.Name = "auth"
			v.Vars["error"] = "Invalid username or password."
			v.RenderSingle(w)
			return
		}

		// Test with MAL
		malclient := myanimelist.NewClient(r.Form.Get("username"), r.Form.Get("password"))
		_, err := malclient.Search("test")
		if err != nil {
			v := view.New(r)
			v.Name = "auth"
			v.Vars["error"] = "Invalid username or password."
			v.RenderSingle(w)
			return
		}

		sess := session.Instance(r)
		sess.Values["invite"] = true
		sess.Values["loggedin"] = true
		sess.Values["maluser"] = r.Form.Get("username")
		sess.Values["malpass"] = r.Form.Get("password")

		exists, err := model.UserExists(sess.Values["maluser"].(string))
		if err != nil {
			v := view.New(r)
			v.Name = "auth"
			v.Vars["error"] = "An unexpected error occured."
			v.RenderSingle(w)
			return
		}

		if exists {
			sess.Values["invite"] = false
			log.Debugln("Returning user logged in:", sess.Values["maluser"])
		}

		sess.Save(r, w)

		http.Redirect(w, r, "/", 302)
		return
	}

	v := view.New(r)
	v.Name = "auth"
	v.RenderSingle(w)
}

func LoginInvite(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if r.ParseForm() != nil {
			v := view.New(r)
			v.Name = "invite"
			v.Vars["error"] = "An unexpected error occured."
			v.RenderSingle(w)
			return
		}

		if r.Form.Get("invite") == "" {
			v := view.New(r)
			v.Name = "invite"
			v.Vars["error"] = "Invalid invite code."
			v.RenderSingle(w)
			return
		}

		valid, err := model.InviteValid(r.Form.Get("invite"))
		if err != nil {
			v := view.New(r)
			v.Name = "invite"
			v.Vars["error"] = "An unexpected error occured."
			v.RenderSingle(w)
			return
		}

		if !valid {
			v := view.New(r)
			v.Name = "invite"
			v.Vars["error"] = "Invalid invite code."
			v.RenderSingle(w)
			return
		}

		model.InviteDelete(r.Form.Get("invite"))

		sess := session.Instance(r)

		if model.UserAdd(sess.Values["maluser"].(string), false) != nil {
			v := view.New(r)
			v.Name = "invite"
			v.Vars["error"] = "An unexpected error occured."
			v.RenderSingle(w)
			return
		}

		sess.Values["invite"] = false
		sess.Save(r, w)

		log.Debugln("New user registered:", sess.Values["maluser"].(string))

		http.Redirect(w, r, "/", 302)
		return
	}

	v := view.New(r)
	v.Name = "invite"
	v.RenderSingle(w)
}
