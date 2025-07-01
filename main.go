package main

import (
	"context"
	"fmt"
	"go/src/github.com/rilgilang/gowebcapture/bootstrap"
	"go/src/github.com/rilgilang/gowebcapture/pkg"
	"go/src/github.com/rilgilang/gowebcapture/service"
)

func main() {

	//https://satumomen.com/preview/peresmian-rs
	//https://joinedwithshan.viding.co/
	//https://app.sangmempelai.id/pilihan-tema/sunda-01
	//https://adirara.webnikah.com/?templatecoba=156/kepada:Budi%20dan%20Ani-Bandung
	//https://ourmoment.my.id/art-6

	bootstrapClienter, config, err := bootstrap.Setup()
	if err != nil {
		panic(err)
	}

	cache := pkg.NewCache(bootstrapClienter.Cache)

	storage := pkg.NewStorage(bootstrapClienter.Storage)

	crawler := service.NewCrawler(storage, config)

	ctx := context.Background()

	for {
		urlLink, err := cache.BRpop(ctx, "video_queue")
		if err != nil {
			fmt.Println("err redis --> ", err)
			break
		}

		// Launch browser and interact
		err = crawler.RunBrowserAndInteract(ctx, urlLink)

		if err != nil {
			fmt.Println("err processing --> ", err)
		}
	}
}
