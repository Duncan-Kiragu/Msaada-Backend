package service

import (
	"context"
	"errors"
	"os"

	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/dto"

	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/domain"
)

func NewAuthService(r domain.UserRepository) domain.AuthService {
	return &authService{
		userRepository: r,
	}
}

type authService struct {
	userRepository domain.UserRepository
}

func (s *authService) generateUserOutputDTO(user *domain.User) *dto.UserOutputDTO {
	if user == nil {
		return nil
	}

	return &dto.UserOutputDTO{
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Status: user.Status,
		Profile: dto.ProfileOutputDTO{
			Id:   user.ProfileID,
			Name: user.Profile.Name,
			Permissions: &dto.PermissionsOutputDTO{
				UserModule:    user.Profile.Permissions.UserModule,
				ProfileModule: user.Profile.Permissions.ProfileModule,
				ProductModule: user.Profile.Permissions.ProductModule,
			},
		},
	}
}

func (s *authService) generateAuthOutputDTO(user *domain.User, ip string) *dto.AuthOutputDTO {
	accessTime, refreshTime := "-", "-"
	if user.Expire {
		accessTime = os.Getenv("ACCESS_TOKEN_EXPIRE")
		refreshTime = os.Getenv("RFRESH_TOKEN_EXPIRE")
	}

	accessToken, _ := user.GenerateToken(accessTime, os.Getenv("ACCESS_TOKEN_PRIVAT"), ip)
	refreshToken, _ := user.GenerateToken(refreshTime, os.Getenv("RFRESH_TOKEN_PRIVAT"), ip)

	return &dto.AuthOutputDTO{
		User:         s.generateUserOutputDTO(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func (s *authService) Login(ctx context.Context, credentials *dto.AuthInputDTO, ip string) (*dto.AuthOutputDTO, error) {
	user, err := s.userRepository.GetUserByMail(ctx, credentials.Login)
	if err != nil {
		return nil, err
	}

	if !user.ValidatePassword(credentials.Password) {
		return nil, errors.New("invalid password")
	}

	if !user.Status || user.New {
		return nil, errors.New("invalid user")
	}

	user.Expire = credentials.Expire
	return s.generateAuthOutputDTO(user, ip), nil
}

func (s *authService) Me(user *domain.User) *dto.UserOutputDTO {
	return s.generateUserOutputDTO(user)
}

func (s *authService) Refresh(user *domain.User, ip string) *dto.AuthOutputDTO {
	return s.generateAuthOutputDTO(user, ip)
}
