package main

import (
	"BE_FAKE_REVIEW/config"
	"BE_FAKE_REVIEW/routes"
	"BE_FAKE_REVIEW/services"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db, err := services.ConnectDatabase(cfg)
	if err != nil {
		log.Fatal("Gagal koneksi database:", err)
	}

	if err := services.SetupDatabase(db); err != nil {
		log.Fatal("Gagal setup database:", err)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	routes.SetupRoutes(r, cfg, db)

	log.Println("Server berjalan di port", cfg.AppPort)

	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
