package matcher

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/cnt0/gotrntmetainfoparser"
	"github.com/xrash/smetrics"
)

var rxnum = regexp.MustCompile(`\d+`)

func findStringIndex(match string, list []string) (int, error) {
	for k, s := range list {
		if s == match {
			return k, nil
		}
	}
	return 0, fmt.Errorf("Not found")
}

func isInIntArray(match int, list []int) bool {
	for _, i := range list {
		if i == match {
			return true
		}
	}
	return false
}

func isInceremeting(in []int) bool {
	lowest := -1
	for _, i := range in {
		if i < lowest || lowest == -1 {
			lowest = i
		}
	}
	for i := 0; i < len(in); i++ {
		if !isInIntArray(lowest+i, in) {
			return false
		}
	}
	return true
}

func groupNumberStrings(list []string, required int) [][]string {
	matches := map[int][]string{}
	for _, s := range list {
		index := rxnum.FindAllStringIndex(s, -1)
		for _, i := range index {
			matches[i[0]] = append(matches[i[0]], s)
		}
	}

	out := [][]string{}
	for _, v := range matches {
		if len(v) >= required {
			out = append(out, v)
		}
	}

	return out
}

// We will scan for all the numbers and then check to see if they are incrementing
func isStringGroupIncrementing(list []string) bool {
	intlist := map[int][]int{}
	for _, s := range list {
		matches := rxnum.FindAllString(s, -1)
		for i := 0; i < len(matches); i++ {
			item, err := strconv.ParseInt(matches[i], 10, 64)
			if err != nil {
				// Bad things
				panic(err)
			}

			intlist[i] = append(intlist[i], int(item))
		}
	}

	for _, v := range intlist {
		if isInceremeting(v) {
			return true
		}
	}

	return false
}

func MatchEpisodes(title string, total int, files []gotrntmetainfoparser.FileDict) ([]string, error) {
	episodes := []string{}
	for _, file := range files {
		if filepath.Ext(file.Path[len(file.Path)-1]) == ".mkv" || filepath.Ext(file.Path[len(file.Path)-1]) == ".mp4" || filepath.Ext(file.Path[len(file.Path)-1]) == ".avi" {
			episodes = append(episodes, file.Path[len(file.Path)-1])
		}
	}

	buffer := []string{}
	for _, ep := range episodes {
		temp := cleanTitle(ep)

		temp = strings.Replace(temp, "mkv", "", -1)
		temp = strings.Replace(temp, "mp4", "", -1)
		temp = strings.Replace(temp, "avi", "", -1)
		temp = strings.TrimSpace(temp)

		buffer = append(buffer, temp)
	}

	matches := groupNumberStrings(buffer, total)

	type ScoredStringList struct {
		Score   int
		Strings []string
	}

	scores := []ScoredStringList{}
	for _, v := range matches {
		total := 0
		if isStringGroupIncrementing(v) {
			output := []string{}
			for _, s := range v {
				total += smetrics.WagnerFischer(title, s, 1, 1, 2)
				index, err := findStringIndex(s, buffer)
				if err != nil {
					return []string{}, err
				}
				output = append(output, episodes[index])
			}
			scores = append(scores, ScoredStringList{total, output})
		}
	}

	if len(scores) > 0 {
		lowest := -1
		var best ScoredStringList
		for _, score := range scores {
			if lowest == -1 || score.Score < lowest {
				best = score
				lowest = score.Score
			}
		}

		sort.Strings(best.Strings)

		return best.Strings, nil
	}

	return []string{}, fmt.Errorf("No episodes found.")
}
