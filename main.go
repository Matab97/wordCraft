package main

import (
	"wordCraft/db"
	"wordCraft/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env file if present (optional in k8s — env vars injected via Secret)
	godotenv.Load()
}

func main() {
	db.InitDB()
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost", "https://abbad.app", "https://www.abbad.app", "https://wordcraft.abbad.app"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// API routes only - no static files or templates
	routes.RegisterRoutes(r)
	r.Run(":8080") // API server on 8080
}
