package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go/src/github.com/rilgilang/gowebcapture/bootstrap"
	"go/src/github.com/rilgilang/gowebcapture/entities"
	"go/src/github.com/rilgilang/gowebcapture/pkg"
	"go/src/github.com/rilgilang/gowebcapture/repositories"
	"go/src/github.com/rilgilang/gowebcapture/service"
)

func main() {

	bootstrapClienter, config, err := bootstrap.Setup()
	if err != nil {
		panic(err)
	}

	cache := pkg.NewCache(bootstrapClienter.Cache)

	storage := pkg.NewStorage(bootstrapClienter.Storage)

	videoRepo := repositories.NewVideoRepo(bootstrapClienter.DB)

	crawler := service.NewCrawler(storage, videoRepo, config)

	ctx := context.Background()

	for {
		payload := entities.VideoQueuePayload{}

		redisBytesPayload, err := cache.BRpop(ctx, "video_queue")
		if err != nil {
			fmt.Println("err redis --> ", err)
			break
		}

		err = json.Unmarshal([]byte(redisBytesPayload), &payload)
		if err != nil {
			fmt.Println("err process redis payload --> ", err)
			continue
		}

		// Launch browser and interact
		err = crawler.RunBrowserAndInteract(ctx, payload.URL)

		if err != nil {
			fmt.Println("err processing --> ", err)
		}
	}
}
