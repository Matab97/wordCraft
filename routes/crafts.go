package routes

import (
	"context"
	"net/http"
	"wordCraft/db"
	"wordCraft/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetCrafts(c *gin.Context) {
	crafts := models.GetCrafts()
	c.JSON(http.StatusOK, crafts)
}

func CreateCraft(c *gin.Context) {
	first := c.Query("first")
	second := c.Query("second")

	if first == "" || second == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "first and second query params required"})
		return
	}

	craft := models.GetCraftByName(first, second)
	if craft.ID == primitive.NilObjectID {
		newCraft := models.AskAIToCraftPair(first, second)
		newCraft.Save()
		c.JSON(http.StatusCreated, newCraft)
		return
	}
	c.JSON(http.StatusOK, craft)
}

// GET /crafts/tree — full graph for discovery map
func GetCraftTree(c *gin.Context) {
	tree := models.GetCraftTree()
	c.JSON(http.StatusOK, tree)
}

// GET /crafts/:name/describe — encyclopedia entry
func DescribeCraft(c *gin.Context) {
	name := c.Param("name")

	// Try to get from DB first
	var craft models.Craft
	err := db.GetCollection("crafts").FindOne(context.Background(), bson.M{"name": name}).Decode(&craft)
	if err == nil && craft.Description != "" {
		c.JSON(http.StatusOK, gin.H{"name": craft.Name, "emoji": craft.Emoji, "description": craft.Description})
		return
	}

	// Generate fresh
	parts := []string{"", ""}
	if craft.Combination != "" {
		p := splitCombination(craft.Combination)
		parts = p
	}
	desc := models.GetCraftDescription(name, parts[0], parts[1])
	c.JSON(http.StatusOK, gin.H{"name": name, "emoji": craft.Emoji, "description": desc})
}

func splitCombination(combo string) []string {
	for i, ch := range combo {
		if ch == '+' {
			return []string{combo[:i], combo[i+1:]}
		}
	}
	return []string{combo, ""}
}

// GET /themes — all theme packs
func GetThemes(c *gin.Context) {
	c.JSON(http.StatusOK, models.ThemePacks)
}
