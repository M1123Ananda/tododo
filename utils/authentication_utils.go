package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserClaims struct {
	Email string
	Exp   int64
	jwt.RegisteredClaims
}

func GenerateToken(email string) (token string, err error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256,
		UserClaims{
			Email: email,
			Exp:   time.Now().Add(time.Hour * 24).Unix(),
		})

	key := os.Getenv("JWT_SECRET")
	byteKey := []byte(key)

	token, err = claims.SignedString(byteKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func VerifyToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
