package v1

import (
	"server/global"
	"server/model/system"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DashboardStats 仪表板统计数据结构
type DashboardStats struct {
	UserCount   int64 `json:"userCount"`   // 用户总数
	RoleCount   int64 `json:"roleCount"`   // 角色总数
	MenuCount   int64 `json:"menuCount"`   // 菜单总数
	OnlineCount int64 `json:"onlineCount"` // 在线用户数（暂时模拟）
}

// SystemInfo 系统信息结构
type SystemInfo struct {
	SystemVersion string `json:"systemVersion"`
	GoVersion     string `json:"goVersion"`
	GinVersion    string `json:"ginVersion"`
	VueVersion    string `json:"vueVersion"`
	ElementPlus   string `json:"elementPlus"`
}

// GetDashboardStats 获取仪表板统计数据
func GetDashboardStats(c *gin.Context) {
	var stats DashboardStats

	// 获取用户总数
	global.DB.Model(&system.SysUser{}).Count(&stats.UserCount)

	// 获取角色总数
	global.DB.Model(&system.SysAuthority{}).Count(&stats.RoleCount)

	// 获取菜单总数
	global.DB.Model(&system.SysBaseMenu{}).Count(&stats.MenuCount)

	// 在线用户数（这里暂时使用模拟数据，实际可以通过Redis或其他方式统计）
	stats.OnlineCount = 12

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": stats,
		"msg":  "获取成功",
	})
}

// GetSystemInfo 获取系统信息
func GetSystemInfo(c *gin.Context) {
	info := SystemInfo{
		SystemVersion: "v1.0.0",
		GoVersion:     "1.21",
		GinVersion:    "1.9.1",
		VueVersion:    "3.4.21",
		ElementPlus:   "2.6.3",
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": info,
		"msg":  "获取成功",
	})
}
