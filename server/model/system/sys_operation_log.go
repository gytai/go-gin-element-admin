package system

import (
	"time"

	"gorm.io/gorm"
)

// SysOperationLog 操作日志表
type SysOperationLog struct {
	gorm.Model
	UserID       uint      `json:"userId" gorm:"index;comment:用户ID"`
	Username     string    `json:"username" gorm:"comment:用户名"`
	Method       string    `json:"method" gorm:"comment:请求方法"`
	Path         string    `json:"path" gorm:"comment:请求路径"`
	OperationType string   `json:"operationType" gorm:"comment:操作类型(CREATE/UPDATE/DELETE)"`
	Description  string    `json:"description" gorm:"comment:操作描述"`
	RequestBody  string    `json:"requestBody" gorm:"type:text;comment:请求参数"`
	ResponseBody string    `json:"responseBody" gorm:"type:text;comment:响应结果"`
	IP           string    `json:"ip" gorm:"comment:请求IP"`
	UserAgent    string    `json:"userAgent" gorm:"comment:用户代理"`
	Status       int       `json:"status" gorm:"comment:响应状态码"`
	ErrorMessage string    `json:"errorMessage" gorm:"comment:错误信息"`
	Latency      int64     `json:"latency" gorm:"comment:请求耗时(毫秒)"`
	OperationTime time.Time `json:"operationTime" gorm:"comment:操作时间"`
}

func (SysOperationLog) TableName() string {
	return "sys_operation_logs"
}

// OperationLogRequest 操作日志查询请求
type OperationLogRequest struct {
	PageInfo
	UserID        uint   `json:"userId" form:"userId"`
	Username      string `json:"username" form:"username"`
	Method        string `json:"method" form:"method"`
	Path          string `json:"path" form:"path"`
	OperationType string `json:"operationType" form:"operationType"`
	Status        int    `json:"status" form:"status"`
	StartTime     string `json:"startTime" form:"startTime"`
	EndTime       string `json:"endTime" form:"endTime"`
}

// OperationLogResponse 操作日志响应
type OperationLogResponse struct {
	List     []SysOperationLog `json:"list"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"pageSize"`
}

// PageInfo 分页信息
type PageInfo struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
}
