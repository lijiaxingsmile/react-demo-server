package model

import (
	"react-demo-server/db"

	"gorm.io/gorm"
)

type Menu struct {
	gorm.Model

	ParentID uint   `json:"parentId" comment:"父级ID"`
	Name     string `json:"name" binding:"required" comment:"菜单名称"`

	Path      string `json:"path" binding:"required" comment:"菜单路径"`
	Component string `json:"component" binding:"required" comment:"菜单组件"`
	Icon      string `json:"icon" comment:"菜单图标"`
	Layout    string `json:"layout" comment:"菜单布局"`

	Sort float64 `json:"sort" binding:"required" comment:"排序"`
}
type MenuWithChildren struct {
	Menu
	Children []MenuWithChildren `json:"children" comment:"子菜单"`
}

func GetMenus() ([]Menu, error) {
	var menus []Menu
	err := db.DB.Find(&menus).Error
	return menus, err
}

func GetMenu(menu *Menu) (Menu, error) {
	var menu1 Menu
	err := db.DB.Where(menu).First(&menu1).Error
	return menu1, err
}

func CreateMenu(menu *Menu) error {
	return db.DB.Create(menu).Error
}
func UpdateMenu(menu *Menu) error {
	return db.DB.Save(menu).Error
}

func DeleteMenu(menuIDs []uint) error {
	return db.DB.Where("id in (?)", menuIDs).Delete(Menu{}).Error
}
