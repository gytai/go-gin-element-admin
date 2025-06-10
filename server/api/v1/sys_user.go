package v1

import (
	"fmt"
	"net/http"
	"server/global"
	"server/model/system"
	"strconv"

	"server/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
	var user system.SysUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 检查用户名是否已存在
	var existUser system.SysUser
	if err := global.DB.Where("username = ?", user.Username).First(&existUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "用户名已存在",
		})
		return
	}

	// 加密密码
	user.UUID = uuid.New()
	user.Password = utils.BcryptHash(user.Password)

	if err := global.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建用户失败: " + err.Error(),
		})
		return
	}

	// 清除密码后返回用户信息
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "创建成功",
		"data": user,
	})
}

// GetUserList 获取用户列表
func GetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	username := c.Query("username")
	nickName := c.Query("nickName")

	var users []system.SysUser
	var total int64

	db := global.DB.Model(&system.SysUser{})

	// 搜索条件
	if username != "" {
		db = db.Where("username LIKE ?", "%"+username+"%")
	}
	if nickName != "" {
		db = db.Where("nick_name LIKE ?", "%"+nickName+"%")
	}

	// 获取总数
	db.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询失败: " + err.Error(),
		})
		return
	}

	// 清空密码字段
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "获取成功",
		"data": gin.H{
			"list":     users,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// GetUserById 根据ID获取用户
func GetUserById(c *gin.Context) {
	id := c.Param("id")
	var user system.SysUser

	if err := global.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "用户不存在",
		})
		return
	}

	user.Password = "" // 清空密码

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "获取成功",
		"data": user,
	})
}

// UpdateUser 更新用户
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user system.SysUser

	if err := global.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "用户不存在",
		})
		return
	}

	var updateData system.SysUser
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 调试信息：打印接收到的密码
	fmt.Printf("接收到的密码: '%s'\n", updateData.Password)

	// 检查是否是重置密码请求
	if updateData.Password == "RESET_PASSWORD_123456" {
		fmt.Println("检测到重置密码请求")
		// 特殊标识，表示重置密码为123456
		updateData.Password = utils.BcryptHash("123456")

		if err := global.DB.Model(&user).Update("password", updateData.Password).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "密码重置失败: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "密码重置成功，新密码为：123456",
		})
		return
	}

	// 如果更新密码，需要加密
	if updateData.Password != "" {
		updateData.Password = utils.BcryptHash(updateData.Password)
	}

	if err := global.DB.Model(&user).Updates(updateData).Error; err != nil {
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

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user system.SysUser

	if err := global.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "用户不存在",
		})
		return
	}

	if err := global.DB.Delete(&user).Error; err != nil {
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

// GetUserInfo 获取当前用户信息（个人中心用）
func GetUserInfo(c *gin.Context) {
	// 从JWT token中获取用户ID
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未找到用户信息",
		})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "用户信息格式错误",
		})
		return
	}

	var user system.SysUser
	if err := global.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "用户不存在",
		})
		return
	}

	user.Password = "" // 清空密码

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "获取成功",
		"data": user,
	})
}

