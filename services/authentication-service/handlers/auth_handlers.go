package handlers

import (
	"context"
	"errors"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/temuka-authentication-service/config"
	"github.com/temuka-authentication-service/models"
	"github.com/temuka-authentication-service/pb"
	"golang.org/x/crypto/bcrypt"
)

func (s *server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	db := config.GetDBInstance()

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, err
	}

	newUser := models.User{
		Username:       req.Username,
		Email:          req.Email,
		Password:       string(hashedPwd),
		ProfilePicture: "",
		CoverPicture:   "",
	}

	if err := db.Create(&newUser).Error; err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		Message: "New user has been registered",
	}, nil
}

func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	db := config.GetDBInstance()

	var user models.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, err
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})

	tokenString, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Message: "User has login successfully",
		Token:   tokenString,
	}, nil
}

func (s *server) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	db := config.GetDBInstance()

	token, err := jwt.Parse(req.ResetToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		emailClaim, emailOk := claims["email"].(string)
		if !emailOk {
			return nil, errors.New("Invalid token: email claim not found or not a string")
		}

		// Fetch user from the database using the email from the claims
		var user models.User
		if err := db.Where("email = ?", emailClaim).First(&user).Error; err != nil {
			return nil, err
		}

		if req.Email != emailClaim {
			return nil, errors.New("Email in the request does not match the email in the token")
		}

		if req.NewPassword != req.NewPasswordConfirmation {
			return nil, errors.New("Password and password confirmation do not match")
		}

		hashedNewPwd, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 10)
		if err != nil {
			return nil, err
		}

		user.Password = string(hashedNewPwd)
		if err := db.Save(&user).Error; err != nil {
			return nil, err
		}

		return &pb.ResetPasswordResponse{
			Message: "Password was reset successfully",
		}, nil
	}

	return nil, errors.New("Invalid or expired token")
}
