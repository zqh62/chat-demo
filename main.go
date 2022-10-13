package main

import (
	"chat/cache"
	"chat/conf"
	"chat/router"
)

func main() {
	conf.Init()
	cache.Init()
	r := router.NewRouter()
	_ = r.Run(conf.HttpPort)
}
