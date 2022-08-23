package utils

import (
	"GoStatusClient/config"
	"GoStatusClient/logger"
	"GoStatusClient/system"
	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"net/url"
	"strconv"
	"time"
)

type UpdateRequest struct {
	ClientId           string
	DisplayName        string
	CountryCode        string
	DynamicInformation system.DynamicInformation
	UpdateTime         time.Time
	Online             bool
}

func Update() {
	for {
		Send()
		time.Sleep(time.Second * 2)
	}
}

func GetDynamicInformation() {
	cc, _ := cpu.Percent(time.Millisecond*300, false)
	v, _ := mem.VirtualMemory()
	d, _ := disk.Usage("/")

	IOCountersStat, _ := net.IOCounters(true)
	var tempDownStreamDataSize uint64
	var tempUpStreamDataSize uint64

	system.StaticDynamicInformation.CPUAvg = cc[0]
	system.StaticDynamicInformation.MemAll = v.Total
	system.StaticDynamicInformation.MenUsed = v.Used
	system.StaticDynamicInformation.MenFree = v.Total - v.Used
	system.StaticDynamicInformation.MemUsedPercent = float64(v.Used) / float64(v.Total) * 100.0
	system.StaticDynamicInformation.DiskInformation.Total = d.Total / 1024 / 1024 / 1024
	system.StaticDynamicInformation.DiskInformation.Used = d.Used / 1024 / 1024 / 1024
	system.StaticDynamicInformation.DiskInformation.Percent = uint64(d.UsedPercent)

	for i := 0; i < len(IOCountersStat); i++ {
		tempDownStreamDataSize += IOCountersStat[i].BytesRecv
		tempUpStreamDataSize += IOCountersStat[i].BytesSent
	}

	system.StaticDynamicInformation.NowDownStreamDataSize = int((tempDownStreamDataSize - system.StaticDynamicInformation.TotalDownStreamDataSize) / 1024)
	system.StaticDynamicInformation.NowUpStreamDataSize = int((tempUpStreamDataSize - system.StaticDynamicInformation.TotalUpStreamDataSize) / 1024)
	system.StaticDynamicInformation.TotalDownStreamDataSize = tempDownStreamDataSize
	system.StaticDynamicInformation.TotalUpStreamDataSize = tempUpStreamDataSize
	system.StaticDynamicInformation.CT.AvgRTT, system.StaticDynamicInformation.CT.PacketsReceive, system.StaticDynamicInformation.CT.PacketsSent = Ping(config.Config.CT.Addr)
	system.StaticDynamicInformation.CU.AvgRTT, system.StaticDynamicInformation.CU.PacketsReceive, system.StaticDynamicInformation.CU.PacketsSent = Ping(config.Config.CU.Addr)
	system.StaticDynamicInformation.CM.AvgRTT, system.StaticDynamicInformation.CM.PacketsReceive, system.StaticDynamicInformation.CM.PacketsSent = Ping(config.Config.CM.Addr)
}

func Send() {
	u := url.URL{Scheme: "ws", Host: config.Config.Server.IP + ":" + strconv.Itoa(config.Config.Server.Port), Path: "/update"}
	logger.Info("Connecting to "+u.String(), nil)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logger.Error("Connect to server fail", err)
		return
	}
	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
			logger.Warning("Close websocket error", err)
			return
		}
	}(c)

	err = c.WriteJSON(config.Config.SecretKey)
	if err != nil {
		logger.Error("Websocket writing json error", err)
		return
	}

	for true {
		GetDynamicInformation()
		updateRequest := UpdateRequest{
			ClientId:           config.Config.ClientId,
			DisplayName:        config.Config.DisplayName,
			CountryCode:        system.StaticBasicInformation.CountryCode,
			DynamicInformation: system.StaticDynamicInformation,
			UpdateTime:         time.Now(),
		}
		err := c.WriteJSON(updateRequest)
		if err != nil {
			logger.Error("Websocket writing json error", err)
			err := c.Close()
			if err != nil {
				logger.Error("Websocket close error", err)
				return
			}
			return
		}
		time.Sleep(time.Second * 1)
	}
}
