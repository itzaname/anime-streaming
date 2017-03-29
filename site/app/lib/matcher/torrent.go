package matcher

import (
	"sort"

	"strings"

	"fmt"

	"github.com/itzaname/anime-streaming/site/app/lib/nyaa"
	"github.com/xrash/smetrics"
	"regexp"
)

type nyaaSeedSorter []nyaa.NyaaItem

func (a nyaaSeedSorter) Len() int           { return len(a) }
func (a nyaaSeedSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a nyaaSeedSorter) Less(i, j int) bool { return a[i].Seeders > a[j].Seeders }

type scoredNyaa struct {
	Torrent nyaa.NyaaItem
	Score   int
}

type scoredNyaaSorter []scoredNyaa

func (a scoredNyaaSorter) Len() int           { return len(a) }
func (a scoredNyaaSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a scoredNyaaSorter) Less(i, j int) bool { return a[i].Score < a[j].Score }

func performMultiSearchSorted(titles []string) ([]nyaa.NyaaItem, error) {
	torrents := map[string]nyaa.NyaaItem{}

	for _, title := range titles {
		results, err := nyaa.PerformSearch(title)
		if err != nil {
			return nil, err
		}

		for _, torrent := range results {
			if torrent.Seeders >= 3 {
				if _, ok := torrents[torrent.Id]; ok {
					existing := smetrics.WagnerFischer(strings.ToLower(torrents[torrent.Id].Search), hyperCleanTitle(torrents[torrent.Id].Title), 1, 1, 2)
					new := smetrics.WagnerFischer(strings.ToLower(torrent.Search), hyperCleanTitle(torrent.Title), 1, 1, 2)
					if new < existing {
						torrents[torrent.Id] = torrent
					}
				} else {
					torrents[torrent.Id] = torrent
				}
			}
		}
	}

	sorted := nyaaSeedSorter{}

	for _, torrent := range torrents {
		sorted = append(sorted, torrent)
	}

	sort.Sort(sorted)

	return sorted, nil
}

func isValidSearchItem(search string, title string, episode int) bool {
	title = strings.ToLower(title)
	search = strings.ToLower(search)
	title = strings.Replace(title, search, "", -1)
	for _,s := range regexp.MustCompile(`\d+`).FindAllString(title, -1) {
		if s[0] == '0' {
			return s[1:] == fmt.Sprintf("%d", episode)
		}
		return s == fmt.Sprintf("%d", episode)
	}
	return false//strings.Contains(title, fmt.Sprintf("%d", episode))
}

func isValidTorrentEpisodes(info nyaa.NyaaItem, episodes int) (bool, error) {
	meta, err := nyaa.TorrentInfo(info.Id)
	if err != nil {
		return false, err
	}

	list, err := MatchEpisodes(info.Search, episodes, meta.Info.Files)
	if err != nil || len(list) < episodes {
		return false, nil
	}

	return true, nil
}

func isValidTorrentSingleEpisode(info nyaa.NyaaItem) (bool, error) {
	meta, err := nyaa.TorrentInfo(info.Id)
	if err != nil {
		return false, err
	}

	if len(meta.Info.Files) > 0 {
		result := MatchMovie(info.Search, meta.Info.Files)
		if result == "" {
			return false, nil
		}
	}
	return true, nil
}

func checkScoreEpisodesRange(score int, items []scoredNyaa, episodes int) []scoredNyaa {
	buffer := []scoredNyaa{}
	for _, item := range items {
		if item.Score == score {
			valid, err := isValidTorrentEpisodes(item.Torrent, episodes)
			if err == nil && valid {
				buffer = append(buffer, item)
			}
		}
	}

	return buffer
}

func checkScoreEpisodeRange(score int, items []scoredNyaa) []scoredNyaa {
	buffer := []scoredNyaa{}
	for _, item := range items {
		if item.Score == score {
			valid, err := isValidTorrentSingleEpisode(item.Torrent)
			if err == nil && valid {
				buffer = append(buffer, item)
			}
		}
	}

	return buffer
}

