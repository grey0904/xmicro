package u_jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func CreateToken(claims Claims, salt []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(salt)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenStr string, claims *Claims, salt string) error {
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(salt), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	// 类型断言以确保得到的是正确的Claims类型
	if _, ok := token.Claims.(*Claims); !ok {
		return fmt.Errorf("invalid claims type")
	}

	return nil
}
