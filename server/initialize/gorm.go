package initialize

import (
	"errors"
	"fmt"
	"log"
	"os"

	"server/global"
	"server/model/system"
	"server/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	Mysql = "mysql"
)

// Gorm 初始化数据库并产生数据库全局变量
func Gorm() *gorm.DB {
	switch global.CONFIG.System.DbType {
	case Mysql:
		return GormMysql()
	default:
		return GormMysql()
	}
}

// RegisterTables 注册数据库表专用
func RegisterTables() {
	db := global.DB

	// 检查表是否需要初始化
	needsInit := checkIfNeedsInitialization(db)

	// 自动迁移表结构（不删除现有数据）
	err := db.AutoMigrate(
		// 系统模块表 - 注意顺序：先创建被引用的表
		system.SysAuthority{},
		system.SysBaseMenu{},
		system.SysUser{},
		system.SysBaseMenuParameter{},
		system.SysBaseMenuBtn{},
		system.SysAuthorityMenu{},
		system.SysOperationLog{},
	)

	if err != nil {
		fmt.Printf("数据库表迁移失败: %v\n", err)
		os.Exit(0)
	}
	fmt.Println("数据库表结构检查完成")

	// 只有在需要初始化时才初始化基础数据
	if needsInit {
		fmt.Println("检测到数据库为空，开始初始化基础数据...")
		initData()
		fmt.Println("基础数据初始化完成")
	} else {
		fmt.Println("数据库已存在数据，跳过初始化")
	}
}

// checkIfNeedsInitialization 检查是否需要初始化数据库
func checkIfNeedsInitialization(db *gorm.DB) bool {
	// 检查关键表是否存在数据
	var userCount int64
	var authorityCount int64
	var menuCount int64

	// 如果表不存在，HasTable会返回false，这种情况下我们需要初始化
	if !db.Migrator().HasTable(&system.SysUser{}) ||
		!db.Migrator().HasTable(&system.SysAuthority{}) ||
		!db.Migrator().HasTable(&system.SysBaseMenu{}) {
		return true
	}

	// 检查关键表的数据量
	db.Model(&system.SysUser{}).Count(&userCount)
	db.Model(&system.SysAuthority{}).Count(&authorityCount)
	db.Model(&system.SysBaseMenu{}).Count(&menuCount)

	// 如果任一关键表为空，则需要初始化
	return userCount == 0 || authorityCount == 0 || menuCount == 0
}

// initData 初始化基础数据
func initData() {
	db := global.DB

	// 检查是否已经有管理员角色，如果没有则创建
	var authority system.SysAuthority
	if err := db.Where("authority_id = ?", 888).First(&authority).Error; err != nil {
		authority = system.SysAuthority{
			AuthorityId:   888,
			AuthorityName: "管理员",
			AuthorityCode: "admin",
			DefaultRouter: "dashboard",
		}
		fmt.Printf("准备创建管理员角色: %+v\n", authority)
		if err := db.Create(&authority).Error; err != nil {
			log.Println("创建管理员角色失败:", err)
		} else {
			log.Println("创建管理员角色成功")
		}
	} else {
		// 如果角色已存在但没有角色编码，则更新
		if authority.AuthorityCode == "" {
			authority.AuthorityCode = "admin"
			if err := db.Save(&authority).Error; err != nil {
				log.Println("更新管理员角色编码失败:", err)
			} else {
				log.Println("更新管理员角色编码成功")
			}
		}
	}

	// 检查是否已经有管理员用户，如果没有则创建
	var adminUser system.SysUser
	if err := db.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 生成正确的密码哈希
			hashedPassword := utils.BcryptHash("123456")

			adminUser = system.SysUser{
				UUID:        uuid.New(),
				Username:    "admin",
				NickName:    "管理员",
				Password:    hashedPassword, // 使用正确的密码哈希
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
		}
	} else {
		// 如果用户存在，强制更新为新的MD5密码格式
		fmt.Println("检测到管理员用户已存在，正在更新为新的MD5密码格式...")
		hashedPassword := utils.BcryptHash("123456")
		db.Model(&adminUser).Update("password", hashedPassword)
		fmt.Println("管理员密码已更新为MD5格式")
	}

	// 创建基础菜单
	initMenus(db)

	// 为管理员角色分配所有菜单权限
	initAdminPermissions(db)
}

