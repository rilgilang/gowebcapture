package main

import (
	"context"
	"fmt"
	"go-web-screen-record/pkg"
	"go-web-screen-record/service"
)

func main() {

	//https://satumomen.com/preview/peresmian-rs
	//https://joinedwithshan.viding.co/
	//https://app.sangmempelai.id/pilihan-tema/sunda-01
	//https://adirara.webnikah.com/?templatecoba=156/kepada:Budi%20dan%20Ani-Bandung
	//https://ourmoment.my.id/art-6

	cache := pkg.NewCache()

	ctx := context.Background()

	cache.Lpush(ctx, "video_queue", "https://ourmoment.my.id/art-6")

	for {
		urlLink, err := cache.BRpop(ctx, "video_queue")
		if err != nil {
			fmt.Println("err redis --> ", err)
			break
		}

		// Launch browser and interact
		err = service.RunBrowser(urlLink)

		if err != nil {
			fmt.Println("err processing --> ", err)
		}
	}
}
