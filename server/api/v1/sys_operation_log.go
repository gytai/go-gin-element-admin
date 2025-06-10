package v1

import (
	"fmt"
	"net/http"
	system2 "server/model/system"
	"server/service/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

var operationLogService = system.OperationLogService{}

// GetOperationLogList 获取操作日志列表
// @Tags      操作日志
// @Summary   获取操作日志列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     system.OperationLogRequest  true  "查询参数"
// @Success   200   {object}  response.Response{data=system.OperationLogResponse,msg=string}  "获取成功"
// @Router    /system/operation-log/list [get]
func GetOperationLogList(c *gin.Context) {
	var req system2.OperationLogRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	list, err := operationLogService.GetOperationLogList(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "获取成功",
		"data": gin.H{
			"list":     list.List,
			"total":    list.Total,
			"page":     list.Page,
			"pageSize": list.PageSize,
		},
	})
}

// GetOperationLogById 根据ID获取操作日志详情
// @Tags      操作日志
// @Summary   根据ID获取操作日志详情
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     id   path      int                         true  "操作日志ID"
// @Success   200  {object}  response.Response{data=system.SysOperationLog,msg=string}  "获取成功"
// @Router    /system/operation-log/{id} [get]
func GetOperationLogById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "ID格式错误",
		})
		return
	}

	log, err := operationLogService.GetOperationLogById(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "操作日志不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "获取成功",
		"data": log,
	})
}

// DeleteOperationLog 删除操作日志
// @Tags      操作日志
// @Summary   删除操作日志
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     id   path      int             true  "操作日志ID"
// @Success   200  {object}  response.Response{msg=string}  "删除成功"
// @Router    /system/operation-log/{id} [delete]
func DeleteOperationLog(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "ID格式错误",
		})
		return
	}

	err = operationLogService.DeleteOperationLog(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "删除失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "删除成功",
	})
}

// DeleteOperationLogsByIds 批量删除操作日志
// @Tags      操作日志
// @Summary   批量删除操作日志
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.IdsReq  true  "操作日志ID列表"
// @Success   200   {object}  response.Response{msg=string}  "删除成功"
// @Router    /system/operation-log/batch [delete]
func DeleteOperationLogsByIds(c *gin.Context) {
	var req struct {
		Ids []uint `json:"ids" binding:"required"`
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	err = operationLogService.DeleteOperationLogsByIds(req.Ids)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "删除失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "删除成功",
	})
}

// ClearOperationLogs 清空所有操作日志
// @Tags      操作日志
// @Summary   清空所有操作日志
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{msg=string}  "清空成功"
// @Router    /system/operation-log/clear [delete]
func ClearOperationLogs(c *gin.Context) {
	err := operationLogService.ClearOperationLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "清空失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "清空成功",
	})
}

// ClearOperationLogsByDays 清理指定天数前的操作日志
// @Tags      操作日志
// @Summary   清理指定天数前的操作日志
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.DaysReq  true  "保留天数"
// @Success   200   {object}  response.Response{msg=string}  "清理成功"
// @Router    /system/operation-log/clear-by-days [delete]
func ClearOperationLogsByDays(c *gin.Context) {
	var req struct {
		Days int `json:"days" binding:"required,min=1"`
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	err = operationLogService.ClearOperationLogsByDays(req.Days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "清理失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "清理成功",
	})
}

// GetOperationStats 获取操作统计信息
// @Tags      操作日志
// @Summary   获取操作统计信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=map[string]interface{},msg=string}  "获取成功"
// @Router    /system/operation-log/stats [get]
func GetOperationStats(c *gin.Context) {
	stats, err := operationLogService.GetOperationStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "获取成功",
		"data": stats,
	})
}

// ExportOperationLogs 导出操作日志
// @Tags      操作日志
// @Summary   导出操作日志
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     system.OperationLogRequest  true  "查询参数"
// @Success   200   {object}  response.Response{msg=string}  "导出成功"
// @Router    /system/operation-log/export [get]
func ExportOperationLogs(c *gin.Context) {
	var req system2.OperationLogRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 设置大的页面大小以获取所有数据
	req.PageSize = 10000
	req.Page = 1

	list, err := operationLogService.GetOperationLogList(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "导出失败: " + err.Error(),
		})
		return
	}

	// 设置响应头为CSV文件
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=operation_logs.csv")

	// 写入CSV头部
	c.Writer.WriteString("\xEF\xBB\xBF") // UTF-8 BOM
	c.Writer.WriteString("ID,用户名,操作类型,请求方法,请求路径,操作描述,状态码,IP地址,耗时(ms),操作时间,错误信息\n")

	// 写入数据
	for _, log := range list.List {
		c.Writer.WriteString(fmt.Sprintf("%d,%s,%s,%s,%s,%s,%d,%s,%d,%s,%s\n",
			log.ID,
			log.Username,
			log.OperationType,
			log.Method,
			log.Path,
			log.Description,
			log.Status,
			log.IP,
			log.Latency,
			log.OperationTime.Format("2006-01-02 15:04:05"),
			log.ErrorMessage,
		))
	}

	c.Status(http.StatusOK)
}
