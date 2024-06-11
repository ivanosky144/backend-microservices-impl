package models

import (
	"time"

	"gorm.io/gorm"
)

type UserVote struct {
	gorm.Model
	UserID    int       `gorm:"column:user_id"`
	CommentID int       `gorm:"column:comment_id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (u *UserVote) TableName() string {
	return "user_votes"
}
