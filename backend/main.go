package main

import (
	"bedroomfm/internal/api"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: false,
	}))

	v1 := r.Group("/api")
	{
		v1.POST("/room/create", api.CreateRoom)
		v1.POST("/room/join", api.JoinRoom)
		v1.GET("/room/:id", api.GetRoom)
		v1.GET("/music/search", api.SearchMusic)
		v1.GET("/music/url", api.GetMusicURL)
		v1.GET("/music/lyric", api.GetLyric)
		v1.GET("/music/login/qr/key", api.QRKey)
		v1.GET("/music/login/qr/create", api.QRCreate)
		v1.GET("/music/login/qr/check", api.QRCheck)
	}

	r.GET("/ws/:roomId", api.WSHandler)

	log.Printf("BedroomFM backend listening on :%s", port)
	r.Run(":" + port)
}
