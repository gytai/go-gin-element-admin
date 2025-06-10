package system

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SysUser struct {
	gorm.Model
	UUID        uuid.UUID `json:"uuid" gorm:"index;comment:用户UUID"`
	Username    string    `json:"username" gorm:"index;comment:用户登录名"`
	Password    string    `json:"password" gorm:"comment:用户登录密码"`
	NickName    string    `json:"nickName" gorm:"default:系统用户;comment:用户昵称"`
	SideMode    string    `json:"sideMode" gorm:"default:dark;comment:用户侧边主题"`
	HeaderImg   string    `json:"headerImg" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"`
	BaseColor   string    `json:"baseColor" gorm:"default:#fff;comment:基础颜色"`
	ActiveColor string    `json:"activeColor" gorm:"default:#1890ff;comment:活跃颜色"`
	AuthorityId uint      `json:"authorityId" gorm:"default:888;index;comment:用户角色ID"`
	Phone       string    `json:"phone" gorm:"comment:用户手机号"`
	Email       string    `json:"email" gorm:"comment:用户邮箱"`
	Enable      int       `json:"enable" gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"`
}

func (SysUser) TableName() string {
	return "sys_users"
}
