package main

import (
	"context"
	"encoding/json"
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"github.com/labstack/echo/v4"
	"go/src/github.com/rilgilang/gowebcapture/bootstrap"
	"go/src/github.com/rilgilang/gowebcapture/entities"
	"go/src/github.com/rilgilang/gowebcapture/pkg"
	"go/src/github.com/rilgilang/gowebcapture/repositories"
	"go/src/github.com/rilgilang/gowebcapture/service"
	"log"
)

func main() {

	bootstrapClienter, config, err := bootstrap.Setup()
	if err != nil {
		panic(err)
	}

	cache := pkg.NewCache(bootstrapClienter.Cache)

	storage := pkg.NewStorage(bootstrapClienter.Storage)

	videoRepo := repositories.NewVideoRepo(bootstrapClienter.DB)

	server := socketio.NewServer(nil)

	socket := pkg.NewSocket(server)

	crawler := service.NewCrawler(storage, socket, videoRepo, config)

	ctx := context.Background()

	go socketServer(server)

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
		err = crawler.RunBrowserAndInteract(ctx, payload.UniqueId, payload.URL)

		if err != nil {
			fmt.Println("err processing --> ", err)
		}
	}
}

func socketServer(server *socketio.Server) {
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())
		return nil
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("closed", reason)
	})

	go server.Serve()
	defer server.Close()

	e := echo.New()
	e.HideBanner = true

	e.Static("/", "../asset")
	e.Any("/socket.io/", func(context echo.Context) error {
		server.ServeHTTP(context.Response(), context.Request())
		return nil
	})
	e.Logger.Fatal(e.Start(":8000"))
}
