package initialize

import (
	"fmt"
	"net/http"

	v1 "server/api/v1"
	"server/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// 初始化总路由
func Routers() *gin.Engine {
	Router := gin.Default()

	// 跨域配置
	Router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
	}))

	// 操作日志中间件
	Router.Use(middleware.OperationLogMiddlewareWithBody())

	// 健康监测
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"status":  "running",
		})
	})

	// 静态文件服务 - 头像文件访问
	Router.Static("/uploads", "./uploads")

	// 公开路由组
	PublicGroup := Router.Group("/api")
	{
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "ok",
				"status":  "running",
			})
		})

		// 基础路由（无需认证）
		BaseGroup := PublicGroup.Group("/base")
		{
			BaseGroup.POST("/login", v1.Login)
			BaseGroup.POST("/logout", v1.Logout)
			BaseGroup.GET("/captcha", v1.Captcha)
		}

		// 文件上传路由
		UploadGroup := PublicGroup.Group("/upload")
		{
			UploadGroup.POST("/avatar", v1.UploadAvatar)
		}

		// 仪表盘路由
		DashboardGroup := PublicGroup.Group("/dashboard")
		{
			DashboardGroup.GET("/stats", v1.GetDashboardStats)
			DashboardGroup.GET("/systemInfo", v1.GetSystemInfo)
		}
		fmt.Println("Dashboard routes initialized")

		// 系统管理路由
		SystemGroup := PublicGroup.Group("/system")
		{
			// 用户管理
			UserGroup := SystemGroup.Group("/user")
			{
				UserGroup.GET("/list", v1.GetUserList)
				UserGroup.GET("/info", middleware.JWTAuth(), v1.GetUserInfo)
				UserGroup.GET("/menus", middleware.JWTAuth(), v1.GetUserMenus)
				UserGroup.PUT("/info", middleware.JWTAuth(), v1.UpdateUserInfo)
				UserGroup.PUT("/password", middleware.JWTAuth(), v1.ChangePassword)

				UserGroup.POST("", v1.CreateUser)
				UserGroup.GET("/:id", v1.GetUserById)
				UserGroup.PUT("/:id", v1.UpdateUser)
				UserGroup.DELETE("/:id", v1.DeleteUser)
			}

			// 角色管理
			RoleGroup := SystemGroup.Group("/role")
			{
				RoleGroup.GET("/list", v1.GetAuthorityList)
				RoleGroup.GET("/all", v1.GetAllAuthorities)
				RoleGroup.POST("", v1.CreateAuthority)
				RoleGroup.GET("/:id/menus", v1.GetAuthorityMenus)
				RoleGroup.POST("/:id/menus", v1.AssignMenus)
				RoleGroup.GET("/:id", v1.GetAuthorityById)
				RoleGroup.PUT("/:id", v1.UpdateAuthority)
				RoleGroup.DELETE("/:id", v1.DeleteAuthority)
			}

			// 菜单管理
			MenuGroup := SystemGroup.Group("/menu")
			{
				MenuGroup.GET("/list", v1.GetMenuList)
				MenuGroup.GET("/tree", v1.GetMenuTree)
				MenuGroup.POST("", v1.CreateMenu)
				MenuGroup.GET("/:id", v1.GetMenuById)
				MenuGroup.PUT("/:id", v1.UpdateMenu)
				MenuGroup.DELETE("/:id", v1.DeleteMenu)
			}

			// 操作日志管理 - 暂时注释掉，需要重新实现
			// 操作日志管理路由
			OperationLogGroup := SystemGroup.Group("/operation-log")
			{
				OperationLogGroup.GET("/list", v1.GetOperationLogList)
				OperationLogGroup.GET("/stats", v1.GetOperationStats)
				OperationLogGroup.GET("/export", v1.ExportOperationLogs)
				OperationLogGroup.GET("/:id", v1.GetOperationLogById)
				OperationLogGroup.DELETE("/:id", v1.DeleteOperationLog)
				OperationLogGroup.DELETE("/batch", v1.DeleteOperationLogsByIds)
				OperationLogGroup.DELETE("/clear", v1.ClearOperationLogs)
				OperationLogGroup.DELETE("/clear-by-days", v1.ClearOperationLogsByDays)
			}
		}
		fmt.Println("System routes initialized")

		// 暂时移除其他业务路由，后续需要时再添加
		// router.InitBusinessRouter(PublicGroup)
	}

	fmt.Println("        欢迎使用 go-gin-element-admin")
	fmt.Println("        当前版本:v1.0.0")
	fmt.Println("        默认自动化文档地址:http://127.0.0.1:8888/swagger/index.html")
	fmt.Println("        默认前端文件运行地址:http://127.0.0.1:8080")

	return Router
}
