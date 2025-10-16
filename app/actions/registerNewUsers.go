package actions

import (
	"main/database"
	"main/database/models"
	"main/util"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type RegisterNewUsers struct {
	Name   string
	Client tgbotapi.BotAPI
}

func (r RegisterNewUsers) Run(update tgbotapi.Update) error {
	user, err := util.GetOrCreateUser(util.UserInfo{}.New().FromAPIUser(update.Message.From))
	if err != nil {
		return err
	}

	group, err := util.GetOrCreateGroup(update.Message.Chat.ID)
	if err != nil {
		return err
	}

	chatMemberConfig := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: group.TgId,
			UserID: user.TgId,
		},
	}

	chatMember, err := r.Client.GetChatMember(chatMemberConfig)
	if err != nil {
		return err
	}

	db := database.GetDB()
	_, err = db.Model(&models.GroyupParticipants{
		UserTgId: user.TgId,
		GroupTgId: group.TgId,
		IsAdmin: chatMember.IsCreator() || chatMember.IsAdministrator(),
	}).OnConflict("DO NOTHING").Insert()

	return err
}