// UpdateUserInfo 更新当前用户信息（个人中心用）
func UpdateUserInfo(c *gin.Context) {
	var updateData system.SysUser
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 从JWT token中获取用户ID
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未找到用户信息",
		})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "用户信息格式错误",
		})
		return
	}

	var user system.SysUser
	if err := global.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "用户不存在",
		})
		return
	}

	// 如果更新密码，需要加密
	if updateData.Password != "" {
		updateData.Password = utils.BcryptHash(updateData.Password)
	}

	// 使用Select方法明确指定要更新的字段，这样可以更新零值字段
	fieldsToUpdate := []string{"nick_name", "email", "phone", "header_img"}
	if updateData.Password != "" {
		fieldsToUpdate = append(fieldsToUpdate, "password")
	}

	if err := global.DB.Model(&user).Select(fieldsToUpdate).Updates(updateData).Error; err != nil {
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

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword     string `json:"oldPassword" binding:"required"`
		NewPassword     string `json:"newPassword" binding:"required"`
		ConfirmPassword string `json:"confirmPassword" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 验证新密码和确认密码是否一致
	if req.NewPassword != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "新密码和确认密码不一致",
		})
		return
	}

	// 从JWT token中获取用户ID
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未找到用户信息",
		})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "用户信息格式错误",
		})
		return
	}

	var user system.SysUser
	if err := global.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "用户不存在",
		})
		return
	}

	// 验证原密码
	if !utils.BcryptCheck(req.OldPassword, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "原密码错误",
		})
		return
	}

	// 加密新密码
	hashedPassword := utils.BcryptHash(req.NewPassword)

	// 更新密码
	if err := global.DB.Model(&user).Update("password", hashedPassword).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "密码修改失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "密码修改成功",
	})
}

// GetUserMenus 获取用户菜单（根据用户权限动态返回）
func GetUserMenus(c *gin.Context) {
	// 从JWT token中获取用户信息，如果没有则从查询参数获取（兼容性）
	var user system.SysUser
	var username string

	// 优先从JWT token获取用户信息
	if userIDInterface, exists := c.Get("userID"); exists {
		if userID, ok := userIDInterface.(uint); ok {
			if err := global.DB.Where("id = ?", userID).First(&user).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"code": 404,
					"msg":  "用户不存在",
				})
				return
			}
			username = user.Username
		}
	} else {
		// 如果没有JWT token，从查询参数获取用户名（兼容旧版本）
		username = c.DefaultQuery("username", "admin")
		if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  "用户不存在",
			})
			return
		}
	}

	// 获取用户角色的菜单权限
	var authorityMenus []system.SysAuthorityMenu
	if err := global.DB.Where("authority_id = ?", user.AuthorityId).Find(&authorityMenus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询用户权限失败: " + err.Error(),
		})
		return
	}

	// 调试信息
	fmt.Printf("用户: %s, 角色ID: %d, 权限数量: %d\n", username, user.AuthorityId, len(authorityMenus))

	// 如果没有配置权限，根据角色类型决定返回内容
	if len(authorityMenus) == 0 {
		// 只有管理员角色（888）才能看到所有菜单
		if user.AuthorityId == 888 {
			var allMenus []system.SysBaseMenu
			if err := global.DB.Order("sort ASC").Find(&allMenus).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": 500,
					"msg":  "查询菜单失败: " + err.Error(),
				})
				return
			}

			menuTree := buildUserMenuTree(allMenus, nil)
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  "获取成功（管理员全部权限）",
				"data": menuTree,
			})
			return
		} else {
			// 普通用户没有配置权限，只返回基础菜单（首页和个人中心）
			var basicMenus []system.SysBaseMenu
			if err := global.DB.Where("id IN ?", []uint{1, 2}).Order("sort ASC").Find(&basicMenus).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": 500,
					"msg":  "查询菜单失败: " + err.Error(),
				})
				return
			}

			menuTree := buildUserMenuTree(basicMenus, nil)
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  "获取成功（基础权限）",
				"data": menuTree,
			})
			return
		}
	}

	// 获取有权限的菜单ID列表
	var menuIds []uint
	for _, am := range authorityMenus {
		menuIds = append(menuIds, am.BaseMenuId)
	}

	// 调试信息
	fmt.Printf("菜单ID列表: %v\n", menuIds)

	// 查询对应的菜单信息
	var menus []system.SysBaseMenu
	if err := global.DB.Where("id IN ?", menuIds).Order("sort ASC").Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询菜单失败: " + err.Error(),
		})
		return
	}

	// 构建菜单树结构
	menuTree := buildUserMenuTree(menus, nil)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "获取成功",
		"data": menuTree,
	})
}

// buildUserMenuTree 构建用户菜单树结构（用于前端侧边栏显示和权限控制）
func buildUserMenuTree(menus []system.SysBaseMenu, parentId *uint) []system.SysBaseMenu {
	var tree []system.SysBaseMenu

	for _, menu := range menus {
		// 比较父ID：都为nil或者值相等
		if (menu.ParentId == nil && parentId == nil) ||
			(menu.ParentId != nil && parentId != nil && *menu.ParentId == *parentId) {
			menuId := menu.ID
			children := buildUserMenuTree(menus, &menuId)
			menu.Children = children
			tree = append(tree, menu)
		}
	}

	return tree
}

/*
// ResetPassword 重置用户密码为默认密码123456
// 注释掉，现在使用UpdateUser中的特殊处理来实现密码重置
func ResetPassword(c *gin.Context) {
	id := c.Param("id")
	var user system.SysUser

	if err := global.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "用户不存在",
		})
		return
	}

	// 设置默认密码为123456
	defaultPassword := "123456"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "密码加密失败",
		})
		return
	}

	// 更新密码
	if err := global.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "密码重置失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "密码重置成功，新密码为：123456",
	})
}
*/
