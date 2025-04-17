package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/iammrsea/social-app/internal/shared/config"
)

func GetFakeUser() *Authenticateduser {
	return &Authenticateduser{
		Email: "johndoe@example.com",
		Id:    "user-123",
	}
}

func GenerateTestToken(user *Authenticateduser) string {
	claims := AuthClaims{
		UserId: user.Id,
		Email:  user.Email,
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