// initMenus 初始化基础菜单
func initMenus(db *gorm.DB) {
	parentId3 := uint(3)   // 为子菜单定义父ID
	parentId4 := uint(4)   // 人员管理页面ID
	parentId5 := uint(5)   // 角色管理页面ID
	parentId6 := uint(6)   // 菜单管理页面ID
	parentId22 := uint(22) // 操作日志页面ID

	menus := []system.SysBaseMenu{
		// 顶级菜单
		{
			Model:     gorm.Model{ID: 1},
			ParentId:  nil, // 顶级菜单，父ID为nil
			Path:      "/dashboard",
			Name:      "Dashboard",
			Component: "views/dashboard/index.vue",
			Sort:      1,
			Title:     "首页",
			Icon:      "House",
			Hidden:    false,
			KeepAlive: true,
			MenuType:  "menu",
		},
		{
			Model:     gorm.Model{ID: 2},
			ParentId:  nil, // 顶级菜单，父ID为nil
			Path:      "/profile",
			Name:      "Profile",
			Component: "views/profile/index.vue",
			Sort:      2,
			Title:     "个人中心",
			Icon:      "User",
			Hidden:    false,
			KeepAlive: true,
			MenuType:  "menu",
		},
		{
			Model:     gorm.Model{ID: 3},
			ParentId:  nil, // 顶级菜单，父ID为nil
			Path:      "/system",
			Name:      "System",
			Component: "layout/index.vue",
			Sort:      3,
			Title:     "系统设置",
			Icon:      "Setting",
			Hidden:    false,
			KeepAlive: false,
			MenuType:  "menu",
		},

		// 系统设置子菜单
		{
			Model:     gorm.Model{ID: 4},
			ParentId:  &parentId3, // 子菜单，父ID为3
			Path:      "/system/user",
			Name:      "SystemUser",
			Component: "views/system/user/index.vue",
			Sort:      1,
			Title:     "人员管理",
			Icon:      "User",
			Hidden:    false,
			KeepAlive: true,
			MenuType:  "menu",
		},
		{
			Model:     gorm.Model{ID: 5},
			ParentId:  &parentId3, // 子菜单，父ID为3
			Path:      "/system/role",
			Name:      "SystemRole",
			Component: "views/system/role/index.vue",
			Sort:      2,
			Title:     "角色管理",
			Icon:      "UserFilled",
			Hidden:    false,
			KeepAlive: true,
			MenuType:  "menu",
		},
		{
			Model:     gorm.Model{ID: 6},
			ParentId:  &parentId3, // 子菜单，父ID为3
			Path:      "/system/menu",
			Name:      "SystemMenu",
			Component: "views/system/menu/index.vue",
			Sort:      3,
			Title:     "菜单管理",
			Icon:      "Menu",
			Hidden:    false,
			KeepAlive: true,
			MenuType:  "menu",
		},
		{
			Model:     gorm.Model{ID: 22},
			ParentId:  &parentId3, // 子菜单，父ID为3
			Path:      "/system/operation-log",
			Name:      "SystemOperationLog",
			Component: "views/system/operation-log/index.vue",
			Sort:      4,
			Title:     "操作日志",
			Icon:      "Document",
			Hidden:    false,
			KeepAlive: true,
			MenuType:  "menu",
		},

		// 人员管理按钮权限
		{
			Model:          gorm.Model{ID: 7},
			ParentId:       &parentId4,
			Title:          "新增用户",
			Sort:           1,
			Hidden:         true, // 按钮权限不在菜单中显示
			MenuType:       "button",
			PermissionCode: "user:create",
		},
		{
			Model:          gorm.Model{ID: 8},
			ParentId:       &parentId4,
			Title:          "编辑用户",
			Sort:           2,
			Hidden:         true,
			MenuType:       "button",
			PermissionCode: "user:update",
		},
		{
			Model:          gorm.Model{ID: 9},
			ParentId:       &parentId4,
			Title:          "删除用户",
			Sort:           3,
			Hidden:         true,
			MenuType:       "button",
			PermissionCode: "user:delete",
		},
		{
			Model:          gorm.Model{ID: 10},
			ParentId:       &parentId4,
			Title:          "查看用户详情",
			Sort:           4,
			Hidden:         true,
			MenuType:       "button",
			PermissionCode: "user:view",
		},
		{
			Model:          gorm.Model{ID: 11},
			ParentId:       &parentId4,
			Title:          "重置密码",
			Sort:           5,
			Hidden:         true,
			MenuType:       "button",
			PermissionCode: "user:reset_password",
		},

		// 角色管理按钮权限
		{
			Model:          gorm.Model{ID: 12},
			ParentId:       &parentId5,
			Title:          "新增角色",
			Sort:           1,
			Hidden:         true,
			MenuType:       "button",
			PermissionCode: "role:create",
		},
		{
			Model:          gorm.Model{ID: 13},
			ParentId:       &parentId5,
			Title:          "编辑角色",
			Sort:           2,
			Hidden:         true,
			MenuType:       "button",
			PermissionCode: "role:update",
		},
		{
			Model:          gorm.Model{ID: 14},
			ParentId:       &parentId5,
			Title:          "删除角色",
			Sort:           3,
			Hidden:         true,
			MenuType:       "button",
			PermissionCode: "role:delete",
		},
		{
			Model:          gorm.Model{ID: 15},
			ParentId:       &parentId5,
			Title:          "分配菜单权限",
			Sort:           4,
			Hidden:         true,
			MenuType:       "button",
			PermissionCode: "role:assign_menu",
		},
		{
			Model:          gorm.Model{ID: 16},
			ParentId:       &parentId5,
			Title:          "查看角色详情",
			Sort:           5,
			Hidden:         true,
			MenuType:       "button",
			PermissionCode: "role:view",
		},

		// 菜单管理按钮权限
		{
			Model:          gorm.Model{ID: 17},
			ParentId:       &parentId6,
			Title:          "新增菜单",
			Sort:           1,
			Hidden:         true,
			MenuType:       "button",
			PermissionCode: "menu:create",
		},
		{
			Model:          gorm.Model{ID: 18},
			ParentId:       &parentId6,
			Title:          "编辑菜单",
			Sort:           2,
			Hidden:         true,
			MenuType:       "button",
			PermissionCode: "menu:update",
		},
		{
			Model:          gorm.Model{ID: 19},
			ParentId:       &parentId6,
			Title:          "删除菜单",
			Sort:           3,
			Hidden:         true,
			MenuType:       "button",
			PermissionCode: "menu:delete",
		},
		{
			Model:          gorm.Model{ID: 20},
			ParentId:       &parentId6,
			Title:          "查看菜单详情",
			Sort:           4,
			Hidden:         true,
			MenuType:       "button",
			PermissionCode: "menu:view",
		},
		{
			Model:          gorm.Model{ID: 21},
			ParentId:       &parentId6,
			Title:          "添加子菜单",
			Sort:           5,
			Hidden:         true,
			MenuType:       "button",
			PermissionCode: "menu:create_child",
		},

		// 操作日志按钮权限
		{
			Model:          gorm.Model{ID: 23},
			ParentId:       &parentId22,
			Title:          "查看日志详情",
			Sort:           1,
			Hidden:         true,
			MenuType:       "button",
			PermissionCode: "operation-log:view",
		},
		{
			Model:          gorm.Model{ID: 24},
			ParentId:       &parentId22,
			Title:          "删除日志",
			Sort:           2,
			Hidden:         true,
			MenuType:       "button",
			PermissionCode: "operation-log:delete",
		},
		{
			Model:          gorm.Model{ID: 25},
			ParentId:       &parentId22,
			Title:          "批量删除日志",
			Sort:           3,
			Hidden:         true,
			MenuType:       "button",
			PermissionCode: "operation-log:batch-delete",
		},
		{
			Model:          gorm.Model{ID: 26},
			ParentId:       &parentId22,
			Title:          "清空日志",
			Sort:           4,
			Hidden:         true,
			MenuType:       "button",
			PermissionCode: "operation-log:clear",
		},
		{
			Model:          gorm.Model{ID: 27},
			ParentId:       &parentId22,
			Title:          "导出日志",
			Sort:           5,
			Hidden:         true,
			MenuType:       "button",
			PermissionCode: "operation-log:export",
		},
	}

	for _, menu := range menus {
		var existMenu system.SysBaseMenu
		if err := db.Where("id = ?", menu.ID).First(&existMenu).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&menu).Error; err != nil {
					fmt.Printf("创建菜单失败: %v\n", err)
				}
			}
		}
	}
	fmt.Println("基础菜单初始化完成")
}

