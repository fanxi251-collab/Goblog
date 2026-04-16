package service

import (
	"Goblog/internal/config"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// FileService 文件服务接口
type FileService interface {
	Upload(file *multipart.FileHeader) (string, error)
	Delete(path string) error
}

// LocalFileService 本地文件服务
type LocalFileService struct {
	basePath string
}

// NewLocalFileService 创建本地文件服务
func NewLocalFileService() *LocalFileService {
	cfg := config.Get()
	if cfg == nil {
		// 默认路径
		return &LocalFileService{basePath: "web/static/uploads"}
	}
	return &LocalFileService{basePath: cfg.Upload.Path}
}

// Upload 上传文件
func (s *LocalFileService) Upload(file *multipart.FileHeader) (string, error) {
	// 创建日期目录
	datePath := time.Now().Format("2006/01")
	uploadPath := filepath.Join(s.basePath, datePath)

	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", err
	}

	// 生成文件名
	ext := filepath.Ext(file.Filename)
	filename := time.Now().Format("20060102150405") + ext
	fullPath := filepath.Join(uploadPath, filename)

	// 打开并保存文件
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	// 返回相对路径
	return filepath.Join(datePath, filename), nil
}

// Delete 删除文件
func (s *LocalFileService) Delete(path string) error {
	fullPath := filepath.Join(s.basePath, path)
	return os.Remove(fullPath)
}
