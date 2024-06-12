package handlers

import (
	"context"
	"net/http"

	"github.com/temuka-authentication-service/config"
	"github.com/temuka-authentication-service/models"
	"github.com/temuka-authentication-service/pb"
)

func (s *server) SearchUsers(ctx context.Context, req *pb.SearchUsersRequest) (*pb.SearchUsersResponse, error) {
	db := config.GetDBInstance()

	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, httpError(http.StatusNotFound, "No user found")
	}

	var pbUsers []*pb.User
	for _, user := range users {
		pbUsers = append(pbUsers, &pb.User{
			Id:             int32(user.ID),
			Username:       user.Username,
			Email:          user.Email,
			ProfilePicture: user.ProfilePicture,
			Desc:           user.Desc,
		})
	}

	return &pb.SearchUsersResponse{
		Message: "Users have been retrieved",
		Data:    pbUsers,
	}, nil
}

func (s *server) GetUserDetail(ctx context.Context, req *pb.UserDetailRequest) (*pb.UserDetailResponse, error) {
	db := config.GetDBInstance()

	userID := req.Id
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, httpError(http.StatusNotFound, "User not found")
	}

	return &pb.UserDetailResponse{
		Id:             int32(user.ID),
		Username:       user.Username,
		Email:          user.Email,
		ProfilePicture: user.ProfilePicture,
		Desc:           user.Desc,
	}, nil
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	db := config.GetDBInstance()

	newUser := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := db.Create(&newUser).Error; err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{
		Message: "User has been created",
		Data: &pb.User{
			Id:             int32(newUser.ID),
			Username:       newUser.Username,
			Email:          newUser.Email,
			ProfilePicture: newUser.ProfilePicture,
			Desc:           newUser.Desc,
		},
	}, nil
}

func (s *server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	db := config.GetDBInstance()

	userID := req.Id
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, httpError(http.StatusNotFound, "User not found")
	}

	user.Username = req.Username
	user.Email = req.Email
	user.ProfilePicture = req.ProfilePicture
	user.Desc = req.Desc

	if err := db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &pb.UpdateUserResponse{
		Message: "User has been updated",
		Data: &pb.User{
			Id:             int32(user.ID),
			Username:       user.Username,
			Email:          user.Email,
			ProfilePicture: user.ProfilePicture,
			Desc:           user.Desc,
		},
	}, nil
}

func httpError(statusCode int, message string) error {
	return &err{statusCode, message}
}

type err struct {
	statusCode int
	message    string
}

func (e *err) Error() string {
	return e.message
}
