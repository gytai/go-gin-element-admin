package router

import (
	v1 "server/api/v1"

	"github.com/gin-gonic/gin"
)

func InitDashboardRouter(Router *gin.RouterGroup) {
	dashboardRouter := Router.Group("dashboard")
	{
		dashboardRouter.GET("stats", v1.GetDashboardStats)  // 获取仪表板统计数据
		dashboardRouter.GET("systemInfo", v1.GetSystemInfo) // 获取系统信息
	}
}
