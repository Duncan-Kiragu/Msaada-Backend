package domain

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/dto"
	"github.com/Duncan-Kiragu/Msaada-Backend/pkg/filter"
	"github.com/Duncan-Kiragu/Msaada-Backend/pkg/helpers"
	"github.com/Duncan-Kiragu/Msaada-Backend/pkg/validator"
)

const UserTableName string = "users"

type (
	User struct {
		Base
		Name      string   `json:"name" gorm:"column:name;type:varchar(90);not null;" validate:"required,min=5"`
		Email     string   `json:"mail" gorm:"column:mail;type:varchar(50);not null;unique;index;" validate:"required,email"`
		Status    bool     `json:"status" gorm:"column:status;type:bool;not null;"`
		New       bool     `json:"new" gorm:"column:new;type:bool;not null;"`
		ProfileID uint     `json:"-" gorm:"column:profile_id;type:bigint;not null;index;" validate:"required,min=1"`
		Token     *string  `json:"-" gorm:"column:token;type:varchar(255);unique;index"`
		Password  *string  `json:"-" gorm:"column:password;type:varchar(255);"`
		Profile   *Profile `json:"profile,omitempty"`
		Expire    bool     `json:"-" gorm:"-"`
	}

	UserRepository interface {
		CountUsers(context.Context, *filter.UserFilter) (int64, error)
		GetUserByID(context.Context, uint) (*User, error)
		GetUsers(context.Context, *filter.UserFilter) (*[]User, error)
		GetUserByMail(context.Context, string) (*User, error)
		GetUserByToken(context.Context, string) (*User, error)
		CreateUser(context.Context, *dto.UserInputDTO) (*User, error)
		UpdateUser(context.Context, *User, *dto.UserInputDTO) error
		DeleteUser(context.Context, *User) error
		ResetUserPassword(context.Context, *User) error
		SetUserPassword(context.Context, *User, *dto.PasswordInputDTO) error
	}

	UserService interface {
		GetUserByID(context.Context, uint) (*dto.UserOutputDTO, error)
		GetUsers(context.Context, *filter.UserFilter) (*dto.ListItemsOutputDTO, error)
		GetUserByMail(context.Context, string) (*User, error)
		GetUserByToken(context.Context, string) (*User, error)
		CreateUser(context.Context, *dto.UserInputDTO) (*dto.UserOutputDTO, error)
		UpdateUser(context.Context, *User, *dto.UserInputDTO) (*dto.UserOutputDTO, error)
		DeleteUser(context.Context, *User) error
		ResetUserPassword(context.Context, *User) error
		SetUserPassword(context.Context, *User, *dto.PasswordInputDTO) error
	}
)

func (u *User) TableName() string {
	return UserTableName
}

func (u *User) ToMap() *map[string]interface{} {
	mapped := &map[string]interface{}{
		"name":       u.Name,
		"mail":       u.Email,
		"status":     u.Status,
		"profile_id": u.ProfileID,
		"new":        u.New,
		"token":      nil,
		"password":   nil,
	}

	if u.Password != nil {
		(*mapped)["token"] = *u.Token
		(*mapped)["password"] = *u.Password
	}

	return mapped
}

func (u *User) Bind(userDTO *dto.UserInputDTO) error {
	if userDTO.Name != nil {
		u.Name = *userDTO.Name
	}
	if userDTO.Email != nil {
		u.Email = *userDTO.Email
	}
	if userDTO.Status != nil {
		u.Status = *userDTO.Status
	}
	if userDTO.ProfileID != nil {
		u.ProfileID = *userDTO.ProfileID
	}

	return validator.StructValidator.Validate(u)
}

func (u *User) ValidatePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(*u.Password), []byte(password)) == nil
}

func (u *User) GenerateToken(expire, originalKey, ip string) (string, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(originalKey)
	if err != nil {
		return "", fmt.Errorf("could not decode key: %v", err.Error())
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedKey)
	if err != nil {
		return "", fmt.Errorf("could not parse key: %v", err.Error())
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"token": u.Token,
		"ip":    ip,
		"iat":   now.Unix(),
	}

	life, err := helpers.DurationFromString(expire, time.Minute)
	if err == nil {
		claims["exp"] = now.Add(life).Unix()
	}
	claims["expire"] = err == nil

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(key)
}
