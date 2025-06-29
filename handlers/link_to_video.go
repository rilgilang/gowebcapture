package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go-web-screen-record/service"
)

func LinkToVideo(c echo.Context) string {
	// Example links
	//https://satumomen.com/preview/peresmian-rs
	//https://joinedwithshan.viding.co/
	//https://app.sangmempelai.id/pilihan-tema/sunda-01
	//https://adirara.webnikah.com/?templatecoba=156/kepada:Budi%20dan%20Ani-Bandung
	//https://ourmoment.my.id/art-6

	// Launch browser and interact
	err := service.RunBrowser("https://ourmoment.my.id/art-6")

	if err != nil {
		fmt.Println("err --> ", err)
	}

	fmt.Println("Recording complete.")

	return "blok goblok"
}
