package v1

import (
	"net/http"
	"server/global"
	"server/model/system"
	"server/utils"

	"github.com/gin-gonic/gin"
)

// LoginRequest 登录请求结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应结构体
type LoginResponse struct {
	User  system.SysUser `json:"user"`
	Token string         `json:"token"`
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 查找用户
	var user system.SysUser
	if err := global.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "用户名或密码错误",
		})
		return
	}

	// 验证密码
	if !checkPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "用户名或密码错误",
		})
		return
	}

	// 检查用户是否被禁用
	if user.Enable != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "用户已被禁用",
		})
		return
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, user.AuthorityId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "生成token失败",
		})
		return
	}

	// 清除密码字段
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "登录成功",
		"data": LoginResponse{
			User:  user,
			Token: token,
		},
	})
}

// checkPassword 验证密码
func checkPassword(password, hashedPassword string) bool {
	return utils.BcryptCheck(password, hashedPassword)
}

// Logout 用户登出
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "登出成功",
	})
}

// Captcha 获取验证码（暂时返回固定值）
func Captcha(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "获取成功",
		"data": gin.H{
			"captchaId":      "captcha_" + utils.RandomString(8),
			"pictureBases64": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==",
		},
	})
}
