package utils

import (
	"GoStatusClient/logger"
	"GoStatusClient/net"
	"GoStatusClient/system"
	"encoding/json"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"strings"
)

func Init() {
	var err error

	c, _ := cpu.Info()
	n, _ := host.Info()

	for _, subCpu := range c {
		system.StaticBasicInformation.CPUs = append(system.StaticBasicInformation.CPUs, system.CPU{
			CPUName:  subCpu.ModelName,
			CPUCores: subCpu.Cores,
		})
	}
	system.StaticBasicInformation.CPUPhysicalCnt, _ = cpu.Counts(false)
	system.StaticBasicInformation.CPULogicalCnt, _ = cpu.Counts(true)
	system.StaticBasicInformation.OS = n.Platform + n.PlatformFamily + n.PlatformVersion
	system.StaticBasicInformation.Hostname = n.Hostname
	if system.StaticBasicInformation.IP, err = net.GetRequest("http://myip.ipip.net/s"); err != nil {
		logger.Warning("Requesting IP error:", err)
	}
	system.StaticBasicInformation.IP = strings.Split(system.StaticBasicInformation.IP, "\n")[0]
	if content, err := net.GetRequest("http://ip-api.com/json/" + system.StaticBasicInformation.IP + "?lang=zh-CN"); err != nil {
		logger.Warning("Requesting IP information error:", err)
	} else {
		if err := json.Unmarshal([]byte(content), &system.StaticBasicInformation); err != nil {
			logger.Warning("Reading IP information error:", err)
		}
	}
	logger.Success("Load basic information success", nil)
}
