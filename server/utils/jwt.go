package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID      uint   `json:"user_id"`
	Username    string `json:"username"`
	AuthorityId uint   `json:"authority_id"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte("go-gin-element-admin-secret-key")

// GenerateToken 生成JWT token
func GenerateToken(userID uint, username string, authorityId uint) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(7 * 24 * time.Hour) // 7天过期

	claims := Claims{
		UserID:      userID,
		Username:    username,
		AuthorityId: authorityId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    "go-gin-element-admin",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken 解析JWT token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

// ValidateToken 验证token是否有效
func ValidateToken(token string) error {
	claims, err := ParseToken(token)
	if err != nil {
		return err
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		return errors.New("token已过期")
	}

	return nil
}
