package actions

import (
	"fmt"
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

	// Обрабатываем ошибку преобразования группы в супергруппу
	if util.IsSupergroupUpgradeError(err) {
		newChatID, handleErr := util.HandleSupergroupUpgrade(err, group.TgId)
		if handleErr != nil {
			return fmt.Errorf("failed to handle supergroup upgrade: %w", handleErr)
		}

		// Обновляем chat_id в конфигурации и повторяем запрос
		chatMemberConfig.ChatID = newChatID
		chatMember, err = r.Client.GetChatMember(chatMemberConfig)

		// Обновляем group.TgId для использования в базе данных
		group.TgId = newChatID
	}

	if err != nil {
		return err
	}

	db := database.GetDB()
	_, err = db.Model(&models.GroupParticipants{
		UserTgId:  user.TgId,
		GroupTgId: group.TgId,
		IsAdmin:   chatMember.IsCreator() || chatMember.IsAdministrator(),
	}).OnConflict("DO NOTHING").Insert()

	return err
}

func (r RegisterNewUsers) GetName() string {
	return r.Name
}
