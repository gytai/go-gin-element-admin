package router

import (
	v1 "server/api/v1"

	"github.com/gin-gonic/gin"
)

func InitSystemRouter(Router *gin.RouterGroup) {
	systemRouter := Router.Group("system")
	{
		// 用户管理路由
		userRouter := systemRouter.Group("user")
		{
			userRouter.GET("list", v1.GetUserList)        // 获取用户列表
			userRouter.GET("info", v1.GetUserInfo)        // 获取当前用户信息
			userRouter.GET("menus", v1.GetUserMenus)      // 获取用户菜单
			userRouter.PUT("info", v1.UpdateUserInfo)     // 更新当前用户信息
			userRouter.PUT("password", v1.ChangePassword) // 修改密码
			userRouter.GET(":id", v1.GetUserById)         // 根据ID获取用户
			userRouter.POST("", v1.CreateUser)            // 创建用户
			userRouter.PUT(":id", v1.UpdateUser)          // 更新用户
			userRouter.DELETE(":id", v1.DeleteUser)       // 删除用户
		}

		// 角色管理路由
		roleRouter := systemRouter.Group("role")
		{
			roleRouter.GET("list", v1.GetAuthorityList)       // 获取角色列表
			roleRouter.GET("all", v1.GetAllAuthorities)       // 获取所有角色（下拉选项）
			roleRouter.POST("", v1.CreateAuthority)           // 创建角色
			roleRouter.GET(":id/menus", v1.GetAuthorityMenus) // 获取角色菜单权限
			roleRouter.POST(":id/menus", v1.AssignMenus)      // 分配角色菜单权限
			roleRouter.GET(":id", v1.GetAuthorityById)        // 根据ID获取角色
			roleRouter.PUT(":id", v1.UpdateAuthority)         // 更新角色
			roleRouter.DELETE(":id", v1.DeleteAuthority)      // 删除角色
		}

		// 菜单管理路由
		menuRouter := systemRouter.Group("menu")
		{
			menuRouter.GET("list", v1.GetMenuList)  // 获取菜单列表
			menuRouter.GET("tree", v1.GetMenuTree)  // 获取菜单树结构
			menuRouter.POST("", v1.CreateMenu)      // 创建菜单
			menuRouter.GET(":id", v1.GetMenuById)   // 根据ID获取菜单
			menuRouter.PUT(":id", v1.UpdateMenu)    // 更新菜单
			menuRouter.DELETE(":id", v1.DeleteMenu) // 删除菜单
		}
	}
}
