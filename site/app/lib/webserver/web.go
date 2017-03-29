package webserver

import (
	"log"
	"net/http"

	"github.com/go-zoo/bone"
)

// Next used for iterator in Groups
type Next func()

// Middleware used for defining middleware functions
type Middleware func(http.ResponseWriter, *http.Request, Next)

// WebServer struct containing http instance as well as handlers
type WebServer struct {
	listen string
	router *bone.Mux
}

// New Returns a new webserver object.
func New(info Info) *WebServer {
	router := bone.New()
	// Add more to this later on
	return &WebServer{info.Listen + ":" + info.Port, router}
}

// Group returns a new group for the webserver object
func (web *WebServer) Group() *Groups {
	group := Groups{}
	group.webserver = web
	return &group
}

// Start starts up the webserver
func (web *WebServer) Start() {
	log.Println(http.ListenAndServe(web.listen, web.router))
}

// Get adds new get handler for route
func (web *WebServer) Get(path string, handler http.HandlerFunc) {
	web.router.Get(path, handler)
}

// Post adds new post handler for route
func (web *WebServer) Post(path string, handler http.HandlerFunc) {
	web.router.Post(path, handler)
}

// NotFound sets the 404 handler
func (web *WebServer) NotFound(handler http.HandlerFunc) {
	web.router.NotFound(handler)
}

// Info config data
type Info struct {
	Listen string
	Port   string
}
