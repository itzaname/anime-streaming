package nyaa

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type listItem struct {
	Title   string
	Url     string
	Seeders int
}

func performSearch(query string) (io.ReadCloser, error) {
	urlString := fmt.Sprintf("https://www.nyaa.se/?page=search&cats=1_37&filter=0&sort=2&term=%s", url.QueryEscape(query))
	resp, err := http.Get(urlString)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error %d", resp.Status)
	}

	return resp.Body, nil
}

func readToken(list *[]listItem, z *html.Tokenizer, tt html.TokenType) {
	t := z.Token()
	url := ""
	title := ""
	if t.Data == "td" {
		for _, td := range t.Attr {
			if td.Key == "class" && strings.Contains(td.Val, "tlistname") {
				z.Next()
				t = z.Token()
				for _, a := range t.Attr {
					if a.Key == "href" {
						url = a.Val
					}
				}
				z.Next()
				t = z.Token()
				title = t.String()
				for {
					tt = z.Next()
					t = z.Token()
					if tt == html.ErrorToken {
						return
					}
					if t.Data == "td" {
						for _, td := range t.Attr {
							if td.Key == "class" && strings.Contains(td.Val, "tlistsn") {
								z.Next()
								t = z.Token()
								seeders, err := strconv.ParseInt(t.String(), 10, 64)
								if err != nil {
									return
								}
								*list = append(*list, listItem{title, url, int(seeders)})
								return
							}
						}
					}
				}
			}
		}
	}
}

func getNyaaData(query string) ([]listItem, error) {
	var listData []listItem

	kek, err := performSearch(query)
	if err != nil {
		return listData, err
	}

	defer kek.Close()

	z := html.NewTokenizer(kek)
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			return listData, nil
		case html.StartTagToken:
			readToken(&listData, z, tt)
			break
		}
	}
}
