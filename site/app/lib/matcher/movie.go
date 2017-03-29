package matcher

import (
	"github.com/cnt0/gotrntmetainfoparser"
	"github.com/xrash/smetrics"
	"path/filepath"
	"strings"
)

func MatchMovie(title string, files []gotrntmetainfoparser.FileDict) string {
	episodes := []string{}
	for _, file := range files {
		if filepath.Ext(file.Path[len(file.Path)-1]) == ".mkv" || filepath.Ext(file.Path[len(file.Path)-1]) == ".mp4" || filepath.Ext(file.Path[len(file.Path)-1]) == ".avi" {
			episodes = append(episodes, file.Path[len(file.Path)-1])
		}
	}

	type ScoredString struct {
		Score  int
		String string
	}

	scored := []ScoredString{}
	for _, ep := range episodes {
		temp := cleanTitle(ep)

		temp = strings.Replace(temp, "mkv", "", -1)
		temp = strings.Replace(temp, "mp4", "", -1)
		temp = strings.Replace(temp, "avi", "", -1)
		temp = strings.TrimSpace(temp)

		scored = append(scored, ScoredString{smetrics.WagnerFischer(title, temp, 1, 1, 2), ep})
	}

	bestscore := -1
	bestitem := ScoredString{}
	for _, score := range scored {
		if bestscore == -1 || score.Score < bestscore {
			bestscore = score.Score
			bestitem = score
		}
	}

	return bestitem.String
}
