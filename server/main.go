package main

import (
	"fmt"
	"server/core"
	"server/global"
	"server/initialize"
)

// @title Go-Gin-Element-Admin API
// @version 1.0
// @description 基于 Gin + Vue3 + Element-Admin 的后台管理系统
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @BasePath /
func main() {
	global.VP = core.Viper() // 初始化Viper
	initialize.OtherInit()
	global.DB = initialize.Gorm() // gorm连接数据库
	initialize.Timer()
	initialize.DBList()
	if global.DB != nil {
		initialize.RegisterTables() // 初始化表
		// 程序结束前关闭数据库链接
		db, _ := global.DB.DB()
		defer db.Close()
	}
	core.RunWindowsServer()
}

func init() {
	fmt.Println(`
	欢迎使用 go-gin-element-admin
	当前版本:v1.0.0
	默认自动化文档地址:http://127.0.0.1:8888/swagger/index.html
	默认前端文件运行地址:http://127.0.0.1:8080
	`)
}