func getMaxScore(items []scoredNyaa) int {
	max := 0
	for _, item := range items {
		if item.Score > max {
			max = item.Score
		}
	}
	return max
}

func getHighestSeedRes(items []scoredNyaa, resolution string) (string, string) {
	bestseed := -1
	bestitem := scoredNyaa{}
	for _, item := range items {
		if item.Torrent.Seeders > bestseed && extractResolution(item.Torrent.Title) == resolution {
			bestseed = item.Torrent.Seeders
			bestitem = item
		}
	}

	if bestseed > 0 {
		return bestitem.Torrent.Id, bestitem.Torrent.Search
	}

	return "", ""
}

func getHighestSeed(items []scoredNyaa) (string, string) {
	bestseed := -1
	bestitem := scoredNyaa{}
	for _, item := range items {
		if item.Torrent.Seeders > bestseed {
			bestseed = item.Torrent.Seeders
			bestitem = item
		}
	}

	if bestseed > 0 {
		return bestitem.Torrent.Id, bestitem.Torrent.Search
	}

	return "", ""
}

// MatchEpisodeTorrent finds torrent on nyaa with info and returns NyaaID, search, error
func MatchEpisodeTorrent(titles []string, episodes int) (string, string, error) {
	results, err := performMultiSearchSorted(titles)
	if err != nil {
		return "", "", err
	}

	scored := scoredNyaaSorter{}

	for _, item := range results {
		scored = append(scored, scoredNyaa{item, smetrics.WagnerFischer(strings.ToLower(item.Search), hyperCleanTitle(item.Title), 1, 1, 2)})
	}

	sort.Sort(scored)

	maxscore := getMaxScore(scored)
	for i := 0; i < maxscore+1; i++ {
		items := checkScoreEpisodesRange(i, scored, episodes)
		if len(items) > 0 {
			best, search := getHighestSeedRes(items, "1080")
			if best != "" {
				return best, search, nil
			}
			best, search = getHighestSeedRes(items, "720")
			if best != "" {
				return best, search, nil
			}
			best, search = getHighestSeedRes(items, "480")
			if best != "" {
				return best, search, nil
			}
			best, search = getHighestSeed(items)
			if best != "" {
				return best, search, nil
			}
			return "", "", fmt.Errorf("Something bad happened...")
		}
	}

	return "", "", fmt.Errorf("None found...")
}

func MatchSingleEpisodeTorrent(titles []string, episode int, prefix bool) (string, string, error) {
	if prefix {
		for i := 0; i < len(titles); i++ {
			if episode < 10 {
				titles[i] += fmt.Sprintf(" 0%d",episode)
			} else {
				titles[i] += fmt.Sprintf(" %d",episode)
			}
		}
	}

	results, err := performMultiSearchSorted(titles)
	if err != nil {
		return "", "", err
	}

	scored := scoredNyaaSorter{}

	for _, item := range results {
		if isValidSearchItem(item.Search, hyperCleanTitle(item.Title), episode) {
			scored = append(scored, scoredNyaa{item, smetrics.WagnerFischer(strings.ToLower(item.Search), hyperCleanTitle(item.Title), 1, 1, 2)})
		}
	}

	sort.Sort(scored)

	maxscore := getMaxScore(scored)
	for i := 0; i < maxscore+1; i++ {
		items := checkScoreEpisodeRange(i, scored)
		if len(items) > 0 {
			best, search := getHighestSeedRes(items, "1080")
			if best != "" {
				return best, search, nil
			}
			best, search = getHighestSeedRes(items, "720")
			if best != "" {
				return best, search, nil
			}
			best, search = getHighestSeedRes(items, "480")
			if best != "" {
				return best, search, nil
			}
			best, search = getHighestSeed(items)
			if best != "" {
				return best, search, nil
			}
			return "", "", fmt.Errorf("Something bad happened...")
		}
	}

	return "", "", fmt.Errorf("None found...")
}
