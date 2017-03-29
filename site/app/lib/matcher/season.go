package matcher

import (
	"regexp"
	"strconv"
	"strings"
)

var rxs1 = regexp.MustCompile(`season\s*?(\d)`)
var rxs2 = regexp.MustCompile(`(\d)\w{2,3}\s*?season`)
var rxs3 = regexp.MustCompile(`s(\d)`)

// IsFirstSeason returns if a title is likely the first season and if not the number
func IsFirstSeason(title string) (bool, int) {
	title = strings.ToLower(title)
	results := rxs1.FindAllStringSubmatch(title, 1)
	if len(results) > 0 {
		season, err := strconv.ParseInt(results[0][1], 10, 64)
		if err != nil {
			panic(err) // Regex fucked up big time
		}
		return false, int(season)
	}

	results = rxs2.FindAllStringSubmatch(title, 1)
	if len(results) > 0 {
		season, err := strconv.ParseInt(results[0][1], 10, 64)
		if err != nil {
			panic(err) // Regex fucked up big time
		}
		return false, int(season)
	}

	results = rxs3.FindAllStringSubmatch(title, 1)
	if len(results) > 0 {
		season, err := strconv.ParseInt(results[0][1], 10, 64)
		if err != nil {
			panic(err) // Regex fucked up big time
		}
		return false, int(season)
	}

	if strings.Contains(strings.Replace(title, " ", "", -1), "secondseason") {
		return false, 2
	}
	if strings.Contains(strings.Replace(title, " ", "", -1), "thirdseason") {
		return false, 3
	}
	if strings.Contains(strings.Replace(title, " ", "", -1), "fourthseason") {
		return false, 4
	}

	if strings.Contains(strings.Replace(title, " ", "", -1), "seasontwo") {
		return false, 2
	}
	if strings.Contains(strings.Replace(title, " ", "", -1), "seasonthree") {
		return false, 3
	}
	if strings.Contains(strings.Replace(title, " ", "", -1), "seasonfour") {
		return false, 4
	}

	return true, 1
}
