package route

import (
	"github.com/itzaname/anime-streaming/site/app/controller"
	"github.com/itzaname/anime-streaming/site/app/lib/log"
	"github.com/itzaname/anime-streaming/site/app/lib/webserver"
	"github.com/itzaname/anime-streaming/site/app/route/middleware"
)

func Init(srv *webserver.WebServer) {
	// User views
	stream := srv.Group()
	stream.Add(middleware.RequireAuth)
	stream.Get("/", controller.StreamHome)
	stream.Post("/", controller.StreamHome)
	stream.Get("/search", controller.StreamSearch)
	stream.Get("/watch/:id/:ep", controller.WatchVideo)

	// Public Api Stuff
	pubapi := srv.Group()
	pubapi.Add(middleware.RequireAuth)
	pubapi.Get("/api/anime/list", controller.APIAnimeList)

	// Server Api Stuff
	srvapi := srv.Group()
	srvapi.Add(middleware.None)
	srvapi.Get("/api/download/queue", controller.ServerAPIDownloadList)

	srv.Get("/invite", controller.DebugKey)

	// Serve static files
	srv.Get("/public/*filepath", controller.ServeStatic)

	log.Debugln("Routes initialized.")
}
