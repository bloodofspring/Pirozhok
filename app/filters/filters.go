package filters

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func IsGroup(update tgbotapi.Update) bool {
	return update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup"
}

func IsMessageFromGroup(update tgbotapi.Update) bool {
	return update.Message != nil && IsGroup(update)
}

func CallCommand(update tgbotapi.Update) bool {
	return update.Message.IsCommand() && update.Message.Command() == "call"
}

func StartCommand(update tgbotapi.Update) bool {
	return update.Message.IsCommand() && update.Message.Command() == "start"
}

func NewChatMember(update tgbotapi.Update) bool {
	return update.Message != nil && update.Message.NewChatMembers != nil && len(update.Message.NewChatMembers) > 0
}

func LeftChatMember(update tgbotapi.Update) bool {
	return update.Message != nil && update.Message.LeftChatMember != nil
}
