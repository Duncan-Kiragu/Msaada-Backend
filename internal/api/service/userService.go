package service

import (
	"context"

	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/domain"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/dto"
	"github.com/Duncan-Kiragu/Msaada-Backend/pkg/filter"
)

func NewUserService(r domain.UserRepository) domain.UserService {
	return &userService{
		userRepository: r,
	}
}

type userService struct {
	userRepository domain.UserRepository
}

func (s *userService) generateUserOutputDTO(user *domain.User) *dto.UserOutputDTO {
	return &dto.UserOutputDTO{
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Status: user.Status,
		Profile: dto.ProfileOutputDTO{
			Id:   user.ProfileID,
			Name: user.Profile.Name,
		},
	}
}

// GetUserByID Implementation of 'GetUserByID'.
func (s *userService) GetUserByID(ctx context.Context, userID uint) (*dto.UserOutputDTO, error) {
	user, err := s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.generateUserOutputDTO(user), nil
}

// GetUsers Implementation of 'GetUsers'.
func (s *userService) GetUsers(ctx context.Context, filter *filter.UserFilter) (*dto.ListItemsOutputDTO, error) {
	count, err := s.userRepository.CountUsers(ctx, filter)
	if err != nil {
		return nil, err
	}

	users, err := s.userRepository.GetUsers(ctx, filter)
	if err != nil {
		return nil, err
	}

	outputUsers := &[]dto.UserOutputDTO{}
	for _, user := range *users {
		*outputUsers = append(*outputUsers, *s.generateUserOutputDTO(&user))
	}

	return &dto.ListItemsOutputDTO{
		Count: count,
		Items: outputUsers,
	}, nil
}

// GetUserByMail Implementation of 'GetUserByMail'.
func (s *userService) GetUserByMail(ctx context.Context, userMail string) (*domain.User, error) {
	return s.userRepository.GetUserByMail(ctx, userMail)
}

// GetUserByToken Implementation of 'GetUserByToken'.
func (s *userService) GetUserByToken(ctx context.Context, token string) (*domain.User, error) {
	return s.userRepository.GetUserByToken(ctx, token)
}

// CreateUser Implementation of 'CreateUser'.
func (s *userService) CreateUser(ctx context.Context, data *dto.UserInputDTO) (*dto.UserOutputDTO, error) {
	user, err := s.userRepository.CreateUser(ctx, data)
	if err != nil {
		return nil, err
	}

	user, err = s.userRepository.GetUserByID(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	return s.generateUserOutputDTO(user), nil
}

// UpdateUser Implementation of 'UpdateUser'.
func (s *userService) UpdateUser(ctx context.Context, user *domain.User, data *dto.UserInputDTO) (*dto.UserOutputDTO, error) {
	if err := s.userRepository.UpdateUser(ctx, user, data); err != nil {
		return nil, err
	}

	user, err := s.userRepository.GetUserByID(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	return s.generateUserOutputDTO(user), nil
}

// DeleteUser Implementation of 'DeleteUser'.
func (s *userService) DeleteUser(ctx context.Context, user *domain.User) error {
	return s.userRepository.DeleteUser(ctx, user)
}

// ResetUserPassword Implementation of 'ResetUserPassword'.
func (s *userService) ResetUserPassword(ctx context.Context, user *domain.User) error {
	return s.userRepository.ResetUserPassword(ctx, user)
}

// SetUserPassword Implementation of 'SetUserPassword'.
func (s *userService) SetUserPassword(ctx context.Context, user *domain.User, pass *dto.PasswordInputDTO) error {
	return s.userRepository.SetUserPassword(ctx, user, pass)
}
