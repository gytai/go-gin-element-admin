package system

import (
	"time"

	"gorm.io/gorm"
)

type SysAuthority struct {
	CreatedAt     time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt     time.Time      `json:"updatedAt" gorm:"comment:更新时间"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
	AuthorityId   uint           `json:"authorityId" gorm:"not null;unique;primaryKey;autoIncrement;comment:角色ID;size:90"`
	AuthorityName string         `json:"authorityName" gorm:"comment:角色名"`
	AuthorityCode string         `json:"authorityCode" gorm:"column:authority_code;comment:角色编码;unique"`
	ParentId      *uint          `json:"parentId" gorm:"comment:父角色ID"`
	DefaultRouter string         `json:"defaultRouter" gorm:"comment:默认菜单;default:dashboard"`
}

func (SysAuthority) TableName() string {
	return "sys_authorities"
}
