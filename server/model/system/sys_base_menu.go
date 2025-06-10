package system

import (
	"gorm.io/gorm"
)

type SysBaseMenu struct {
	gorm.Model
	MenuLevel   uint   `json:"-"`
	ParentId    *uint  `json:"parentId" gorm:"comment:父菜单ID"`
	Path        string `json:"path" gorm:"comment:路由path"`
	Name        string `json:"name" gorm:"comment:路由name"`
	Hidden      bool   `json:"hidden" gorm:"comment:是否在列表隐藏"`
	Component   string `json:"component" gorm:"comment:对应前端文件路径"`
	Sort        int    `json:"sort" gorm:"comment:排序标记"`
	ActiveName  string `json:"activeName" gorm:"comment:高亮菜单"`
	KeepAlive   bool   `json:"keepAlive" gorm:"comment:是否缓存"`
	DefaultMenu bool   `json:"defaultMenu" gorm:"comment:是否是基础路由(开发中)"`
	Title       string `json:"title" gorm:"comment:菜单名"`
	Icon        string `json:"icon" gorm:"comment:菜单图标"`
	CloseTab    bool   `json:"closeTab" gorm:"comment:自动关闭tab"`

	// 新增权限编码字段，用于按钮权限控制
	PermissionCode string `json:"permissionCode" gorm:"comment:权限编码，按钮类型必填"`
	MenuType       string `json:"menuType" gorm:"comment:菜单类型 menu:菜单 button:按钮"`

	Children   []SysBaseMenu          `json:"children" gorm:"-"`
	Parameters []SysBaseMenuParameter `json:"parameters"`
	MenuBtn    []SysBaseMenuBtn       `json:"menuBtn"`
}

type SysBaseMenuParameter struct {
	gorm.Model
	SysBaseMenuID uint   `json:"sysBaseMenuID" gorm:"comment:菜单ID"`
	Type          string `json:"type" gorm:"comment:地址栏携带参数为params还是query"`
	Key           string `json:"key" gorm:"comment:地址栏携带参数的key"`
	Value         string `json:"value" gorm:"comment:地址栏携带参数的值"`
}

type SysBaseMenuBtn struct {
	gorm.Model
	Name          string `json:"name" gorm:"comment:按钮关键key"`
	Desc          string `json:"desc" gorm:"comment:按钮备注"`
	SysBaseMenuID uint   `json:"sysBaseMenuID" gorm:"comment:菜单ID"`
}

func (SysBaseMenu) TableName() string {
	return "sys_base_menus"
}
