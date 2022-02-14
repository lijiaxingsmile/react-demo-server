package api

import (
	"react-demo-server/model"
	"react-demo-server/util"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 测试
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// 获取用户信息
func GetUser(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID := claims["id"].(float64)
	user, err := model.GetUser(uint(userID), false)
	user.Password = ""

	if err != nil {
		c.JSON(200, util.Response(false, "无效的用户信息", nil))
		return
	}

	c.JSON(200, util.Response(true, "获取用户信息成功", gin.H{
		"user": user,
	}))
}

// 获取菜单树
func GetMenus(c *gin.Context) {

	menus, err := model.GetMenus()
	if err != nil {
		c.JSON(200, util.Response(false, "获取菜单信息失败", nil))
		return
	}
	// 构造菜单树
	menuTree := GenerateMenuTree(0, menus)

	c.JSON(200, util.Response(true, "获取菜单成功", menuTree))
}

// 获取菜单以及子菜单
func GetMenu(c *gin.Context) {

	var menu model.Menu
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		c.JSON(200, util.Response(false, "菜单信息不能为空", nil))
		return
	}

	menus, err := model.GetMenus()
	if err != nil {
		c.JSON(200, util.Response(false, "获取菜单信息失败", nil))
		return
	}

	menuTree := GenerateMenuTree(menu.ID, menus)

	c.JSON(200, util.Response(true, "获取菜单成功", menuTree[0]))
}

// 生成菜单树
func GenerateMenuTree(menuID uint, menus []model.Menu) []model.MenuWithChildren {
	var menuTree []model.MenuWithChildren
	for _, menu := range menus {
		if menu.ParentID == menuID {
			menuTree = append(menuTree, model.MenuWithChildren{
				Menu:     menu,
				Children: GenerateMenuTree(menu.ID, menus),
			})
		}
	}
	return menuTree
}

// 创建菜单
func CreateMenu(c *gin.Context) {

	var menu model.Menu
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		c.JSON(200, util.Response(false, "菜单信息不能为空", nil))
		return
	}

	err = model.CreateMenu(&menu)
	if err != nil {
		c.JSON(200, util.Response(false, "创建菜单失败", nil))
		return
	}

	c.JSON(200, util.Response(true, "创建菜单成功", nil))
}

// 更新菜单
func UpdateMenu(c *gin.Context) {

	var menu model.Menu
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		c.JSON(200, util.Response(false, "菜单信息不能为空", nil))
		return
	}

	err = model.UpdateMenu(&menu)
	if err != nil {
		c.JSON(200, util.Response(false, "更新菜单失败", nil))
		return
	}

	c.JSON(200, util.Response(true, "更新菜单成功", nil))
}

// 获取菜单的所有子菜单ID
func GetMenuChildrenIds(menuId uint, menus []model.Menu, menuIds []uint) []uint {
	for _, menu := range menus {
		if menu.ParentID == menuId {
			menuIds = append(menuIds, menu.ID)
			menuIds = GetMenuChildrenIds(menu.ID, menus, menuIds)
		}
	}
	return menuIds
}

// 删除菜单以及子菜单
func DeleteMenu(c *gin.Context) {

	var menu model.Menu
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		c.JSON(200, util.Response(false, "菜单信息不能为空", nil))
		return
	}
	menus, err := model.GetMenus()
	if err != nil {
		c.JSON(200, util.Response(false, "获取菜单信息失败", nil))
		return
	}

	menuChildrenSlice := GetMenuChildrenIds(menu.ID, menus, []uint{})

	err = model.DeleteMenu(menuChildrenSlice)
	if err != nil {
		c.JSON(200, util.Response(false, "删除菜单失败", nil))
		return
	}

	c.JSON(200, util.Response(true, "删除菜单成功", nil))
}
