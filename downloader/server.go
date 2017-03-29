package main

import (
	"log"

	"github.com/itzaname/anime-streaming/api"
)

const DownloadWorkers = 3
const DiskSpace = 1.45e11

func loopGetDownloadQueue() []api.Download {
	for {
		queue, err := api.GetDownloadQueue()
		if err != nil {
			log.Println("GetDownloadQueue:", err)
			continue
		}

		return queue
	}
}

func main() {
	log.Println("Getting download queue")

	queue := loopGetDownloadQueue()

	log.Println("Starting download")

	log.Println(startDownload(queue))

	log.Println("Downloads complete")
}