package main

import (
	"GoStatusClient/config"
	"GoStatusClient/utils"
)

func main() {
	utils.Welcome()
	config.Read()
	utils.Init()
	utils.Register()
	utils.Update()
}
