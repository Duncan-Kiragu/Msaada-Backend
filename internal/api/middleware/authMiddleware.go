package middleware

import (
	"encoding/base64"
	"errors"

	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/domain"
	"github.com/Duncan-Kiragu/Msaada-Backend/pkg/helpers"
	"github.com/golang-jwt/jwt/v5"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"

	httphelper "github.com/Duncan-Kiragu/Msaada-Backend/pkg/http-helper"
)

var (
	MidAccess  fiber.Handler
	MidRefresh fiber.Handler
)

func Auth(base64key string, repo domain.UserRepository) fiber.Handler {
	decodedKey, err := base64.StdEncoding.DecodeString(base64key)
	helpers.PanicIfErr(err)

	parsedKey, err := jwt.ParseRSAPublicKeyFromPEM(decodedKey)
	helpers.PanicIfErr(err)

	return keyauth.New(keyauth.Config{
		KeyLookup:  "header:" + fiber.HeaderAuthorization,
		AuthScheme: "Bearer",
		ContextKey: "token",
		Next: func(c *fiber.Ctx) bool {
			// Filter request to skip middleware
			// true to skip, false to not skip
			return false
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return httphelper.NewHTTPResponse(c, fiber.StatusUnauthorized, err)
		},
		Validator: func(c *fiber.Ctx, key string) (bool, error) {
			parsedToken, err := jwt.Parse(key, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, err
				}

				return parsedKey, nil
			})
			if err != nil {
				return false, err
			}

			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			if !ok || !parsedToken.Valid {
				return false, errors.New("invalid jwt token")
			}

			user, err := repo.GetUserByToken(c.Context(), claims["token"].(string))
			if err != nil {
				return false, err
			}
			if val, ok := claims["ip"]; !ok || val.(string) != c.IP() {
				return false, domain.ErrInvalidIpAssociation
			}
			if val, ok := claims["expire"]; ok {
				user.Expire = val.(bool)
			}

			if !user.Status {
				return false, errors.New("invalid user")
			}

			c.Locals(httphelper.LocalUser, user)
			return true, nil
		},
	})
}
