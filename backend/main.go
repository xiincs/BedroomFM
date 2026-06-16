package main

import (
	"bedroomfm/internal/api"
	"bedroomfm/internal/models"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Persistent user store — JSON file, survives restarts
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "data"
	}
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("cannot create data dir %s: %v", dataDir, err)
	}
	models.InitUserStore(filepath.Join(dataDir, "users.json"))

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: false,
	}))

	v1 := r.Group("/api")
	{
		v1.POST("/auth/register", api.Register)
		v1.POST("/auth/login", api.Login)
		v1.GET("/auth/me", api.GetMe)
		v1.POST("/user/xp", api.AddXP)

		v1.POST("/room/create", api.CreateRoom)
		v1.POST("/room/join", api.JoinRoom)
		v1.POST("/room/leave", api.LeaveRoom)
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
