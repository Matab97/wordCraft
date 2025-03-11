package routes

import (
	"net/http"
	"wordCraft/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetCrafts(c *gin.Context) {
	// TODO: Implement getting crafts
	crafts := models.GetCrafts()
	c.JSON(http.StatusOK, crafts)
}

func CreateCraft(c *gin.Context) {
	craft := models.GetCraftByName(c.Query("first"), c.Query("second"))
	if craft.ID == primitive.NilObjectID { // Check for zero value since Craft is a struct
		newCraft := models.AskAIToCraftPair(c.Query("first"), c.Query("second"))
		newCraft.Save()
		c.JSON(http.StatusCreated, newCraft)
		return
	}
	c.JSON(http.StatusOK, craft)
}
