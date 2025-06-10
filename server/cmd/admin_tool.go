package main

import (
	"flag"
	"fmt"
	"server/core"
	"server/global"
	"server/initialize"
	"server/model/system"
	"os"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func main() {
	var action = flag.String("action", "", "操作类型: reset-admin-permissions, create-admin, show-admin-permissions")
	flag.Parse()

	if *action == "" {
		fmt.Println("使用方法:")
		fmt.Println("  go run cmd/admin_tool.go -action=reset-admin-permissions    # 重置管理员权限")
		fmt.Println("  go run cmd/admin_tool.go -action=create-admin               # 创建管理员用户")
		fmt.Println("  go run cmd/admin_tool.go -action=show-admin-permissions     # 显示管理员权限")
		os.Exit(1)
	}

	// 初始化配置和数据库
	global.VP = core.Viper()
	initialize.OtherInit()
	global.DB = initialize.Gorm()

	if global.DB == nil {
		fmt.Println("数据库连接失败")
		os.Exit(1)
	}

	switch *action {
	case "reset-admin-permissions":
		resetAdminPermissions()
	case "create-admin":
		createAdminUser()
	case "show-admin-permissions":
		showAdminPermissions()
	default:
		fmt.Printf("未知操作: %s\n", *action)
		os.Exit(1)
	}
}

// resetAdminPermissions 重置管理员权限
func resetAdminPermissions() {
	fmt.Println("开始重置管理员权限...")
	initialize.ResetAdminPermissions()
	fmt.Println("管理员权限重置完成")
}

// createAdminUser 创建管理员用户
func createAdminUser() {
	db := global.DB

	// 检查管理员角色是否存在
	var adminAuth system.SysAuthority
	err := db.Where("authority_id = ?", 888).First(&adminAuth).Error
	if err == gorm.ErrRecordNotFound {
		// 创建管理员角色
		adminAuth = system.SysAuthority{
			AuthorityId:   888,
			AuthorityName: "管理员",
			DefaultRouter: "dashboard",
		}
		if err := db.Create(&adminAuth).Error; err != nil {
			fmt.Printf("创建管理员角色失败: %v\n", err)
			return
		}
		fmt.Println("创建管理员角色成功")
	}

	// 检查管理员用户是否存在
	var adminUser system.SysUser
	if err := db.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			adminUser = system.SysUser{
				UUID:        uuid.New(),
				Username:    "admin",
				NickName:    "管理员",
				Password:    "$2a$10$ve6jMQzd7klAudy.LdDYOOEGMlOqq8zfvuOAut6FuRqAQzjgHh2LG", // 123456
				AuthorityId: 888,
				HeaderImg:   "https://qmplusimg.henrongyi.top/gva_header.jpg",
				SideMode:    "dark",
				Enable:      1,
			}
			if err := db.Create(&adminUser).Error; err != nil {
				fmt.Printf("创建管理员用户失败: %v\n", err)
				return
			}
			fmt.Println("创建管理员用户成功 (用户名: admin, 密码: 123456)")
		} else {
			fmt.Printf("查询管理员用户失败: %v\n", err)
			return
		}
	} else {
		fmt.Println("管理员用户已存在")
	}
}

// showAdminPermissions 显示管理员权限
func showAdminPermissions() {
	db := global.DB

	// 获取管理员角色信息
	var adminAuth system.SysAuthority
	if err := db.Where("authority_id = ?", 888).First(&adminAuth).Error; err != nil {
		fmt.Printf("管理员角色不存在: %v\n", err)
		return
	}

	fmt.Printf("管理员角色信息:\n")
	fmt.Printf("  角色ID: %d\n", adminAuth.AuthorityId)
	fmt.Printf("  角色名称: %s\n", adminAuth.AuthorityName)
	fmt.Printf("  默认路由: %s\n", adminAuth.DefaultRouter)

	// 获取管理员权限
	var permissions []struct {
		MenuId         uint   `json:"menuId"`
		Title          string `json:"title"`
		MenuType       string `json:"menuType"`
		PermissionCode string `json:"permissionCode"`
		Path           string `json:"path"`
		Sort           int    `json:"sort"`
	}

	err := db.Table("sys_authority_menus").
		Select("DISTINCT sys_base_menus.id as menu_id, sys_base_menus.title, sys_base_menus.menu_type, sys_base_menus.permission_code, sys_base_menus.path, sys_base_menus.sort").
		Joins("LEFT JOIN sys_base_menus ON sys_authority_menus.base_menu_id = sys_base_menus.id").
		Where("sys_authority_menus.authority_id = ?", 888).
		Order("sys_base_menus.menu_type, sys_base_menus.sort").
		Find(&permissions).Error

	if err != nil {
		fmt.Printf("获取管理员权限失败: %v\n", err)
		return
	}

	fmt.Printf("\n管理员权限列表 (共 %d 个):\n", len(permissions))
	fmt.Println("菜单权限:")
	for _, perm := range permissions {
		if perm.MenuType == "menu" {
			fmt.Printf("  [%d] %s (%s)\n", perm.MenuId, perm.Title, perm.Path)
		}
	}

	fmt.Println("\n按钮权限:")
	for _, perm := range permissions {
		if perm.MenuType == "button" {
			fmt.Printf("  [%d] %s (%s)\n", perm.MenuId, perm.Title, perm.PermissionCode)
		}
	}
}
