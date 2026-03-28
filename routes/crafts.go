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

// GET /crafts/:name/describe — encyclopedia entry (DB-cached)
func DescribeCraft(c *gin.Context) {
	name := c.Param("name")

	// Try DB first
	var craft models.Craft
	err := db.GetCollection("crafts").FindOne(context.Background(), bson.M{"name": name}).Decode(&craft)
	if err == nil && craft.Description != "" {
		// Cache hit — no AI call
		c.JSON(http.StatusOK, gin.H{"name": craft.Name, "emoji": craft.Emoji, "description": craft.Description, "cached": true})
		return
	}

	// Cache miss — generate and persist
	parts := []string{"", ""}
	if craft.Combination != "" {
		parts = splitCombination(craft.Combination)
	}
	desc := models.GetCraftDescription(name, parts[0], parts[1])

	// Save to DB so next request is free
	models.UpdateDescription(name, desc)

	c.JSON(http.StatusOK, gin.H{"name": name, "emoji": craft.Emoji, "description": desc, "cached": false})
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
