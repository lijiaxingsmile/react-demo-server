package model

import (
	"react-demo-server/db"
)

type Role struct {
	BaseModel
	Name string `json:"name" binding:"required" comment:"角色名称"`
	Code string `json:"code" binding:"required" comment:"角色编码"`
}

func (r *Role) Get() error {
	return db.DB.Take(r, r.ID).Error
}

func GetRoles() ([]Role, error) {
	var roles []Role
	err := db.DB.Find(&roles).Error
	return roles, err
}

func (r *Role) Update() error {
	return db.DB.Save(r).Error
}

func (r *Role) Delete() error {
	return db.DB.Delete(r).Error
}

func (r *Role) Create() error {
	return db.DB.Create(r).Error
}

func DeleteRoles(roleIDs []uint) error {
	return db.DB.Delete(&Role{}, roleIDs).Error
}
