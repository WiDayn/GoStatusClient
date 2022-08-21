package utils

import (
	"GoStatusClient/logger"
	"github.com/go-ping/ping"
	"time"
)

func Ping(addr string) (int64, int, int) {
	pinger, err := ping.NewPinger(addr)
	pinger.Timeout = time.Second * 2
	pinger.SetPrivileged(true)
	pinger.Count = 1
	err = pinger.Run()
	if err != nil {
		logger.Error("PINGER ERROR", err)
	}
	stats := pinger.Statistics()
	return stats.AvgRtt.Milliseconds(), stats.PacketsRecv, stats.PacketsSent
}
