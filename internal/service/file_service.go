package service

import (
	"Goblog/internal/config"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nfnt/resize"
)

// FileService 文件服务接口
type FileService interface {
	Upload(file *multipart.FileHeader) (string, error)
	UploadCover(file *multipart.FileHeader, postID uint) (string, error)
	Delete(path string) error
	DeleteCover(postID uint) error
}

// 封面尺寸
const (
	CoverWidth  = 200
	CoverHeight = 200
)

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

	// 返回相对路径（将反斜杠替换为正斜杠，兼容 URL）
	path := filepath.Join(datePath, filename)
	path = strings.ReplaceAll(path, "\\", "/")
	return path, nil
}

// Delete 删除文件
func (s *LocalFileService) Delete(path string) error {
	fullPath := filepath.Join(s.basePath, path)
	return os.Remove(fullPath)
}

// UploadCover 上传并裁剪封面图片
func (s *LocalFileService) UploadCover(file *multipart.FileHeader, postID uint) (string, error) {
	// 创建封面目录
	coverPath := "web/static/covers"
	if err := os.MkdirAll(coverPath, 0755); err != nil {
		return "", err
	}

	// 打开上传的图片
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 解码图片
	img, format, err := image.Decode(src)
	if err != nil {
		return "", err
	}

	// 缩放到目标尺寸
	thumbnail := resize.Resize(CoverWidth, CoverHeight, img, resize.Lanczos3)

	// 生成文件名
	ext := ".jpg"
	if format == "png" {
		ext = ".png"
	}
	filename := "cover-" + formatPostID(postID) + ext
	fullPath := filepath.Join(coverPath, filename)

	// 保存裁剪后的图片
	var dst *os.File
	if ext == ".png" {
		dst, err = os.Create(fullPath)
		if err != nil {
			return "", err
		}
		defer dst.Close()
		err = png.Encode(dst, thumbnail)
	} else {
		dst, err = os.Create(fullPath)
		if err != nil {
			return "", err
		}
		defer dst.Close()
		err = jpeg.Encode(dst, thumbnail, &jpeg.Options{Quality: 85})
	}

	if err != nil {
		return "", err
	}

	// 返回相对于 static 的路径
	return "covers/" + filename, nil
}

// formatPostID 将 postID 格式化为 3 位数字字符串
func formatPostID(id uint) string {
	return fmt.Sprintf("%03d", id)
}

// DeleteCover 删除封面图片
func (s *LocalFileService) DeleteCover(postID uint) error {
	coverPath := "web/static/covers"
	// 尝试删除 jpg 和 png
	jpgPath := filepath.Join(coverPath, "cover-"+formatPostID(postID)+".jpg")
	pngPath := filepath.Join(coverPath, "cover-"+formatPostID(postID)+".png")

	if _, err := os.Stat(jpgPath); err == nil {
		return os.Remove(jpgPath)
	}
	if _, err := os.Stat(pngPath); err == nil {
		return os.Remove(pngPath)
	}
	return nil
}
