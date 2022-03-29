package model

import (
	"react-demo-server/db"
)

type User struct {
	BaseModel

	UserName string `json:"username" comment:"用户名"`
	Password string `json:"password" comment:"密码"`
	Nickname string `json:"nickName" comment:"昵称"`
	Avatar   string `json:"avatar" comment:"头像"`

	// 权限
	Role string `json:"role" comment:"角色"`
}

func GetUserByUserName(username string) (User, error) {
	var user User
	err := db.DB.Take(&user, "username = ?", username).Error
	return user, err
}

func GetUser(id uint, needPassword bool) (User, error) {
	var user User
	omitColumns := []string{}
	if !needPassword {
		omitColumns = append(omitColumns, "password")
	}
	err := db.DB.Take(&user, id).Omit(omitColumns...).Error
	return user, err
}

func CreateUser(user User) error {
	return db.DB.Create(&user).Error
}

func UpdateUser(user User) error {
	return db.DB.Save(&user).Error
}

func DeleteUser(user User) error {
	return db.DB.Delete(&user).Error
}
