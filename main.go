package main

import (
	"github.com/0xJacky/Homework-api/live"
	"github.com/0xJacky/Homework-api/pkg"
	"github.com/0xJacky/Homework-api/routers"
	"github.com/0xJacky/Homework-api/settings"
	"github.com/gin-gonic/gin"

	"log"
)

func main() {
	gin.SetMode(settings.ServerSettings.RunMode)
	r := routers.InitRouter()

	pkg.ExistsOrCreate("upload")
	pkg.ExistsOrCreate("upload/avatar")

	go live.SyncLastActive()

	err := r.Run(":" + settings.ServerSettings.HttpPort)

	if err != nil {
		log.Fatal(err)
	}
}
