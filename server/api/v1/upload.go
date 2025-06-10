package v1

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	MaxFileSize = 10 * 1024 * 1024 // 10MB
	UploadPath  = "./uploads/avatar/"
)

// UploadAvatar 上传头像
func UploadAvatar(c *gin.Context) {
	// 限制文件大小为10MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxFileSize)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		if err.Error() == "http: request body too large" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "文件大小不能超过10MB",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请选择要上传的文件",
		})
		return
	}
	defer file.Close()

	// 验证文件类型
	if !isValidImageType(header) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "只支持上传 JPG、JPEG、PNG、GIF 格式的图片",
		})
		return
	}

	// 验证文件大小
	if header.Size > MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "文件大小不能超过10MB",
		})
		return
	}

	// 创建上传目录
	if err := os.MkdirAll(UploadPath, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建上传目录失败",
		})
		return
	}

	// 生成新的文件名
	ext := filepath.Ext(header.Filename)
	newFileName := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	filePath := filepath.Join(UploadPath, newFileName)

	// 保存文件
	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "保存文件失败",
		})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "保存文件失败",
		})
		return
	}

	// 生成访问URL
	fileURL := fmt.Sprintf("/uploads/avatar/%s", newFileName)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "上传成功",
		"data": gin.H{
			"url":      fileURL,
			"filename": newFileName,
		},
	})
}

// isValidImageType 验证是否为有效的图片类型
func isValidImageType(header *multipart.FileHeader) bool {
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
	}

	// 检查Content-Type
	contentType := header.Header.Get("Content-Type")
	if allowedTypes[contentType] {
		return true
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}

	return allowedExts[ext]
}
