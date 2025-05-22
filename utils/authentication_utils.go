package utils

import (
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
)
	
func GenerateToken(email string) (token string, err error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, 
        jwt.MapClaims{ 
        "email": email, 
        "exp": time.Now().Add(time.Hour * 24).Unix(), 
        })

	key := os.Getenv("JWT_SECRET")
	byteKey := []byte(key)

	token, err = claims.SignedString(byteKey)
	if err != nil{
		return "", err
	}
	
	return token, nil
}