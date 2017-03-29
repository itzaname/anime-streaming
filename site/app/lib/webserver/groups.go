package webserver

import "net/http"

// Groups struct containing group info
type Groups struct {
	webserver  *WebServer
	middleware []Middleware
}

// Add adds middleware to be used
func (group *Groups) Add(mw Middleware) {
	group.middleware = append(group.middleware, mw)
}

func (group *Groups) iterator(index int, w http.ResponseWriter, r *http.Request, handler http.HandlerFunc) func() {
	return func() {
		if index > len(group.middleware)-1 {
			handler(w, r)
		} else {
			group.middleware[index](w, r, group.iterator(index+1, w, r, handler))
		}
	}
}

func (group *Groups) handler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(group.middleware) > 0 {
			group.middleware[0](w, r, group.iterator(1, w, r, handler))
		}
	}
}

// Get adds new get handler for route
func (group *Groups) Get(path string, handler http.HandlerFunc) {
	group.webserver.router.Get(path, group.handler(handler))
}

// Post adds new post handler for route
func (group *Groups) Post(path string, handler http.HandlerFunc) {
	group.webserver.router.Post(path, group.handler(handler))
}
