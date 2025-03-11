package routes

import (
	"wordCraft/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/crafts/pair", CreateCraft)
	authenticated := r.Group("/")
	authenticated.Use(middlewares.Authenticate)
	authenticated.GET("/users", getUsers)
	authenticated.GET("/crafts", GetCrafts)
	r.POST("/signup", signup)
	r.POST("/login", login)
}
