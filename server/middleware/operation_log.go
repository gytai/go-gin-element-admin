package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"server/service/system"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// OperationLogMiddleware 操作日志中间件
func OperationLogMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 只记录需要记录的操作
		if shouldLogOperation(param.Method, param.Path) {
			go recordOperationLog(param)
		}
		return ""
	})
}

// ResponseBodyWriter 用于捕获响应体的写入器
type ResponseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r ResponseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// OperationLogMiddlewareWithBody 带请求体和响应体记录的操作日志中间件
func OperationLogMiddlewareWithBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只对需要记录的操作进行处理
		if !shouldLogOperation(c.Request.Method, c.Request.URL.Path) {
			c.Next()
			return
		}

		start := time.Now()

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 创建响应体写入器
		responseBodyWriter := &ResponseBodyWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = responseBodyWriter

		// 处理请求
		c.Next()

		// 计算耗时
		latency := time.Since(start).Milliseconds()

		// 异步记录操作日志
		go func() {
			recordDetailedOperationLog(c, requestBody, responseBodyWriter.body.Bytes(), latency)
		}()
	}
}

// shouldLogOperation 判断是否需要记录操作日志
func shouldLogOperation(method, path string) bool {
	// 只记录增删改操作
	if method != "POST" && method != "PUT" && method != "DELETE" {
		return false
	}

	// 排除不需要记录的路径
	excludePaths := []string{
		"/api/health",
		"/api/system/operation-log", // 避免记录日志查询操作本身
	}

	for _, excludePath := range excludePaths {
		if strings.HasPrefix(path, excludePath) {
			return false
		}
	}

	// 只记录API路径
	return strings.HasPrefix(path, "/api/")
}

// recordOperationLog 记录操作日志（简单版本）
func recordOperationLog(param gin.LogFormatterParams) {
	operationLogService := &system.OperationLogService{}

	// 确定操作类型
	operationType := getOperationType(param.Method)

	// 生成操作描述
	description := generateDescription(param.Method, param.Path)

	// 获取用户信息（这里需要从上下文中获取，暂时使用默认值）
	userID := uint(1)   // 默认管理员ID
	username := "admin" // 默认用户名

	operationLogService.LogOperation(
		userID,
		username,
		param.Method,
		param.Path,
		operationType,
		description,
		nil, // 请求体
		nil, // 响应体
		param.ClientIP,
		"", // UserAgent
		param.StatusCode,
		param.Latency.Milliseconds(),
		param.ErrorMessage,
	)
}

// recordDetailedOperationLog 记录详细操作日志
func recordDetailedOperationLog(c *gin.Context, requestBody, responseBody []byte, latency int64) {
	operationLogService := &system.OperationLogService{}

	// 确定操作类型
	operationType := getOperationType(c.Request.Method)

	// 生成操作描述
	description := generateDescription(c.Request.Method, c.Request.URL.Path)

	// 获取用户信息（从JWT token或session中获取）
	userID, username := getUserInfo(c)

	// 解析请求体和响应体
	var reqBodyInterface, respBodyInterface interface{}

	if len(requestBody) > 0 {
		json.Unmarshal(requestBody, &reqBodyInterface)
	}

	if len(responseBody) > 0 {
		json.Unmarshal(responseBody, &respBodyInterface)
	}

	// 获取错误信息
	errorMessage := ""
	if c.Writer.Status() >= 400 {
		if respBodyInterface != nil {
			if respMap, ok := respBodyInterface.(map[string]interface{}); ok {
				if msg, exists := respMap["msg"]; exists {
					errorMessage = msg.(string)
				}
			}
		}
	}

	operationLogService.LogOperation(
		userID,
		username,
		c.Request.Method,
		c.Request.URL.Path,
		operationType,
		description,
		reqBodyInterface,
		respBodyInterface,
		c.ClientIP(),
		c.Request.UserAgent(),
		c.Writer.Status(),
		latency,
		errorMessage,
	)
}

// getOperationType 根据HTTP方法确定操作类型
func getOperationType(method string) string {
	switch method {
	case "POST":
		return "CREATE"
	case "PUT":
		return "UPDATE"
	case "DELETE":
		return "DELETE"
	default:
		return "OTHER"
	}
}

// generateDescription 生成操作描述
func generateDescription(method, path string) string {
	operationType := getOperationType(method)

	// 根据路径生成描述
	if strings.Contains(path, "/user") {
		switch operationType {
		case "CREATE":
			return "创建用户"
		case "UPDATE":
			return "更新用户信息"
		case "DELETE":
			return "删除用户"
		}
	} else if strings.Contains(path, "/role") {
		switch operationType {
		case "CREATE":
			return "创建角色"
		case "UPDATE":
			return "更新角色信息"
		case "DELETE":
			return "删除角色"
		}
	} else if strings.Contains(path, "/menu") {
		switch operationType {
		case "CREATE":
			return "创建菜单"
		case "UPDATE":
			return "更新菜单信息"
		case "DELETE":
			return "删除菜单"
		}
	} else if strings.Contains(path, "/authority") {
		switch operationType {
		case "CREATE":
			return "创建权限"
		case "UPDATE":
			return "更新权限信息"
		case "DELETE":
			return "删除权限"
		}
	} else if strings.Contains(path, "/upload") {
		return "文件上传"
	} else if strings.Contains(path, "/login") {
		return "用户登录"
	} else if strings.Contains(path, "/logout") {
		return "用户登出"
	}

	// 根据具体路径生成更详细的描述
	switch {
	case strings.Contains(path, "/password"):
		return "修改密码"
	case strings.Contains(path, "/info"):
		return "更新个人信息"
	case strings.Contains(path, "/batch"):
		return "批量" + operationType
	case strings.Contains(path, "/assign"):
		return "分配权限"
	default:
		return operationType + " " + extractResourceName(path)
	}
}

// extractResourceName 从路径中提取资源名称
func extractResourceName(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) >= 3 {
		resource := parts[len(parts)-2]
		// 转换为中文描述
		switch resource {
		case "system":
			return "系统管理"
		case "dashboard":
			return "仪表板"
		case "operation-log":
			return "操作日志"
		default:
			return resource
		}
	}
	return path
}

// getUserInfo 从上下文中获取用户信息
func getUserInfo(c *gin.Context) (uint, string) {
	// 从JWT token中获取用户信息
	if userClaims, exists := c.Get("claims"); exists {
		if claims, ok := userClaims.(map[string]interface{}); ok {
			var userID uint
			var username string

			// 获取用户ID
			if id, exists := claims["userID"]; exists {
				if idFloat, ok := id.(float64); ok {
					userID = uint(idFloat)
				}
			}

			// 获取用户名
			if name, exists := claims["username"]; exists {
				if nameStr, ok := name.(string); ok {
					username = nameStr
				}
			}

			// 如果获取到了有效信息，返回
			if userID > 0 && username != "" {
				return userID, username
			}
		}
	}

	// 尝试从其他方式获取用户信息
	if userID, exists := c.Get("userID"); exists {
		if username, exists := c.Get("username"); exists {
			if id, ok := userID.(uint); ok {
				if name, ok := username.(string); ok {
					return id, name
				}
			}
		}
	}

	// 默认返回系统用户信息
	return 0, "system"
}
