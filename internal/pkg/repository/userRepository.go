package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/domain"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/dto"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/postgre"
	"github.com/Duncan-Kiragu/Msaada-Backend/pkg/filter"
)

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

type userRepository struct {
	db *gorm.DB
}

func (s *userRepository) applyFilter(ctx context.Context, filter *filter.UserFilter) *gorm.DB {
	db := s.db.WithContext(ctx)
	if filter.ProfileID != 0 {
		db = db.Where(domain.UserTableName+".profile_id = ?", filter.ProfileID)
	}
	db = db.Joins(fmt.Sprintf("JOIN %v ON %v.id = %v.profile_id", domain.ProfileTableName, domain.ProfileTableName, domain.UserTableName))
	db = filter.ApplySearchLike(db, domain.UserTableName+".name", domain.UserTableName+".mail", domain.ProfileTableName+".name")

	return filter.ApplyOrder(db)
}

func (s *userRepository) CountUsers(ctx context.Context, filter *filter.UserFilter) (int64, error) {
	var count int64
	db := s.applyFilter(ctx, filter)
	return count, db.Model(&domain.User{}).Count(&count).Error
}

func (s *userRepository) GetUsers(ctx context.Context, filter *filter.UserFilter) (*[]domain.User, error) {
	db := filter.ApplyPagination(s.applyFilter(ctx, filter))

	users := &[]domain.User{}
	return users, db.Preload(postgre.ProfilePermission).Find(users).Error
}

func (s *userRepository) GetUserByID(ctx context.Context, userID uint) (*domain.User, error) {
	user := &domain.User{}
	return user, s.db.WithContext(ctx).Preload(postgre.ProfilePermission).First(user, userID).Error
}

func (s *userRepository) GetUserByMail(ctx context.Context, mail string) (*domain.User, error) {
	user := &domain.User{Email: mail}
	return user, s.db.WithContext(ctx).Preload(postgre.ProfilePermission).Where(user).First(user).Error
}

func (s *userRepository) GetUserByToken(ctx context.Context, token string) (*domain.User, error) {
	user := &domain.User{Token: &token}
	return user, s.db.WithContext(ctx).Preload(postgre.ProfilePermission).Where(user).First(user).Error
}

func (s *userRepository) CreateUser(ctx context.Context, data *dto.UserInputDTO) (*domain.User, error) {
	user := &domain.User{New: true}
	if err := user.Bind(data); err != nil {
		return nil, err
	}

	return user, s.db.WithContext(ctx).Create(user).Error
}

func (s *userRepository) UpdateUser(ctx context.Context, user *domain.User, data *dto.UserInputDTO) error {
	if err := user.Bind(data); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Model(user).Updates(user.ToMap()).Error
}

func (s *userRepository) DeleteUser(ctx context.Context, user *domain.User) error {
	return s.db.WithContext(ctx).Delete(user).Error
}

func (s *userRepository) ResetUserPassword(ctx context.Context, user *domain.User) error {
	user.Password = nil
	user.Token = nil
	user.New = true

	return s.UpdateUser(ctx, user, &dto.UserInputDTO{})
}

func (s *userRepository) SetUserPassword(ctx context.Context, user *domain.User, pass *dto.PasswordInputDTO) error {
	user.New = false
	user.Token = new(string)
	*user.Token = uuid.New().String()

	hash, err := bcrypt.GenerateFromPassword([]byte(*pass.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = new(string)
	*user.Password = string(hash)

	return s.UpdateUser(ctx, user, &dto.UserInputDTO{})
}
