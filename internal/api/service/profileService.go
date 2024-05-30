package service

import (
	"context"

	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/domain"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/dto"
	"github.com/Duncan-Kiragu/Msaada-Backend/pkg/filter"
)

func NewProfileService(r domain.ProfileRepository) domain.ProfileService {
	return &profileService{
		profileRepository: r,
	}
}

type profileService struct {
	profileRepository domain.ProfileRepository
}

func (s *profileService) generateProfileOutputDTO(profile *domain.Profile) *dto.ProfileOutputDTO {
	return &dto.ProfileOutputDTO{
		Id:   profile.Id,
		Name: profile.Name,
		Permissions: &dto.PermissionsOutputDTO{
			UserModule:    profile.Permissions.UserModule,
			ProfileModule: profile.Permissions.ProfileModule,
			ProductModule: profile.Permissions.ProductModule,
		},
	}
}

// GetProfileByID Implementation of 'GetProfileByID'.
func (s *profileService) GetProfileByID(ctx context.Context, profileID uint) (*dto.ProfileOutputDTO, error) {
	profile, err := s.profileRepository.GetProfileByID(ctx, profileID)
	if err != nil {
		return nil, err
	}

	return s.generateProfileOutputDTO(profile), nil
}

// GetProfiles Implementation of 'GetProfiles'.
func (s *profileService) GetProfiles(ctx context.Context, filter *filter.Filter) (*dto.ListItemsOutputDTO, error) {
	count, err := s.profileRepository.CountProfiles(ctx, filter)
	if err != nil {
		return nil, err
	}

	profiles, err := s.profileRepository.GetProfiles(ctx, filter)
	if err != nil {
		return nil, err
	}

	outputProfiles := &[]dto.ProfileOutputDTO{}
	for _, profile := range *profiles {
		*outputProfiles = append(*outputProfiles, *s.generateProfileOutputDTO(&profile))
	}
	return &dto.ListItemsOutputDTO{
		Count: count,
		Items: outputProfiles,
	}, nil
}

// CreateProfile Implementation of 'CreateProfile'.
func (s *profileService) CreateProfile(ctx context.Context, data *dto.ProfileInputDTO) (*dto.ProfileOutputDTO, error) {
	profile, err := s.profileRepository.CreateProfile(ctx, data)
	if err != nil {
		return nil, err
	}

	return s.generateProfileOutputDTO(profile), nil
}

// UpdateProfile Implementation of 'UpdateProfile'.
func (s *profileService) UpdateProfile(ctx context.Context, profile *domain.Profile, data *dto.ProfileInputDTO) (*dto.ProfileOutputDTO, error) {
	if err := s.profileRepository.UpdateProfile(ctx, profile, data); err != nil {
		return nil, err
	}

	return s.generateProfileOutputDTO(profile), nil
}

// DeleteProfile Implementation of 'DeleteProfile'.
func (s *profileService) DeleteProfile(ctx context.Context, profile *domain.Profile) error {
	return s.profileRepository.DeleteProfile(ctx, profile)
}
