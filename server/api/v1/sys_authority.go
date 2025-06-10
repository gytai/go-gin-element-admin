package v1

import (
	"fmt"
	"net/http"
	"server/global"
	"server/model/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateAuthority 创建角色
func CreateAuthority(c *gin.Context) {
	var authority system.SysAuthority
	if err := c.ShouldBindJSON(&authority); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 验证必填字段
	if authority.AuthorityName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "角色名称不能为空",
		})
		return
	}

	if authority.AuthorityCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "角色编码不能为空",
		})
		return
	}

	// 检查角色编码是否已存在
	var existingRole system.SysAuthority
	if err := global.DB.Where("authority_code = ?", authority.AuthorityCode).First(&existingRole).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "角色编码已存在",
		})
		return
	}

	// 验证父角色是否存在
	if authority.ParentId != nil {
		var parentRole system.SysAuthority
		if err := global.DB.Where("authority_id = ?", *authority.ParentId).First(&parentRole).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "父角色不存在",
			})
			return
		}
	}

	// 设置默认路由
	if authority.DefaultRouter == "" {
		authority.DefaultRouter = "dashboard"
	}

	// AuthorityId由数据库自动生成，不需要前端提供
	authority.AuthorityId = 0

	if err := global.DB.Create(&authority).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "创建成功",
		"data": authority,
	})
}

// GetAuthorityList 获取角色列表
func GetAuthorityList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	authorityName := c.Query("authorityName")

	var authorities []system.SysAuthority
	var total int64

	db := global.DB.Model(&system.SysAuthority{})

	// 搜索条件
	if authorityName != "" {
		db = db.Where("authority_name LIKE ?", "%"+authorityName+"%")
	}

	// 获取总数
	db.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&authorities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询失败: " + err.Error(),
		})
		return
	}

	// 调试信息：打印查询到的数据
	for i, auth := range authorities {
		fmt.Printf("角色 %d: ID=%d, Name=%s, Code='%s', ParentId=%v\n",
			i+1, auth.AuthorityId, auth.AuthorityName, auth.AuthorityCode, auth.ParentId)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "获取成功",
		"data": gin.H{
			"list":     authorities,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// GetAuthorityById 根据ID获取角色
func GetAuthorityById(c *gin.Context) {
	id := c.Param("id")
	var authority system.SysAuthority

	if err := global.DB.Where("authority_id = ?", id).First(&authority).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "角色不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "获取成功",
		"data": authority,
	})
}

// UpdateAuthority 更新角色
func UpdateAuthority(c *gin.Context) {
	id := c.Param("id")
	var authority system.SysAuthority

	if err := global.DB.Where("authority_id = ?", id).First(&authority).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "角色不存在",
		})
		return
	}

	var updateData system.SysAuthority
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 如果更新角色编码，检查是否已存在
	if updateData.AuthorityCode != "" && updateData.AuthorityCode != authority.AuthorityCode {
		var existingRole system.SysAuthority
		if err := global.DB.Where("authority_code = ? AND authority_id != ?", updateData.AuthorityCode, id).First(&existingRole).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "角色编码已存在",
			})
			return
		}
	}

	// 验证父角色是否存在且不能是自己
	if updateData.ParentId != nil {
		// 不能将自己设置为父角色
		if *updateData.ParentId == authority.AuthorityId {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "不能将自己设置为父角色",
			})
			return
		}

		var parentRole system.SysAuthority
		if err := global.DB.Where("authority_id = ?", *updateData.ParentId).First(&parentRole).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "父角色不存在",
			})
			return
		}

		// 检查是否会形成循环引用
		if checkCircularReference(authority.AuthorityId, *updateData.ParentId) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "不能设置为子角色的父角色，这会形成循环引用",
			})
			return
		}
	}

	if err := global.DB.Model(&authority).Updates(updateData).Error; err != nil {
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

// 检查循环引用
func checkCircularReference(roleId, parentId uint) bool {
	var role system.SysAuthority
	if err := global.DB.Where("authority_id = ?", parentId).First(&role).Error; err != nil {
		return false
	}

	// 如果父角色就是当前角色，形成循环
	if role.ParentId != nil && *role.ParentId == roleId {
		return true
	}

	// 递归检查
	if role.ParentId != nil {
		return checkCircularReference(roleId, *role.ParentId)
	}

	return false
}

// DeleteAuthority 删除角色
func DeleteAuthority(c *gin.Context) {
	id := c.Param("id")
	var authority system.SysAuthority

	if err := global.DB.Where("authority_id = ?", id).First(&authority).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "角色不存在",
		})
		return
	}

	// 检查是否有用户关联此角色
	var userCount int64
	global.DB.Model(&system.SysUser{}).Where("authority_id = ?", id).Count(&userCount)
	if userCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "该角色下有用户，无法删除",
		})
		return
	}

	if err := global.DB.Delete(&authority).Error; err != nil {
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

// GetAllAuthorities 获取所有角色（下拉选项用）
func GetAllAuthorities(c *gin.Context) {
	var authorities []system.SysAuthority

	if err := global.DB.Find(&authorities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "获取成功",
		"data": authorities,
	})
}

