package model

import (
	"react-demo-server/db"
)

type Menu struct {
	BaseModel
	Type     string `json:"type" binding:"required" comment:"菜单类型"`
	ParentID uint   `json:"parentId" default:"0" comment:"父级ID"`
	Name     string `json:"name" binding:"required" comment:"菜单名称"`

	Path string `json:"path" binding:"required" comment:"菜单路径"`
	Icon string `json:"icon" comment:"菜单图标"`

	Sort         float64 `json:"sort" default:"1.0" comment:"排序"`
	ShowSideMenu bool    `json:"showSideMenu" default:"true" comment:"是否显示侧边菜单"`
	Enabled      bool    `json:"enabled" default:"true" comment:"是否启用"`
}
type MenuWithChildren struct {
	Menu
	Children []MenuWithChildren `json:"children" comment:"子菜单"`
}

func (m *Menu) Get() error {
	return db.DB.Take(m, m.ID).Error
}

func GetMenuById(id uint) (Menu, error) {
	var menu Menu
	err := db.DB.Take(&menu, id).Error
	return menu, err
}

func GetMenus() ([]Menu, error) {
	var menus []Menu
	err := db.DB.Find(&menus).Error
	return menus, err
}

func GetMenu(menu *Menu) (Menu, error) {
	var menu1 Menu
	err := db.DB.Take(&menu1, menu.ID).Error
	return menu1, err
}

func CreateMenu(menu *Menu) error {
	return db.DB.Create(menu).Error
}
func UpdateMenu(menu *Menu) error {
	return db.DB.Save(menu).Error
}

func DeleteMenu(menuIDs []uint) error {
	return db.DB.Delete(&Menu{}, menuIDs).Error
}
