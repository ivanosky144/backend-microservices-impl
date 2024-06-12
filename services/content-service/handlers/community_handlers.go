package handlers

import (
	"context"
	"errors"
	"log"

	"github.com/temuka-content-service/config"
	"github.com/temuka-content-service/models"
	pb "github.com/temuka-content-service/pb"
	"gorm.io/gorm"
)

func (s *server) CreateCommunity(ctx context.Context, req *pb.CreateCommunityRequest) (*pb.CreateCommunityResponse, error) {
	db := config.GetDBInstance()

	newCommunity := models.Community{
		Name:        req.GetName(),
		Description: req.GetDesc(),
		LogoPicture: req.GetLogopicture(),
	}

	if err := db.Create(&newCommunity).Error; err != nil {
		log.Printf("Error creating community: %v", err)
		return nil, err
	}

	return &pb.CreateCommunityResponse{
		Message: "Community has been created",
		Data: &pb.Community{
			Id:           int32(newCommunity.ID),
			Name:         newCommunity.Name,
			Desc:         newCommunity.Description,
			Logopicture:  newCommunity.LogoPicture,
			MembersCount: int32(newCommunity.MembersCount),
		},
	}, nil
}

func (s *server) JoinCommunity(ctx context.Context, req *pb.JoinCommunityRequest) (*pb.JoinCommunityResponse, error) {
	db := config.GetDBInstance()

	var existingMember models.CommunityMember
	if err := db.Where("community_id = ? AND user_id = ?", req.GetCommunityId(), req.GetUserId()).First(&existingMember).Error; err == nil {
		return nil, errors.New("User already a member of the community")
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	newMember := models.CommunityMember{
		UserID:      int(req.GetUserId()),
		CommunityID: int(req.GetCommunityId()),
	}

	if err := db.Create(&newMember).Error; err != nil {
		return nil, err
	}

	if err := db.Model(&models.Community{}).Where("id = ?", req.GetCommunityId()).Update("members_count", gorm.Expr("members_count + ?", 1)).Error; err != nil {
		return nil, err
	}

	return &pb.JoinCommunityResponse{
		Message: "Successfully joined the community",
	}, nil
}
