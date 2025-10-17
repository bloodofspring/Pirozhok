package util

import (
	"fmt"
	"main/database"
	"main/database/models"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserInfo struct {
	TgID     int64
	UserName string
	FullName string
}

func (u UserInfo) New() *UserInfo {
	return &UserInfo{}
}

func (u *UserInfo) FromAPIUser(user *tgbotapi.User) *UserInfo {
	return &UserInfo{
		TgID:     user.ID,
		UserName: user.UserName,
		FullName: user.FirstName + " " + user.LastName,
	}
}

func (u UserInfo) ToModel() *models.Users {
	return &models.Users{
		TgId:     u.TgID,
		UserName: u.UserName,
		FullName: u.FullName,
	}
}

func GetOrCreateUser(userInfo *UserInfo) (models.Users, error) {
	db := database.GetDB()

	var user models.Users

	err := db.Model(&user).Where("tg_id = ?", userInfo.TgID).Select()
	if err == nil {
		return user, nil
	}

	userToInsert := userInfo.ToModel()
	_, err = db.Model(userToInsert).OnConflict("DO NOTHING").Returning("*").Insert()
	if err != nil {
		return models.Users{}, err
	}

	return *userToInsert, nil
}

func GetOrCreateGroup(chatID int64) (models.Groups, error) {
	db := database.GetDB()

	var group models.Groups

	err := db.Model(&group).Where("tg_id = ?", chatID).Select()
	if err == nil {
		return group, nil
	}

	groupToInsert := &models.Groups{
		TgId: chatID,
	}

	_, err = db.Model(groupToInsert).OnConflict("DO NOTHING").Returning("*").Insert()
	if err != nil {
		return models.Groups{}, err
	}

	return *groupToInsert, nil
}

// IsSupergroupUpgradeError проверяет, является ли ошибка результатом преобразования группы в супергруппу
func IsSupergroupUpgradeError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "group chat was upgraded to a supergroup chat")
}

// ConvertToSupergroupID конвертирует обычный chat_id в supergroup chat_id
func ConvertToSupergroupID(chatID int64) int64 {
	// Если chat_id уже имеет префикс супергруппы, возвращаем как есть
	if chatID < -1000000000000 {
		return chatID
	}
	// Для отрицательных chat_id (обычные группы) добавляем префикс -100
	if chatID < 0 {
		// Убираем знак минус, добавляем префикс -100, затем возвращаем знак минус
		absID := -chatID
		return -1000000000000 - absID
	}
	// Для положительных chat_id добавляем префикс -100
	return -1000000000000 - chatID
}

// UpdateGroupChatID обновляет chat_id группы в базе данных
func UpdateGroupChatID(oldChatID, newChatID int64) error {
	db := database.GetDB()

	// Обновляем основную таблицу групп
	_, err := db.Model(&models.Groups{}).
		Set("tg_id = ?", newChatID).
		Where("tg_id = ?", oldChatID).
		Update()
	if err != nil {
		return fmt.Errorf("failed to update group chat_id: %w", err)
	}

	// Обновляем таблицу участников группы
	_, err = db.Model(&models.GroupParticipants{}).
		Set("group_tg_id = ?", newChatID).
		Where("group_tg_id = ?", oldChatID).
		Update()
	if err != nil {
		return fmt.Errorf("failed to update group participants chat_id: %w", err)
	}

	return nil
}

// HandleSupergroupUpgrade обрабатывает ошибку преобразования группы в супергруппу
func HandleSupergroupUpgrade(err error, chatID int64) (int64, error) {
	if !IsSupergroupUpgradeError(err) {
		return chatID, err
	}

	newChatID := ConvertToSupergroupID(chatID)

	// Обновляем chat_id в базе данных
	if updateErr := UpdateGroupChatID(chatID, newChatID); updateErr != nil {
		return chatID, fmt.Errorf("failed to handle supergroup upgrade: %w", updateErr)
	}

	return newChatID, nil
}
