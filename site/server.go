package main

import (
	"encoding/json"

	"github.com/itzaname/anime-streaming/site/app/lib/config"
	"github.com/itzaname/anime-streaming/site/app/lib/database"
	"github.com/itzaname/anime-streaming/site/app/lib/log"
	"github.com/itzaname/anime-streaming/site/app/lib/session"
	"github.com/itzaname/anime-streaming/site/app/lib/view"
	"github.com/itzaname/anime-streaming/site/app/lib/webserver"
	"github.com/itzaname/anime-streaming/site/app/route"
	"github.com/itzaname/anime-streaming/site/app/routines"
)

func main() {
	log.Set(log.LevelDebug)

	log.Infof("Loading config file `config.json`")
	err := config.Load("config.json", cfg)
	if err != nil {
		log.Fatalln(err)
	}

	view.Configure(cfg.View)
	view.LoadTemplates(cfg.Template.Root, cfg.Template.Children)

	session.Configure(cfg.Session)

	if err := database.Init(cfg.Database); err != nil {
		log.Fatalln(err)
	}

	//database.DB.FlushDb()

	/*download, err := model.DownloadAdd(27)
	if err != nil {
		panic(err)
	}

	download.Nyaa = "69"
	download.Anime = 27
	download.Update()*/

	routines.SetupRoutines()

	//log.Debugln(model.InviteCreate())

	server := webserver.New(cfg.Webserver)

	route.Init(server)

	log.Infoln("Server is up.")
	server.Start()
}

var cfg = &configuration{}

type configuration struct {
	Webserver webserver.Info  `json:"Webserver"`
	Database  database.Info   `json:"Database"`
	View      view.View       `json:"View"`
	Template  view.Template   `json:"Template"`
	Session   session.Session `json:"Session"`
}

// ParseJSON unmarshals bytes to structs
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}
