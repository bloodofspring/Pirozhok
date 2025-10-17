package main

import (
	"fmt"
	"log"
	"main/actions"
	"main/database"
	"main/filters"
	"main/handlers"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func connect(debug bool, apikey string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(apikey)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %w", err)
	}

	bot.Debug = debug
	log.Printf("Successfully authorized on account @%s", bot.Self.UserName)

	return bot, nil
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
	// Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: could not load .env file: %v", err)
	}

	// Получаем API ключ
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY environment variable is required")
	}

	// Подключаемся к боту
	client, err := connect(true, apiKey)
	if err != nil {
		log.Fatalf("Failed to connect to bot: %v", err)
	}

	// Инициализируем базу данных
	if err := database.InitDb(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Настраиваем graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Получаем обработчики
	act := getBotActions(*client)

	// Настраиваем получение обновлений
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := client.GetUpdatesChan(updateConfig)
	log.Println("Bot started successfully. Waiting for updates...")

	// Основной цикл с обработкой ошибок
	for {
		select {
		case update := <-updates:
			if err := act.HandleAll(update); err != nil {
				log.Printf("Error handling update: %v", err)
			}
		case sig := <-sigChan:
			log.Printf("Received signal %v, shutting down gracefully...", sig)
			return
		}
	}
}
