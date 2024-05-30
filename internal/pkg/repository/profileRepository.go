package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/domain"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/dto"
	"github.com/Duncan-Kiragu/Msaada-Backend/pkg/filter"
)

func NewProfileRepository(db *gorm.DB) domain.ProfileRepository {
	return &profileRepository{
		db: db,
	}
}

type profileRepository struct {
	db *gorm.DB
}

func (s *profileRepository) applyFilter(ctx context.Context, filter *filter.Filter) *gorm.DB {
	db := s.db.WithContext(ctx)
	db = filter.ApplySearchLike(db, "name")

	return filter.ApplyOrder(db)
}

func (s *profileRepository) CountProfiles(ctx context.Context, filter *filter.Filter) (int64, error) {
	var count int64
	db := s.applyFilter(ctx, filter)
	return count, db.Model(&domain.Profile{}).Count(&count).Error
}

func (s *profileRepository) GetProfiles(ctx context.Context, filter *filter.Filter) (*[]domain.Profile, error) {
	db := filter.ApplyPagination(s.applyFilter(ctx, filter))

	profiles := &[]domain.Profile{}
	return profiles, db.Preload(clause.Associations).Find(profiles).Error
}

func (s *profileRepository) GetProfileByID(ctx context.Context, profileID uint) (*domain.Profile, error) {
	profile := &domain.Profile{}
	return profile, s.db.WithContext(ctx).Preload(clause.Associations).First(profile, profileID).Error
}

func (s *profileRepository) CreateProfile(ctx context.Context, data *dto.ProfileInputDTO) (*domain.Profile, error) {
	profile := &domain.Profile{}
	if err := profile.Bind(data); err != nil {
		return nil, err
	}

	return profile, s.db.WithContext(ctx).Create(profile).Error
}

func (s *profileRepository) UpdateProfile(ctx context.Context, profile *domain.Profile, data *dto.ProfileInputDTO) error {
	if err := profile.Bind(data); err != nil {
		return err
	}

	if err := s.db.WithContext(ctx).Model(&profile.Permissions).Updates(profile.Permissions.ToMap()).Error; err != nil {
		return err
	}

	return s.db.WithContext(ctx).Updates(profile).Error
}

func (s *profileRepository) DeleteProfile(ctx context.Context, profile *domain.Profile) error {
	return s.db.WithContext(ctx).Select(clause.Associations).Delete(profile).Error
}
