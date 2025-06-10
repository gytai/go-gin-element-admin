package system

import (
	"encoding/json"
	"fmt"
	"server/global"
	"server/model/system"
	"time"
)

type OperationLogService struct{}

// CreateOperationLog 创建操作日志
func (s *OperationLogService) CreateOperationLog(log *system.SysOperationLog) error {
	return global.DB.Create(log).Error
}

// GetOperationLogList 获取操作日志列表
func (s *OperationLogService) GetOperationLogList(req system.OperationLogRequest) (system.OperationLogResponse, error) {
	var logs []system.SysOperationLog
	var total int64

	db := global.DB.Model(&system.SysOperationLog{})

	// 构建查询条件
	if req.UserID != 0 {
		db = db.Where("user_id = ?", req.UserID)
	}
	if req.Username != "" {
		db = db.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Method != "" {
		db = db.Where("method = ?", req.Method)
	}
	if req.Path != "" {
		db = db.Where("path LIKE ?", "%"+req.Path+"%")
	}
	if req.OperationType != "" {
		db = db.Where("operation_type = ?", req.OperationType)
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}
	if req.StartTime != "" {
		db = db.Where("operation_time >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		db = db.Where("operation_time <= ?", req.EndTime)
	}

	// 获取总数
	err := db.Count(&total).Error
	if err != nil {
		return system.OperationLogResponse{}, err
	}

	// 分页查询
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	offset := (req.Page - 1) * req.PageSize
	err = db.Order("operation_time DESC").Offset(offset).Limit(req.PageSize).Find(&logs).Error
	if err != nil {
		return system.OperationLogResponse{}, err
	}

	return system.OperationLogResponse{
		List:     logs,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// DeleteOperationLog 删除操作日志
func (s *OperationLogService) DeleteOperationLog(id uint) error {
	return global.DB.Delete(&system.SysOperationLog{}, id).Error
}

// DeleteOperationLogsByIds 批量删除操作日志
func (s *OperationLogService) DeleteOperationLogsByIds(ids []uint) error {
	return global.DB.Delete(&system.SysOperationLog{}, ids).Error
}

// ClearOperationLogs 清空操作日志
func (s *OperationLogService) ClearOperationLogs() error {
	return global.DB.Exec("TRUNCATE TABLE sys_operation_logs").Error
}

// ClearOperationLogsByDays 清理指定天数前的操作日志
func (s *OperationLogService) ClearOperationLogsByDays(days int) error {
	cutoffTime := time.Now().AddDate(0, 0, -days)
	return global.DB.Where("operation_time < ?", cutoffTime).Delete(&system.SysOperationLog{}).Error
}

// GetOperationLogById 根据ID获取操作日志
func (s *OperationLogService) GetOperationLogById(id uint) (system.SysOperationLog, error) {
	var log system.SysOperationLog
	err := global.DB.Where("id = ?", id).First(&log).Error
	return log, err
}

// LogOperation 记录操作日志的便捷方法
func (s *OperationLogService) LogOperation(userID uint, username, method, path, operationType, description string, requestBody, responseBody interface{}, ip, userAgent string, status int, latency int64, errorMsg string) {
	// 序列化请求和响应数据
	var reqBodyStr, respBodyStr string

	if requestBody != nil {
		if reqBytes, err := json.Marshal(requestBody); err == nil {
			reqBodyStr = string(reqBytes)
		}
	}

	if responseBody != nil {
		if respBytes, err := json.Marshal(responseBody); err == nil {
			respBodyStr = string(respBytes)
		}
	}

	// 限制字段长度，避免数据过大
	if len(reqBodyStr) > 5000 {
		reqBodyStr = reqBodyStr[:5000] + "...[truncated]"
	}
	if len(respBodyStr) > 5000 {
		respBodyStr = respBodyStr[:5000] + "...[truncated]"
	}

	log := &system.SysOperationLog{
		UserID:        userID,
		Username:      username,
		Method:        method,
		Path:          path,
		OperationType: operationType,
		Description:   description,
		RequestBody:   reqBodyStr,
		ResponseBody:  respBodyStr,
		IP:            ip,
		UserAgent:     userAgent,
		Status:        status,
		ErrorMessage:  errorMsg,
		Latency:       latency,
		OperationTime: time.Now(),
	}

	// 异步记录日志，避免影响主业务
	go func() {
		if err := s.CreateOperationLog(log); err != nil {
			fmt.Printf("记录操作日志失败: %v\n", err)
		}
	}()
}

// GetOperationStats 获取操作统计信息
func (s *OperationLogService) GetOperationStats() (map[string]interface{}, error) {
	var stats map[string]interface{} = make(map[string]interface{})

	// 今日操作数
	today := time.Now().Format("2006-01-02")
	var todayCount int64
	err := global.DB.Model(&system.SysOperationLog{}).
		Where("DATE(operation_time) = ?", today).
		Count(&todayCount).Error
	if err != nil {
		return nil, err
	}
	stats["todayCount"] = todayCount

	// 总操作数
	var totalCount int64
	err = global.DB.Model(&system.SysOperationLog{}).Count(&totalCount).Error
	if err != nil {
		return nil, err
	}
	stats["totalCount"] = totalCount

	// 操作类型统计
	var typeStats []struct {
		OperationType string `json:"operationType"`
		Count         int64  `json:"count"`
	}
	err = global.DB.Model(&system.SysOperationLog{}).
		Select("operation_type, COUNT(*) as count").
		Group("operation_type").
		Find(&typeStats).Error
	if err != nil {
		return nil, err
	}
	stats["typeStats"] = typeStats

	// 最近7天操作趋势
	var dailyStats []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	err = global.DB.Model(&system.SysOperationLog{}).
		Select("DATE(operation_time) as date, COUNT(*) as count").
		Where("operation_time >= ?", time.Now().AddDate(0, 0, -7)).
		Group("DATE(operation_time)").
		Order("date").
		Find(&dailyStats).Error
	if err != nil {
		return nil, err
	}
	stats["dailyStats"] = dailyStats

	return stats, nil
}
