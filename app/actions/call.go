package actions

import (
	"fmt"
	"main/database"
	"main/database/models"
	"main/util"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SummonAllUsers struct {
	Name   string
	Client tgbotapi.BotAPI
}

func (s SummonAllUsers) Run(update tgbotapi.Update) error {
	initiatorUser, err := util.GetOrCreateUser(util.UserInfo{}.New().FromAPIUser(update.Message.From))
	if err != nil {
		return err
	}

	group, err := util.GetOrCreateGroup(update.Message.Chat.ID)
	if err != nil {
		return err
	}

	db := database.GetDB()
	var initiatorUserStatus models.GroupParticipants
	err = db.Model(&initiatorUserStatus).
		Column("is_admin").
		Where("user_tg_id = ?", initiatorUser.TgId).
		Where("group_tg_id = ?", group.TgId).
		Select()
	if err != nil {
		return err
	}

	// Раскомментировать, если нужно проверять права инициатора
	// if !initiatorUserStatus.IsAdmin {
	// 	message := tgbotapi.NewMessage(update.Message.Chat.ID, "You are not an admin")
	// 	message.ReplyToMessageID = update.Message.MessageID
	// 	_, err := s.Client.Send(message)

	// 	return err
	// }

	err = db.Model(&group).
		WherePK().
		Relation("Users").
		Select()
	if err != nil {
		return err
	}

	textMessage := "Summoning all users in the group...\n"
	for _, user := range group.Users {
		textMessage += fmt.Sprintf("<a href=\"tg://user?id=%d\">|</a>", user.TgId)
	}

	message := tgbotapi.NewMessage(update.Message.Chat.ID, textMessage)
	message.ParseMode = "HTML"
	_, err = s.Client.Send(message)

	return err
}

func (s SummonAllUsers) GetName() string {
	return s.Name
}
