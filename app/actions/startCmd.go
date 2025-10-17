package actions

import (
	"fmt"
	"main/util"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StartCommand struct {
	Name   string
	Client tgbotapi.BotAPI
}

func (e StartCommand) getUserInfo(update tgbotapi.Update) error {
	_, err := util.GetOrCreateUser(util.UserInfo{}.New().FromAPIUser(update.Message.From))

	return err
}

func (e StartCommand) fabricateAnswer(update tgbotapi.Update) tgbotapi.MessageConfig {
	err := e.getUserInfo(update)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Error getting user info: "+err.Error())
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Bot is up")

	return msg
}

func (e StartCommand) Run(update tgbotapi.Update) error {
	message := e.fabricateAnswer(update)
	_, err := e.Client.Send(message)

	// Обрабатываем ошибку преобразования группы в супергруппу
	if util.IsSupergroupUpgradeError(err) {
		newChatID, handleErr := util.HandleSupergroupUpgrade(err, update.Message.Chat.ID)
		if handleErr != nil {
			return fmt.Errorf("failed to handle supergroup upgrade: %w", handleErr)
		}

		// Повторяем отправку сообщения с новым chat_id
		retryMessage := tgbotapi.NewMessage(newChatID, message.Text)
		_, err = e.Client.Send(retryMessage)
	}

	if err != nil {
		return err
	}

	return nil
}

func (e StartCommand) GetName() string {
	return e.Name
}
