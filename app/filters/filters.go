package filters

import (
	"main/database"
	"main/database/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func IsGroup(update tgbotapi.Update) bool {
	return update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup"
}

func IsGroupAdmin(update tgbotapi.Update) bool {
	if !IsGroup(update) {
		return false
	}

	db := database.GetDB()
	var participantinfo models.GroyupParticipants
	err := db.Model(&participantinfo).Where("group_tg_id = ? AND user_tg_id = ?", update.Message.Chat.ID, update.Message.From.ID).Select()
	if err != nil {
		return false
	}

	return participantinfo.IsAdmin
}

func CallCommand(update tgbotapi.Update) bool {
	return update.Message.IsCommand() && update.Message.Command() == "call"
}

func StartCommand(update tgbotapi.Update) bool {
	return update.Message.IsCommand() && update.Message.Command() == "start"
}
