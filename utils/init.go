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

type UserAgentInfo struct {
	IP          string `json:"ip"`
	Country     string `json:"country"`
	CountryCode string `json:"short_name"`
}

func Init() {
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
	system.StaticBasicInformation.OS = n.Platform + " " + n.PlatformFamily + " " + n.PlatformVersion
	system.StaticBasicInformation.Hostname = n.Hostname

	system.StaticBasicInformation.IP = strings.Split(system.StaticBasicInformation.IP, "\n")[0]
	if content, err := net.GetRequest("http://ip-api.com/json/?lang=zh-CN"); err != nil {
		logger.Warning("Requesting IP information error", err)
		logger.Info("Trying other API...", nil)
		if content, err := net.GetRequest("https://ip.useragentinfo.com/json"); err != nil {
			logger.Warning("Requesting IP information error", err)
		} else {
			var userAgentInfo UserAgentInfo
			if err := json.Unmarshal([]byte(content), &userAgentInfo); err != nil {
				logger.Warning("Reading IP information error", err)
			}
			system.StaticBasicInformation.Country = userAgentInfo.Country
			system.StaticBasicInformation.CountryCode = userAgentInfo.CountryCode
		}
	} else {
		if err := json.Unmarshal([]byte(content), &system.StaticBasicInformation); err != nil {
			logger.Warning("Reading IP information error", err)
		}
	}
	logger.Success("Load basic information success", nil)
}
