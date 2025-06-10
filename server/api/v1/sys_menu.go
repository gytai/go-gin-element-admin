package v1

import (
	"server/global"
	"server/model/system"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateMenu 创建菜单
func CreateMenu(c *gin.Context) {
	var menu system.SysBaseMenu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	if err := global.DB.Create(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建菜单失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "创建成功",
		"data": menu,
	})
}

// GetMenuList 获取菜单列表
func GetMenuList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	title := c.Query("title")
	path := c.Query("path")

	var menus []system.SysBaseMenu
	var total int64

	db := global.DB.Model(&system.SysBaseMenu{})

	// 搜索条件
	if title != "" {
		db = db.Where("title LIKE ?", "%"+title+"%")
	}
	if path != "" {
		db = db.Where("path LIKE ?", "%"+path+"%")
	}

	// 获取总数
	db.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Order("sort ASC").Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "获取成功",
		"data": gin.H{
			"list":     menus,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// GetMenuById 根据ID获取菜单
func GetMenuById(c *gin.Context) {
	id := c.Param("id")
	var menu system.SysBaseMenu

	if err := global.DB.Where("id = ?", id).First(&menu).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "菜单不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "获取成功",
		"data": menu,
	})
}

// UpdateMenu 更新菜单
func UpdateMenu(c *gin.Context) {
	id := c.Param("id")
	var menu system.SysBaseMenu

	if err := global.DB.Where("id = ?", id).First(&menu).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "菜单不存在",
		})
		return
	}

	var updateData system.SysBaseMenu
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	if err := global.DB.Model(&menu).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "更新成功",
	})
}

// DeleteMenu 删除菜单
func DeleteMenu(c *gin.Context) {
	id := c.Param("id")
	var menu system.SysBaseMenu

	if err := global.DB.Where("id = ?", id).First(&menu).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "菜单不存在",
		})
		return
	}

	// 检查是否有子菜单
	var childCount int64
	global.DB.Model(&system.SysBaseMenu{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "该菜单下有子菜单，无法删除",
		})
		return
	}

	if err := global.DB.Delete(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "删除失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "删除成功",
	})
}

// GetMenuTree 获取菜单树结构
func GetMenuTree(c *gin.Context) {
	var menus []system.SysBaseMenu

	if err := global.DB.Order("sort ASC").Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询失败: " + err.Error(),
		})
		return
	}

	// 构建树结构
	tree := buildMenuTree(menus, nil)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "获取成功",
		"data": tree,
	})
}

// buildMenuTree 构建菜单树结构
func buildMenuTree(menus []system.SysBaseMenu, parentId *uint) []system.SysBaseMenu {
	var tree []system.SysBaseMenu

	for _, menu := range menus {
		// 比较父ID：都为nil或者值相等
		if (menu.ParentId == nil && parentId == nil) ||
			(menu.ParentId != nil && parentId != nil && *menu.ParentId == *parentId) {
			menuId := menu.ID
			children := buildMenuTree(menus, &menuId)
			menu.Children = children
			tree = append(tree, menu)
		}
	}

	return tree
}
