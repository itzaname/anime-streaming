package main

import (
	"time"

	"github.com/anacrolix/torrent"
	"github.com/itzaname/anime-streaming/api"
	"github.com/itzaname/anime-streaming/site/app/lib/log"
	"github.com/itzaname/anime-downloader/nyaa"
	"io"
	"github.com/anacrolix/torrent/metainfo"
)

var client *torrent.Client

func downloadTorrent(data io.ReadCloser, worker int) error {
	defer data.Close()
	meta, err := metainfo.Load(data)
	if err != nil {
		return err
	}

	torrent, err := client.AddTorrent(meta)
	if err != nil {
		return err
	}

	log.Infoln(worker, "DSTART")

	<-torrent.GotInfo()
	torrent.DownloadAll()

	log.Infoln(worker, "DREADY")

	for {
		log.Infoln(worker, torrent.Info().Name, (float64(torrent.BytesCompleted())/float64(torrent.Length()))*100.0)
		time.Sleep(time.Second * 10)
		if torrent.BytesCompleted() >= torrent.Length() {
			return nil
		}
	}
}

func downloadWorker(data chan api.Download, signal chan bool, id int) {
	log.Infoln(id, "Worker started.")
	for {
		log.Infoln(id, "Start here")
		download, ok := <-data
		if !ok {
			break
		}

		log.Infoln(id, "Go here")

		data, err := nyaa.TorrentData(download.Torrent)
		if err != nil {
			// This normally shouldnt ever happen
			// If this does happen then the download will just get skipped
			// We will try again once the queue is downloaded
			log.Errorln(id, "Torrent info failed!", err)
			continue
		}

		err = downloadTorrent(data, id)
		log.Infoln(id, "Be here")
		if err != nil {
			// This normally shouldnt ever happen
			// If this does happen then the download will just get skipped
			// We will try again once the queue is downloaded
			log.Errorln(id, "Torrent download failed!", err)
			continue
		}

		log.Infoln(id, "COMPLETE", download.Title)
	}

	signal <- true
	log.Infoln(id, "Worker Exited.")
}

func startDownload(queue []api.Download) error {
	var err error
	// Create torrent client
	client, err = torrent.NewClient(&torrent.Config{
		DataDir: "downloads",
		Seed:    false,
	})

	if err != nil {
		return err
	}

	signal := make(chan bool)
	data := make(chan api.Download)

	for i := 0; i < DownloadWorkers; i++ {
		go downloadWorker(data, signal, i)
	}

	for _, dl := range queue {
		data <- dl
	}

	close(data)

	for i := 0; i < DownloadWorkers; i++ {
		<-signal
		log.Infoln("Worker stop signal recieved.")
	}

	client.Close()

	return nil
}
