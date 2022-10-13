package main

import (
	"chat/cache"
	"chat/conf"
	"chat/router"
	"chat/service"
)

func main() {
	conf.Init()
	cache.Init()
	go service.Manager.Start()
	r := router.NewRouter()
	_ = r.Run(conf.HttpPort)
}