// AssignMenus 给角色分配菜单权限
func AssignMenus(c *gin.Context) {
	authorityId := c.Param("id")

	var req struct {
		MenuIds []uint `json:"menuIds" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 验证角色是否存在
	var authority system.SysAuthority
	if err := global.DB.Where("authority_id = ?", authorityId).First(&authority).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "角色不存在",
		})
		return
	}

	// 如果有父角色，验证权限不能超过父角色
	if authority.ParentId != nil {
		// 获取父角色的权限
		var parentMenus []system.SysAuthorityMenu
		if err := global.DB.Where("authority_id = ?", *authority.ParentId).Find(&parentMenus).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "查询父角色权限失败: " + err.Error(),
			})
			return
		}

		// 构建父角色权限ID集合
		parentMenuIds := make(map[uint]bool)
		for _, pm := range parentMenus {
			parentMenuIds[pm.BaseMenuId] = true
		}

		// 检查当前分配的权限是否超出父角色
		for _, menuId := range req.MenuIds {
			if !parentMenuIds[menuId] {
				c.JSON(http.StatusBadRequest, gin.H{
					"code": 400,
					"msg":  "子角色权限不能超过父角色权限范围",
				})
				return
			}
		}
	}

	// 开启事务
	tx := global.DB.Begin()

	// 删除原有权限
	if err := tx.Where("authority_id = ?", authorityId).Delete(&system.SysAuthorityMenu{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "删除原权限失败: " + err.Error(),
		})
		return
	}

	// 分配新权限
	for _, menuId := range req.MenuIds {
		authorityMenu := system.SysAuthorityMenu{
			AuthorityId: authority.AuthorityId,
			BaseMenuId:  menuId,
		}
		if err := tx.Create(&authorityMenu).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "分配权限失败: " + err.Error(),
			})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "权限分配成功",
	})
}

// GetAuthorityMenus 获取角色的菜单权限
func GetAuthorityMenus(c *gin.Context) {
	authorityId := c.Param("id")

	// 验证角色是否存在
	var authority system.SysAuthority
	if err := global.DB.Where("authority_id = ?", authorityId).First(&authority).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "角色不存在",
		})
		return
	}

	// 获取角色拥有的菜单ID列表
	var authorityMenus []system.SysAuthorityMenu
	if err := global.DB.Where("authority_id = ?", authorityId).Find(&authorityMenus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询权限失败: " + err.Error(),
		})
		return
	}

	var menuIds []uint
	for _, am := range authorityMenus {
		menuIds = append(menuIds, am.BaseMenuId)
	}

	// 获取所有菜单，构建树结构，并标记已分配的菜单
	var allMenus []system.SysBaseMenu
	if err := global.DB.Order("sort ASC").Find(&allMenus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询菜单失败: " + err.Error(),
		})
		return
	}

	// 构建菜单树结构
	menuTree := buildMenuTreeWithPermission(allMenus, nil, menuIds)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "获取成功",
		"data": gin.H{
			"menuIds": menuIds,
			"menus":   menuTree,
		},
	})
}

// MenuTreeNode 带权限标记的菜单树节点
type MenuTreeNode struct {
	system.SysBaseMenu
	Checked  bool           `json:"checked"`
	Children []MenuTreeNode `json:"children"`
}

// buildMenuTreeWithPermission 构建带权限标记的菜单树
func buildMenuTreeWithPermission(menus []system.SysBaseMenu, parentId *uint, assignedIds []uint) []MenuTreeNode {
	var tree []MenuTreeNode

	for _, menu := range menus {
		// 比较父ID：都为nil或者值相等
		if (menu.ParentId == nil && parentId == nil) ||
			(menu.ParentId != nil && parentId != nil && *menu.ParentId == *parentId) {
			// 检查菜单是否已分配
			checked := false
			for _, id := range assignedIds {
				if menu.ID == id {
					checked = true
					break
				}
			}

			menuId := menu.ID
			node := MenuTreeNode{
				SysBaseMenu: menu,
				Checked:     checked,
				Children:    buildMenuTreeWithPermission(menus, &menuId, assignedIds),
			}
			tree = append(tree, node)
		}
	}

	return tree
}
