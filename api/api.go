package api

import (
	"net/url"
	"net/http"
	"fmt"
	"bytes"
	"strconv"
)

const endpoint = "http://127.0.0.1:8080"

func callAPIPost(resource string, data url.Values) (*http.Response, error){
	u, _ := url.ParseRequestURI(endpoint)
	u.Path = resource
	urlStr := fmt.Sprintf("%v", u) // "https://api.com/user/"

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode())) // <-- URL-encoded payload
	r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	return client.Do(r)
}

func callAPIGet(resource string) (*http.Response, error){
	u, _ := url.ParseRequestURI(endpoint)
	u.Path = resource
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	r, _ := http.NewRequest("GET", urlStr, nil)
	r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")

	return client.Do(r)
}

