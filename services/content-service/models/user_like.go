package models

import (
	"time"

	"gorm.io/gorm"
)

type UserLike struct {
	gorm.Model
	UserID    int       `gorm:"column:user_id"`
	PostID    int       `gorm:"column:post_id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (u *UserLike) TableName() string {
	return "user_likes"
}
