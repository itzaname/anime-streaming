package matcher

import (
	"regexp"
	"strings"
)

var rxcln = regexp.MustCompile(`[\[\({].+?[}\)\]]`)
var rxcln2 = regexp.MustCompile(`\d+-\d+`)
var rxrz1 = regexp.MustCompile(`(\d+?)p`)
var rxrz2 = regexp.MustCompile(`\d+x(\d+)`)

func cleanTitle(title string) string {
	match := rxcln.FindAllString(title, -1)
	for _, m := range match {
		title = strings.Replace(title, m, "", 1)
	}

	match = rxcln2.FindAllString(title, -1)
	for _, m := range match {
		title = strings.Replace(title, m, "", 1)
	}

	title = strings.Replace(title, "~", " ", -1)
	title = strings.Replace(title, ".", " ", -1)
	title = strings.Replace(title, "_", " ", -1)

	return strings.TrimSpace(title)
}

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func hyperCleanTitle(title string) string {
	title = strings.ToLower(cleanTitle(title))

	badStrings := []string{
		"mkv",
		"mp4",
		"bdrip",
		"+ specials/ovas",
		"+specials/ovas",
		"+ specials/ova",
		"+specials/ova",
		"+ ovas/specials",
		"+ovas/specials",
		"+ ova",
		"+ova",
		"+ specials",
		"+specials",
	}

	results := rxrz1.FindAllStringSubmatch(title, -1)
	for _, match := range results {
		title = strings.Replace(title, match[0], "", -1)
	}

	results = rxrz2.FindAllStringSubmatch(title, -1)
	for _, match := range results {
		title = strings.Replace(title, match[0], "", -1)
	}

	for _, bad := range badStrings {
		title = strings.Replace(title, bad, "", -1)
	}

	return standardizeSpaces(strings.TrimSpace(title))
}

func extractResolution(title string) string {
	results := rxrz1.FindAllStringSubmatch(title, -1)
	if len(results) > 0 {
		return results[0][1]
	}

	results = rxrz2.FindAllStringSubmatch(title, -1)
	if len(results) > 0 {
		return results[0][1]
	}

	return ""
}
