package model

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// Visitor 访客模型
type Visitor struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	Token     string `gorm:"uniqueIndex;size:64" json:"token"`
	Nickname  string `gorm:"size:50" json:"nickname"`
	Email     string `gorm:"size:100" json:"email"`
	IP        string `gorm:"size:50" json:"ip"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
}

// TableName 表名
func (Visitor) TableName() string {
	return "visitors"
}

// GenerateToken 生成唯一Token
func GenerateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// IsValid 检查Token是否有效（30天有效期内）
func (v *Visitor) IsValid() bool {
	if v == nil || v.Token == "" {
		return false
	}
	// 30天有效期
	now := time.Now().Unix()
	return (now - v.CreatedAt) < 30*24*60*60
}
