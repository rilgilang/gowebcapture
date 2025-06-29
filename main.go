package main

import (
	"fmt"
	"go-web-screen-record/service"
)

func main() {

	//e := echo.New()
	//e.GET("/", func(c echo.Context) error {
	//	ping := handlers.Ping(c)
	//	return c.JSON(http.StatusOK, struct {
	//		Data string `json:"data"`
	//	}{
	//		Data: ping,
	//	})
	//})
	//
	//e.POST("/to-video", func(c echo.Context) error {
	//	return c.String(http.StatusOK, "Hello, World!")
	//})
	//
	//e.Logger.Fatal(e.Start(":6969"))

	// Example links
	//https://satumomen.com/preview/peresmian-rs
	//https://joinedwithshan.viding.co/
	//https://app.sangmempelai.id/pilihan-tema/sunda-01
	//https://adirara.webnikah.com/?templatecoba=156/kepada:Budi%20dan%20Ani-Bandung
	//https://ourmoment.my.id/art-6

	// Launch browser and interact
	err := service.RunBrowser("https://wdp.namawebsite.net/auto-scroll-duplicate/")

	if err != nil {
		fmt.Println("err --> ", err)
	}

	fmt.Println("Recording complete.")
}