// initAdminPermissions 为管理员角色分配所有菜单权限
func initAdminPermissions(db *gorm.DB) {
	// 检查管理员角色是否已有权限
	var existingCount int64
	if err := db.Model(&system.SysAuthorityMenu{}).Where("authority_id = ?", 888).Count(&existingCount).Error; err != nil {
		fmt.Printf("检查管理员权限失败: %v\n", err)
		return
	}

	// 如果管理员已有权限，则不进行自动初始化（避免覆盖用户的权限设置）
	if existingCount > 0 {
		fmt.Printf("管理员角色已有 %d 个权限，跳过自动初始化\n", existingCount)
		return
	}

	// 只有在管理员角色没有任何权限时，才进行初始化
	fmt.Println("管理员角色暂无权限，开始初始化...")

	// 获取所有菜单ID（包括菜单和按钮）
	var allMenuIds []uint
	if err := db.Model(&system.SysBaseMenu{}).Select("id").Find(&allMenuIds).Error; err != nil {
		fmt.Printf("获取菜单ID失败: %v\n", err)
		return
	}

	// 为管理员角色分配所有菜单权限
	var authorityMenus []system.SysAuthorityMenu
	for _, menuId := range allMenuIds {
		authorityMenus = append(authorityMenus, system.SysAuthorityMenu{
			AuthorityId: 888,
			BaseMenuId:  menuId,
		})
	}

	// 批量插入权限关联关系
	if err := db.Create(&authorityMenus).Error; err != nil {
		fmt.Printf("分配管理员权限失败: %v\n", err)
		return
	}

	fmt.Printf("管理员权限初始化完成，分配了 %d 个菜单权限\n", len(authorityMenus))
}

