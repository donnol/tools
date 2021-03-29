package jwt

import (
	"errors"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	issuer = "jd"
)

// 错误
var (
	ErrBadIssuer = errors.New("Token Bad Issuer")
	ErrExpired   = errors.New("Token Expired")
	ErrNotExist  = errors.New("Token Not Exist")
)

// Token 令牌
type Token struct {
	secret []byte
}

// New 新建
func New(secret []byte) *Token {
	return &Token{
		secret: secret,
	}
}

// CustomClaims 自定义
type CustomClaims struct {
	UserID int `json:"userID"`
	jwt.StandardClaims
}

// Sign 签名
func (t *Token) Sign(userID int) (string, error) {
	// Create the Claims
	claims := &CustomClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(0, 0, 3).Unix(),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(t.secret)
	if err != nil {
		return ss, err
	}

	return ss, nil
}

// Verify 校验
func (t *Token) Verify(tokenString string) (int, error) {
	if strings.TrimSpace(tokenString) == "" {
		return 0, nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return t.secret, nil
	})
	if err != nil {
		return 0, err
	}

	var userID int
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		if claims.Issuer != issuer {
			return 0, ErrBadIssuer
		}
		if claims.ExpiresAt <= time.Now().Unix() {
			return 0, ErrExpired
		}

		userID = claims.UserID
	} else {
		return 0, ErrNotExist
	}

	return userID, nil
}
