package api

import (
	"encoding/json"
	"io/ioutil"
)

type Download struct {
	Download int    `json:"download"`
	Episodes int    `json:"episodes"`
	Search   string `json:"search"`
	Title    string `json:"title"`
	Torrent  string `json:"torrent"`
	Magnet   string `json:"magnet"`
}

func GetDownloadQueue() ([]Download, error) {
	downloads := []Download{}

	resp, err := callAPIGet("/api/download/queue")
	if err != nil {
		return downloads, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return downloads, err
	}

	err = json.Unmarshal(data, &downloads)
	if err != nil {
		return downloads, err
	}

	return downloads, nil
}
