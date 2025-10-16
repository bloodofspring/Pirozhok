package util

import (
	"main/database"
	"main/database/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserInfo struct {
	TgID int64
	UserName string
	FullName string
}

func (u UserInfo) New() *UserInfo {
	return &UserInfo{}
}

func (u *UserInfo) FromAPIUser(user *tgbotapi.User) *UserInfo {
	return &UserInfo{
		TgID: user.ID,
		UserName: user.UserName,
		FullName: user.FirstName + " " + user.LastName,
	}
}

func (u UserInfo) ToModel() *models.Users {
	return &models.Users{
		TgId: u.TgID,
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

	err := db.Model(&group).Where("chat_id = ?", chatID).Select()
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
