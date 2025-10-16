package actions

import (
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
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Error getting user info: " + err.Error())
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Bot is up")

	return msg
}

func (e StartCommand) Run(update tgbotapi.Update) error {
	if _, err := e.Client.Send(e.fabricateAnswer(update)); err != nil {
		return err
	}

	return nil
}

func (e StartCommand) GetName() string {
	return e.Name
}
