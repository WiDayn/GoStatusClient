package utils

import (
	"GoStatusClient/config"
	"GoStatusClient/net"
	"GoStatusClient/system"
	"encoding/json"
	"strconv"
)

type RegisterRequest struct {
	ClientId         string
	DisplayName      string
	BasicInformation system.BasicInformation
}

func Register() {
	post, _ := json.Marshal(RegisterRequest{
		ClientId:         config.Config.ClientId,
		DisplayName:      config.Config.DisplayName,
		BasicInformation: system.StaticBasicInformation,
	})

	feedback, _ := net.PostRequest("http://"+config.Config.Server.IP+":"+strconv.Itoa(config.Config.Server.Port)+"/register?SecretKey="+config.Config.SecretKey, string(post))

	println(feedback)
}
