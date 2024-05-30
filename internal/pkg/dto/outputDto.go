package dto

type (
	ListItemsOutputDTO struct {
		Items interface{} `json:"items"`
		Count int64       `json:"count"`
	}

	PermissionsOutputDTO struct {
		UserModule    bool `json:"user_module" example:"true"`
		ProfileModule bool `json:"profile_module" example:"true"`
		ProductModule bool `json:"product_module" example:"true"`
	}

	ProfileOutputDTO struct {
		Id          uint                  `json:"id" example:"1"`
		Name        string                `json:"name" example:"ADMIN"`
		Permissions *PermissionsOutputDTO `json:"permissions,omitempty"`
	}

	ProductOutputDTO struct {
		Id   uint   `json:"id" example:"1"`
		Name string `json:"name" example:"Product 01"`
	}

	UserOutputDTO struct {
		Id      uint             `json:"id" example:"1"`
		Name    string           `json:"name" example:"John Cena"`
		Email   string           `json:"email" example:"john.cena@email.com"`
		Status  bool             `json:"status" example:"true"`
		Profile ProfileOutputDTO `json:"profile"`
	}

	AuthOutputDTO struct {
		User         *UserOutputDTO `json:"user,omitempty"`
		AccessToken  string         `json:"accesstoken"`
		RefreshToken string         `json:"refreshtoken"`
	}
)