// ResetAdminPermissions 重置管理员权限（强制重新分配所有权限）
func ResetAdminPermissions() {
	db := global.DB
	if db == nil {
		fmt.Println("数据库连接不可用")
		return
	}

	fmt.Println("开始重置管理员权限...")

	// 删除管理员角色的所有现有权限
	if err := db.Where("authority_id = ?", 888).Delete(&system.SysAuthorityMenu{}).Error; err != nil {
		fmt.Printf("删除管理员现有权限失败: %v\n", err)
		return
	}

	// 获取所有菜单ID（包括菜单和按钮）
	var allMenuIds []uint
	if err := db.Model(&system.SysBaseMenu{}).Select("id").Find(&allMenuIds).Error; err != nil {
		fmt.Printf("获取菜单ID失败: %v\n", err)
		return
	}

	// 为管理员角色分配所有菜单权限
	var authorityMenus []system.SysAuthorityMenu
	for _, menuId := range allMenuIds {
		authorityMenus = append(authorityMenus, system.SysAuthorityMenu{
			AuthorityId: 888,
			BaseMenuId:  menuId,
		})
	}

	// 批量插入权限关联关系
	if err := db.Create(&authorityMenus).Error; err != nil {
		fmt.Printf("重新分配管理员权限失败: %v\n", err)
		return
	}

	fmt.Printf("管理员权限重置完成，重新分配了 %d 个菜单权限\n", len(authorityMenus))
}
