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
	IsNew       bool               `bson:"-" json:"IsNew"`
	Description string             `bson:"description,omitempty" json:"Description,omitempty"`
	Theme       string             `bson:"theme,omitempty" json:"Theme,omitempty"`
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

// Theme packs — each has a name, emoji, and starter words
type ThemePack struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Emoji       string   `json:"emoji"`
	Description string   `json:"description"`
	Words       []string `json:"words"`
}

var ThemePacks = []ThemePack{
	{
		ID:          "elements",
		Name:        "Elements",
		Emoji:       "🌊",
		Description: "The four classical elements and forces of nature",
		Words:       []string{"water", "fire", "earth", "air", "metal", "wood", "light", "dark"},
	},
	{
		ID:          "science",
		Name:        "Science",
		Emoji:       "🔬",
		Description: "Explore scientific concepts and discoveries",
		Words:       []string{"atom", "energy", "gravity", "wave", "cell", "gene", "star", "plasma", "magnet", "laser"},
	},
	{
		ID:          "nature",
		Name:        "Nature",
		Emoji:       "🌿",
		Description: "Plants, animals, and the natural world",
		Words:       []string{"forest", "ocean", "mountain", "cloud", "rain", "sun", "seed", "root", "wind", "stone"},
	},
	{
		ID:          "history",
		Name:        "History",
		Emoji:       "🏛️",
		Description: "Ancient civilizations and historical concepts",
		Words:       []string{"empire", "sword", "scroll", "trade", "war", "temple", "coin", "ship", "king", "myth"},
	},
	{
		ID:          "tech",
		Name:        "Technology",
		Emoji:       "💻",
		Description: "Modern tech and digital concepts",
		Words:       []string{"code", "data", "network", "signal", "chip", "cloud", "robot", "screen", "wire", "byte"},
	},
}

func askDeepSeek(prompt string, system string) (string, error) {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("no API key")
	}

	reqBody := DeepSeekRequest{
		Model: "deepseek-chat",
		Messages: []Message{
			{Role: "system", Content: system},
			{Role: "user", Content: prompt},
		},
		Temperature: 0.8,
	}

	jsonData, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var aiResp DeepSeekResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		return "", err
	}
	if len(aiResp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}
	return aiResp.Choices[0].Message.Content, nil
}

func AskAIToCraftPair(first string, second string) Craft {
	system := "You are a creative word association expert. When given two words, combine them to return one existing word or short phrase that directly connects both concepts. Prioritize well-known references from real-world culture, media, literature, brands, or idioms. Always respond ONLY in the format: 'result,emoji'. No explanation. Example: 'land' + 'sea' → 'island,🏝️'"
	prompt := fmt.Sprintf("Combine: '%s' + '%s'", first, second)

	response, err := askDeepSeek(prompt, system)
	if err != nil {
		return Craft{Name: "Error: " + err.Error(), Emoji: "❌"}
	}

	parts := strings.SplitN(response, ",", 2)
	if len(parts) != 2 {
		return Craft{Name: "Error: Invalid format", Emoji: "❌"}
	}

	name := strings.TrimSpace(parts[0])
	emoji := strings.TrimSpace(parts[1])
	combination := first + "+" + second
	isNew := db.GetCollection("crafts").FindOne(context.Background(), bson.M{"name": name}).Err() != nil

	// Get encyclopedia description async-style (generate with same call pattern)
	description := GetCraftDescription(name, first, second)

	return Craft{
		Name:        name,
		Combination: combination,
		Emoji:       emoji,
		IsNew:       isNew,
		Description: description,
	}
}

func GetCraftDescription(name, first, second string) string {
	system := "You are a fun educational assistant. Given a word that was created by combining two other words, provide a short 2-sentence description: one sentence explaining what the word means, one fun fact or connection to real-world knowledge. Keep it engaging and suitable for all ages. Respond with plain text only, no formatting."
	prompt := fmt.Sprintf("The word '%s' was created by combining '%s' and '%s'. Describe it.", name, first, second)

	desc, err := askDeepSeek(prompt, system)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(desc)
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

func GetCraftTree() []map[string]interface{} {
	crafts := GetCrafts()
	var nodes []map[string]interface{}
	for _, c := range crafts {
		parts := strings.SplitN(c.Combination, "+", 2)
		node := map[string]interface{}{
			"id":    c.Name,
			"name":  c.Name,
			"emoji": c.Emoji,
		}
		if len(parts) == 2 {
			node["parents"] = []string{parts[0], parts[1]}
		}
		nodes = append(nodes, node)
	}
	return nodes
}
