package routines

import (
	"strings"

	"github.com/itzaname/anime-streaming/site/app/lib/log"
	"github.com/itzaname/anime-streaming/site/app/lib/matcher"
	"github.com/itzaname/anime-streaming/site/app/model"
	"github.com/itzaname/anime-downloader/nyaa"
)

// 1 Airing 2 Aired 3 Not Aired

func getMALTitles(anime model.Anime) []string {
	output := []string{strings.TrimSpace(strings.Replace(anime.Info.Title, "(TV)", "", -1))}
	synonyms := strings.Split(anime.Info.Synonyms, ";")
	for i := 0; i < len(synonyms); i++ {
		if len(strings.TrimSpace(synonyms[i])) > 2 {
			output = append(output, strings.TrimSpace(synonyms[i]))
		}
	}
	return output
}

func SetupDownloads() {
	animelist, err := model.AnimeList()
	if err != nil {
		log.Errorln("SetupDownloads failed.", err)
		return
	}

	for _, anime := range animelist {
		if !anime.Complete && anime.Info.Status != 3 && /* REMOVE THIS PLEASE */ (anime.ID == 30831 || anime.ID == 33988 ) { // If anime isnt completely downloaded and if it has aired
			if anime.State == 1 || anime.Downloaded > 0  || anime.ID == 33988 {
				if anime.Downloaded >= anime.Info.Episodes && anime.Info.Episodes != 0 {
					break
				}

				// We will download one by one
				episode := 1
				for {
					if episode > anime.Info.Episodes {
						break
					}

					id, search, err := matcher.MatchSingleEpisodeTorrent(getMALTitles(anime), episode, true)
					if err != nil {
						log.Errorln(anime.ID," SetupDownloads - MatchSingleEpisodeTorrent:", err)
						break
					}

					meta, err := nyaa.TorrentInfo(id)
					if err != nil {
						log.Errorln(anime.ID," SetupDownloads:", err)
						continue
					}

					download, err := model.DownloadAdd(anime.ID)
					if err != nil {
						log.Errorln(anime.ID," SetupDownloads:", err)
						continue
					}

					download.Anime = anime.ID
					download.Count = 1
					download.Name = anime.Info.Title
					download.Search = search
					download.Nyaa = id
					download.Magnet = meta.MagnetLink()
					download.Update()

					anime.Downloaded++ // Single download add one
					anime.Update()

					log.Infoln("New download queued!", download.ID, search)

					episode++
				}
			} else {
				// We will find batch downloads but it may be risky
				id, search, err := matcher.MatchEpisodeTorrent(getMALTitles(anime), anime.Info.Episodes)
				if err != nil {
					log.Errorln(anime.ID," SetupDownloads - MatchEpisodeTorrent:", err)
					continue
				}

				meta, err := nyaa.TorrentInfo(id)
				if err != nil {
					log.Errorln(anime.ID," SetupDownloads:", err)
					continue
				}

				download, err := model.DownloadAdd(anime.ID)
				if err != nil {
					log.Errorln(anime.ID," SetupDownloads:", err)
					continue
				}

				download.Anime = anime.ID
				download.Count = anime.Info.Episodes
				download.Name = anime.Info.Title
				download.Search = search
				download.Nyaa = id
				download.Magnet = meta.MagnetLink()
				download.Update()

				anime.Downloaded = download.Count // Set downloaded exact as this is the ONLY download that should take place
				anime.Update()

				log.Infoln("New download queued!", download.ID, search)

			}
		}
	}
}
