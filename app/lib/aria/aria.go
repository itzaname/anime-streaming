package aria

import (
	"anime-multi-stream/src/lib/log"
	"os/exec"
	"path/filepath"
	"strconv"

	"anime-multi-stream/src/lib/config"
	"github.com/NoahShen/aria2rpc"
)

type Download struct {
	Status   string
	Download uint64
	Length   uint64
	Speed    float64
}

func StartRPCServer() {
	c := exec.Command(filepath.FromSlash(config.Cfg.TorrentBinary), "--enable-rpc=true")
	_, err := c.Output()
	if err != nil {
		log.Fatalf("RPC Start failed!")
	}
}

func StartTorrentDownload(file string, outPath string) (string, error) {
	options := map[string]interface{}{}
	options["dir"] = outPath

	options["seed-time"] = "0"
	options["file-allocation"] = "none"
	options["listen-port"] = "8767"

	return aria2rpc.AddTorrent(file, options)
}

func GetDownloadInfo(gid string) (Download, error) {
	var dl Download

	status, err := aria2rpc.GetStatus(gid, []string{"status", "totalLength", "completedLength", "downloadSpeed"})
	if err != nil {
		return dl, err
	}

	comp, err := strconv.ParseUint(status["completedLength"].(string), 10, 64)
	if err != nil {
		return dl, err
	}

	total, err := strconv.ParseUint(status["totalLength"].(string), 10, 64)
	if err != nil {
		return dl, err
	}

	speed, err := strconv.ParseFloat(status["downloadSpeed"].(string), 64)
	if err != nil {
		return dl, err
	}

	dl.Download = comp
	dl.Length = total
	dl.Speed = speed
	dl.Status = status["status"].(string)

	return dl, nil
}
