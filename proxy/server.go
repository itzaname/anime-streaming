package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	// This is the server to proxy
	remote, err := url.Parse("http://media.oblivion.ws")
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	// wrap the reverse proxy in a custom handler below
	http.HandleFunc("/", handler(proxy))

	// serve on port 1234 forever
	log.Fatal(http.ListenAndServe("localhost:8020", nil))
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	// returns a function that satisfies the http.HandlerFunc interface
	return func(w http.ResponseWriter, r *http.Request) {
		// set the user name before performing the proxy
		r.Header.Set("REMOTE_USER", "Bob")
		r.Host = "media.oblivion.ws"
		r.URL.Path = "/watch/stream/4898/3"
		p.ServeHTTP(w, r)
	}
}
