package handlers

import (
	"context"

	"github.com/temuka-content-service/config"
	"github.com/temuka-content-service/models"
	"github.com/temuka-content-service/pb"
)

func (s *server) AddComment(ctx context.Context, req *pb.AddCommentRequest) (*pb.AddCommentResponse, error) {
	db := config.GetDBInstance()

	newComment := models.Comment{
		UserID:  int(req.UserId),
		PostID:  int(req.PostId),
		Content: req.Content,
	}

	if err := db.Create(&newComment).Error; err != nil {
		return nil, err
	}

	response := &pb.AddCommentResponse{
		Message: "Comment has been added",
		Data: &pb.Comment{
			Id:      int32(newComment.ID),
			PostId:  int32(newComment.PostID),
			UserId:  int32(newComment.UserID),
			Content: newComment.Content,
		},
	}

	return response, nil
}
