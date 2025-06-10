package system

import (
	"gorm.io/gorm"
)

// SysAuthorityMenu 角色菜单关联表
type SysAuthorityMenu struct {
	gorm.Model
	AuthorityId uint `json:"authorityId" gorm:"comment:角色ID"`
	BaseMenuId  uint `json:"baseMenuId" gorm:"comment:菜单ID"`
}

func (SysAuthorityMenu) TableName() string {
	return "sys_authority_menus"
}
