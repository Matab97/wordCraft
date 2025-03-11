package main

import (
	"os"
	"wordCraft/db"
	"wordCraft/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Set DeepSeek API key
	os.Setenv("DEEPSEEK_API_KEY", "sk-fc663842cc6c42a4ad1ec1ceac0f8ce2")

	db.InitDB()
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost", "http://your-domain.com"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// API routes only - no static files or templates
	routes.RegisterRoutes(r)
	r.Run(":8080") // API server on 8080
}
