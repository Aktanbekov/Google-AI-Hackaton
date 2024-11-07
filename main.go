package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func getEnv(name string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	envData := os.Getenv(name)
	if envData == "" {
		log.Fatal("Not Found")
	}

	return envData
}

func textRequest(token string, text string) string {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(token))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	prompt := text
	response, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatalf("Failed to generate content: %v", err)
	}

	if response.Candidates != nil {
		for _, v := range response.Candidates {
			if len(v.Content.Parts) > 0 {
				if generatedText, ok := v.Content.Parts[0].(genai.Text); ok {
					return string(generatedText)
				}
			}
		}
	}
	return "None"
}


func main() {
	token := getEnv("TOKEN_GEMINI")
	botToken := getEnv("TOKEN_BOT")

	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	updates, _ := bot.UpdatesViaLongPolling(nil)

	defer bot.StopLongPolling()

	for update := range updates {
		if update.Message != nil {
			chatID := tu.ID(update.Message.Chat.ID)
			aiResponse := textRequest(token, update.Message.Text)
			bot.SendMessage(tu.Message(chatID, aiResponse))
		}
		
	}
}
