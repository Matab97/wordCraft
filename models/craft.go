package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"wordCraft/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Craft struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"Name"`
	Emoji       string             `bson:"emoji" json:"Emoji"`
	Combination string             `bson:"combination" json:"Combination"`
	CreatedAt   time.Time          `bson:"created_at" json:"CreatedAt"`
	IsNew       bool               `bson:"-" json:"IsNew"` // Transient field, not stored in DB
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type DeepSeekRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature"`
}

type DeepSeekResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func CraftPair(first string, second string) Craft {
	combination := first + "+" + second
	result := first + second
	return Craft{
		Name:        result,
		Combination: combination,
		Emoji:       "emoji",
	}
}

func AskAIToCraftPair(first string, second string) Craft {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		return Craft{Name: "Error: No API Key", Emoji: "❌"}
	}

	systemPrompt := "You are a creative word association expert. When given two unrelated words, combine them to return one existing word or short phrase (no new/made-up terms) that directly connects both concepts. Prioritize well-known references from real-world culture, media, literature, brands, or idioms. Avoid explanations or mentioning the original words in the answer, along with a relevant emoji. Always respond in the format: 'result,emoji'. For example: 'land' + 'sea' should give 'island,🏝️' and 'animal' + 'fire' should give 'dragon,🐉' and 'cloud' + 'water' should give 'cloud,🌧️'."
	userPrompt := fmt.Sprintf("Combine these two words creatively: '%s' and '%s'", first, second)

	reqBody := DeepSeekRequest{
		Model: "deepseek-chat",
		Messages: []Message{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: userPrompt,
			},
		},
		Temperature: 0.8,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return Craft{Name: "Error marshaling", Emoji: "❌"}
	}

	req, err := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return Craft{Name: "Error creating request", Emoji: "❌"}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Craft{Name: "Error making request", Emoji: "❌"}
	}
	defer resp.Body.Close()

	var aiResp DeepSeekResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		return Craft{Name: "Error decoding response", Emoji: "❌"}
	}

	if len(aiResp.Choices) == 0 {
		return Craft{Name: "No response from AI", Emoji: "❌"}
	}

	// Parse response (expected format: "word,emoji")
	response := aiResp.Choices[0].Message.Content
	var name, emoji string

	// Split the response on comma and handle potential spaces
	parts := strings.Split(response, ",")
	if len(parts) != 2 {
		return Craft{Name: "Error: Invalid response format", Emoji: "❌"}
	}

	name = strings.TrimSpace(parts[0])
	emoji = strings.TrimSpace(parts[1])

	IsNew := db.GetCollection("crafts").FindOne(context.Background(), bson.M{"name": name}).Err() != nil

	combination := first + "+" + second
	return Craft{
		Name:        name,
		Combination: combination,
		Emoji:       emoji,
		IsNew:       IsNew, // Set to true for new crafts
	}
}

func (c *Craft) Save() error {
	collection := db.GetCollection("crafts")
	c.CreatedAt = time.Now()
	_, err := collection.InsertOne(context.Background(), c)
	return err
}

func GetCrafts() []Craft {
	collection := db.GetCollection("crafts")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return []Craft{}
	}
	defer cursor.Close(context.Background())

	var crafts []Craft
	cursor.All(context.Background(), &crafts)
	return crafts
}

func GetCraftByName(first, second string) Craft {
	collection := db.GetCollection("crafts")
	var craft Craft

	filter := bson.M{"combination": first + "+" + second}
	err := collection.FindOne(context.Background(), filter).Decode(&craft)
	if err != nil {
		return Craft{}
	}
	return craft
}
