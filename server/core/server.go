package core

import (
	"fmt"
	"net/http"
	"time"

	"server/global"
	"server/initialize"

	"github.com/gin-gonic/gin"
)

func initServer(address string, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:           address,
		Handler:        router,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func RunWindowsServer() {
	if global.CONFIG.System.UseMultipoint || global.CONFIG.System.UseRedis {
		// 初始化redis服务
		initialize.Redis()
	}

	Router := initialize.Routers()
	Router.Static("/form-generator", "./resource/page")

	address := fmt.Sprintf(":%d", global.CONFIG.System.Addr)
	s := initServer(address, Router)
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	fmt.Printf(`
	欢迎使用 go-gin-element-admin
	当前版本:v1.0.0
	默认自动化文档地址:http://127.0.0.1%s/swagger/index.html
	默认前端文件运行地址:http://127.0.0.1:8080
	`, address)
	fmt.Printf(`
	项目运行中...
	`)
	fmt.Println(s.ListenAndServe().Error())
}
