package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/iammrsea/social-app/internal/shared/config"
	"github.com/iammrsea/social-app/internal/user/domain"
)

func GetFakeUser() *AuthenticatedUser {
	return &AuthenticatedUser{
		Email: "johndoe@example.com",
		Id:    "user-123",
		Role:  domain.Moderator,
	}
}

func GenerateTestToken(user *AuthenticatedUser) string {
	claims := AuthClaims{
		UserId: user.Id,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(config.Env().AuthSecret())
	signedToken, err := token.SignedString(secret)

	if err != nil {
		panic(err)
	}
	return signedToken
}
