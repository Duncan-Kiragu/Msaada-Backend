package domain

import (
	"context"
	"errors"

	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/dto"
)

var ErrInvalidIpAssociation = errors.New("invalid ip source")

type (
	AuthService interface {
		Login(context.Context, *dto.AuthInputDTO, string) (*dto.AuthOutputDTO, error)
		Refresh(*User, string) *dto.AuthOutputDTO
		Me(*User) *dto.UserOutputDTO
	}
)
