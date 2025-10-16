package main

import (
	"log"
	"main/actions"
	"main/database"
	"main/filters"
	"main/handlers"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func connect(debug bool, apikey string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(apikey)
	if err != nil {
		panic(err)
	}

	bot.Debug = debug
	log.Printf("Successfully authorized on account @%s", bot.Self.UserName)

	return bot
}

func getBotActions(bot tgbotapi.BotAPI) handlers.ActiveHandlers {
	act := handlers.ActiveHandlers{Handlers: []handlers.Handler{
		handlers.CommandHandler.Product(actions.StartCommand{Name: "start-cmd", Client: bot}, []handlers.Filter{filters.StartCommand}),
		handlers.CommandHandler.Product(actions.SummonAllUsers{Name: "summon-all-users-cmd", Client: bot}, []handlers.Filter{filters.CallCommand}),
		handlers.MessageHandler.Product(actions.RegisterNewUsers{Name: "register-new-users", Client: bot}, []handlers.Filter{filters.NewChatMember}),
		handlers.MessageHandler.Product(actions.RegisterNewUsers{Name: "register-new-users", Client: bot}, []handlers.Filter{filters.IsMessageFromGroup}),
		handlers.MessageHandler.Product(actions.RegisterLeftUsers{Name: "register-left-users", Client: bot}, []handlers.Filter{filters.LeftChatMember}),
	}}

	return act
}

func main() {
	_ = godotenv.Load()
	client := connect(true, os.Getenv("API_KEY"))
	act := getBotActions(*client)

	err := database.InitDb()
	if err != nil {
		log.Fatalf("error initializing database: %v", err)
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := client.GetUpdatesChan(updateConfig)
	for update := range updates {
		_ = act.HandleAll(update)
	}
}
