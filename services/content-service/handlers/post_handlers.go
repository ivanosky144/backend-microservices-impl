package handlers

import (
	"context"

	"github.com/temuka-content-service/config"
	"github.com/temuka-content-service/models"
	pb "github.com/temuka-content-service/pb"
)

func (s *server) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {
	db := config.GetDBInstance()

	newPost := models.Post{
		Title:       req.GetTitle(),
		Description: req.GetDesc(),
		UserID:      int(req.GetUserId()),
	}
	db.Create(&newPost)

	return &pb.CreatePostResponse{
		Message: "Post has been created",
		Data: &pb.Post{
			Id:     int32(newPost.ID),
			Title:  newPost.Title,
			Desc:   newPost.Description,
			UserId: int32(newPost.UserID),
		},
	}, nil
}

func (s *server) GetTimelinePosts(ctx context.Context, req *pb.GetTimelinePostsRequest) (*pb.GetTimelinePostsResponse, error) {
	db := config.GetDBInstance()

	var timelinePosts []models.Post
	if err := db.Where("user_id = ?", req.GetUserId()).Find(&timelinePosts).Error; err != nil {
		return nil, err
	}

	var pbPosts []*pb.Post
	for _, post := range timelinePosts {
		pbPosts = append(pbPosts, &pb.Post{
			Id:     int32(post.ID),
			Title:  post.Title,
			Desc:   post.Description,
			UserId: int32(post.UserID),
		})
	}

	return &pb.GetTimelinePostsResponse{
		Message: "Timeline posts has been retrieved",
		Data:    pbPosts,
	}, nil
}

// DeletePost implements the DeletePost gRPC method.
func (s *server) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	db := config.GetDBInstance()

	if err := db.Delete(&models.Post{}, req.GetId()).Error; err != nil {
		return nil, err
	}

	return &pb.DeletePostResponse{
		Message: "Post has been deleted",
	}, nil
}

func (s *server) LikePost(ctx context.Context, req *pb.LikePostRequest) (*pb.LikePostResponse, error) {
	db := config.GetDBInstance()

	var post models.Post
	if err := db.First(&post, req.GetId()).Error; err != nil {
		return nil, err
	}

	alreadyLiked := false
	for _, user := range post.Likes {
		if uint(user.ID) == uint(req.GetUserId()) {
			alreadyLiked = true
			break
		}
	}

	if !alreadyLiked {
		var liker models.User
		if err := db.First(&liker, req.GetUserId()).Error; err != nil {
			return nil, err
		}

		// Append the liker to the post's likes slice
		post.Likes = append(post.Likes, &liker)

		// Save the updated post to the database
		if err := db.Save(&post).Error; err != nil {
			return nil, err
		}

		return &pb.LikePostResponse{
			Message: "You have liked this post",
		}, nil
	}

	return &pb.LikePostResponse{
		Message: "You have already liked this post",
	}, nil
}
