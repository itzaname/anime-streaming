package nyaa

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/cnt0/gotrntmetainfoparser"
)

type NyaaItem struct {
	Title   string
	Id      string
	Search  string
	Seeders int
}

// PerformSearch searches and parses nyaa returning the results
func PerformSearch(query string) ([]NyaaItem, error) {
	var items []NyaaItem

	list, err := getNyaaData(query)
	if err != nil {
		return items, err
	}

	for _, item := range list {
		parsed, err := url.Parse(item.Url)
		if err != nil {
			return items, err
		}

		id := parsed.Query().Get("tid")
		if id == "" {
			return items, fmt.Errorf("Invalid ID")
		}

		items = append(items, NyaaItem{item.Title, id, query, item.Seeders})
	}

	return items, nil
}

// TorrentInfo
func TorrentInfo(id string) (*gotrntmetainfoparser.MetaInfo, error) {
	resp, err := http.Get(fmt.Sprintf("https://www.nyaa.se/?page=download&tid=%s", id))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Non 200 reponse %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	meta, err := gotrntmetainfoparser.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return meta, nil
}
