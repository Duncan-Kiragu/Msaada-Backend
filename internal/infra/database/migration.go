package database

import (
	"context"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/domain"
	"github.com/Duncan-Kiragu/Msaada-Backend/pkg/helpers"
)

func autoMigrate(db *gorm.DB) {
	helpers.PanicIfErr(db.AutoMigrate(&domain.Permissions{}))
	helpers.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
	helpers.PanicIfErr(db.AutoMigrate(&domain.User{}))
	helpers.PanicIfErr(db.AutoMigrate(&domain.Product{}))
}

func createDefaults(db *gorm.DB) {
	profile := &domain.Profile{
		Name: "ROOT",
		Permissions: domain.Permissions{
			UserModule:    true,
			ProfileModule: true,
			ProductModule: true,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	helpers.PanicIfErr(db.WithContext(ctx).FirstOrCreate(profile, "name = ?", profile.Name).Error)

	user := &domain.User{
		Name:      os.Getenv("ADM_NAME"),
		Email:     os.Getenv("ADM_MAIL"),
		Status:    true,
		ProfileID: profile.Id,
		New:       false,
		Token:     new(string),
		Password:  new(string),
	}

	token := uuid.New().String()
	*user.Token = token

	hash, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADM_PASS")), bcrypt.DefaultCost)
	helpers.PanicIfErr(err)
	user.Password = new(string)
	*user.Password = string(hash)

	helpers.PanicIfErr(db.WithContext(ctx).FirstOrCreate(user, "mail = ?", user.Email).Error)
}
