package actions

import (
	"main/database"
	"main/database/models"
	"main/util"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type RegisterLeftUsers struct {
	Name   string
	Client tgbotapi.BotAPI
}

func (r RegisterLeftUsers) Run(update tgbotapi.Update) error {
	user, err := util.GetOrCreateUser(util.UserInfo{}.New().FromAPIUser(update.Message.From))
	if err != nil {
		return err
	}

	group, err := util.GetOrCreateGroup(update.Message.Chat.ID)
	if err != nil {
		return err
	}

	db := database.GetDB()
	_, err = db.Model(&models.GroupParticipants{}).
		Where("user_tg_id = ?", user.TgId).
		Where("group_tg_id = ?", group.TgId).
		Delete()

	return err
}

func (r RegisterLeftUsers) GetName() string {
	return r.Name
}
