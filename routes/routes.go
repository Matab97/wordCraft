package routes

import (
	"wordCraft/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// Public craft routes
	r.GET("/crafts/pair", CreateCraft)
	r.GET("/crafts/tree", GetCraftTree)
	r.GET("/crafts/:name/describe", DescribeCraft)
	r.GET("/themes", GetThemes)

	// Auth routes
	r.POST("/signup", signup)
	r.POST("/login", login)

	// Authenticated routes
	authenticated := r.Group("/")
	authenticated.Use(middlewares.Authenticate)
	authenticated.GET("/users", getUsers)
	authenticated.GET("/crafts", GetCrafts)
}
