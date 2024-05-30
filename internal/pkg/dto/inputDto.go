package dto

type (
	ProductInputDTO struct {
		Name *string `json:"name" example:"Product 01"`
	}

	PermissionsInputDTO struct {
		UserModule    *bool `json:"user_module" example:"true"`
		ProfileModule *bool `json:"profile_module" example:"true"`
		ProductModule *bool `json:"product_module" example:"true"`
	}

	ProfileInputDTO struct {
		Name        *string             `json:"name" example:"ADMIN"`
		Permissions PermissionsInputDTO `json:"permissions"`
	}

	UserInputDTO struct {
		Name      *string `json:"name" example:"John Cena"`
		Email     *string `json:"email" example:"john.cena@email.com"`
		Status    *bool   `json:"status" example:"true"`
		ProfileID *uint   `json:"profile_id" example:"1"`
	}

	PasswordInputDTO struct {
		Password        *string `json:"password" example:"secret"`
		PasswordConfirm *string `json:"password_confirm" example:"secret"`
	}

	AuthInputDTO struct {
		Login    string `json:"login" example:"admin@admin.com"`
		Password string `json:"password" example:"12345678"`
		Expire   bool   `json:"expire" example:"false"`
	}
)

func (p PasswordInputDTO) IsValid() bool {
	if p.Password == nil || p.PasswordConfirm == nil {
		return false
	}

	if len(*p.Password) < 5 || len(*p.PasswordConfirm) < 5 {
		return false
	}

	return *p.Password == *p.PasswordConfirm
}
