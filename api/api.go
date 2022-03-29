package api

import (
	"fmt"
	"log"
	"net/http"
	"react-demo-server/model"
	"react-demo-server/util"
	"strconv"
	"strings"

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
		c.JSON(401, util.Response(false, "无效的用户信息", nil))
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

func GetMenu(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(404, util.Response(false, "找不到对应菜单信息", nil))
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(404, util.Response(false, "找不到对应菜单信息", nil))
		return
	}
	fmt.Println("idInt:", idInt)

	menu := model.Menu{}
	menu.ID = uint(idInt)
	err = menu.Get()
	if err != nil {
		c.JSON(200, util.Response(false, "找不到对应菜单信息", nil))
		return
	}
	c.JSON(200, util.Response(true, "获取菜单成功", menu))
}

// 获取菜单以及子菜单
func GetMenuAndChildren(c *gin.Context) {

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

	if _, err = model.GetMenu(&menu); err != nil {
		c.JSON(200, util.Response(false, "找不到对应菜单信息", nil))
		return
	}

	// 对于菜单路劲不以绝对路径开头的,前面添加父级路径
	if !strings.HasPrefix(menu.Path, "/") {
		parentMenu, err := model.GetMenuById(menu.ParentID)
		if err != nil {
			c.JSON(200, util.Response(false, "找不到对应菜单信息", nil))
			return
		}
		menu.Path = parentMenu.Path + "/" + menu.Path
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

// 批量删除菜单以及子菜单
func BatchDeleteMenus(c *gin.Context) {

	idsQuery := c.Query("ids")
	ids := strings.Split(idsQuery, ",")

	if len(ids) == 0 {
		c.JSON(200, util.Response(false, "菜单ID不能为空", nil))
		return
	}

	idsSlice := make([]uint, 0, len(ids))

	for _, id := range ids {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.Response(false, "菜单ID不能为空", nil))
			return
		}
		idsSlice = append(idsSlice, uint(idInt))
	}

	menus, err := model.GetMenus()
	if err != nil {
		c.JSON(200, util.Response(false, "获取菜单信息失败", nil))
		return
	}

	for _, id := range idsSlice {
		menuChildrenSlice := GetMenuChildrenIds(id, menus, []uint{})
		idsSlice = append(idsSlice, menuChildrenSlice...)
	}

	log.Println("需要删除的ID为:", idsSlice)
	if len(idsSlice) == 0 {
		c.JSON(200, util.Response(false, "获取菜单信息失败", nil))
		return
	}

	if err = model.DeleteMenu(idsSlice); err != nil {
		c.JSON(500, util.Response(false, "删除菜单失败", nil))
		return
	}
	c.JSON(200, util.Response(true, "删除菜单成功", nil))
}

// 获取角色列表
func GetRoles(c *gin.Context) {
	roles, err := model.GetRoles()
	if err != nil {
		c.JSON(200, util.Response(false, "获取角色列表失败", nil))
		return
	}
	c.JSON(200, util.Response(true, "获取角色列表成功", roles))
}

func GetRole(c *gin.Context) {

	id := c.Query("id")
	if id == "" {
		c.JSON(404, util.Response(false, "找不到对应角色信息", nil))
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(404, util.Response(false, "找不到对应角色信息", nil))
		return
	}

	role := &model.Role{}
	role.ID = uint(idInt)

	err = role.Get()
	if err != nil {
		c.JSON(200, util.Response(false, "获取角色信息失败", nil))
		return
	}
	c.JSON(200, util.Response(true, "获取角色信息成功", role))
}

func UpdateRole(c *gin.Context) {
	var role model.Role
	err := c.ShouldBindJSON(&role)
	if err != nil {
		c.JSON(200, util.Response(false, "角色信息不能为空", nil))
		return
	}
	err = role.Update()
	if err != nil {
		c.JSON(200, util.Response(false, "更新角色信息失败", nil))
		return
	}
	c.JSON(200, util.Response(true, "更新角色信息成功", nil))
}

func DeleteRole(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(404, util.Response(false, "找不到对应角色信息", nil))
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(404, util.Response(false, "找不到对应角色信息", nil))
		return
	}

	role := &model.Role{}
	role.ID = uint(idInt)

	err = role.Delete()
	if err != nil {
		c.JSON(200, util.Response(false, "删除角色信息失败", nil))
		return
	}
	c.JSON(200, util.Response(true, "删除角色信息成功", nil))
}

func BatchDeleteRoles(c *gin.Context) {
	idsQuery := c.Query("id")
	ids := strings.Split(idsQuery, ",")

	if len(ids) == 0 {
		c.JSON(200, util.Response(false, "角色ID不能为空", nil))
		return
	}

	idsSlice := make([]uint, 0, len(ids))
	for _, id := range ids {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.Response(false, "角色ID不能为空", nil))
			return
		}
		idsSlice = append(idsSlice, uint(idInt))
	}

	err := model.DeleteRoles(idsSlice)
	if err != nil {
		c.JSON(200, util.Response(false, "删除角色失败", nil))
		return
	}
	c.JSON(200, util.Response(true, "删除角色成功", nil))
}

func CreateRole(c *gin.Context) {
	var role model.Role
	err := c.ShouldBindJSON(&role)
	if err != nil {
		c.JSON(200, util.Response(false, "角色信息不能为空", nil))
		return
	}
	err = role.Create()
	if err != nil {
		c.JSON(200, util.Response(false, "创建角色失败", nil))
		return
	}
	c.JSON(200, util.Response(true, "创建角色成功", nil))
}
