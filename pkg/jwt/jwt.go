package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	MySignKey         = []byte("WangZY")
	TokenTimeDuration = time.Hour * 200
)

type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims        // 包含 ExpiresAt, Issuer 等标准字段
}

// 生成token
func GenToken(userID int64, userName string) (aToken string, rToken string) {
	c := MyClaims{
		UserID:   userID,
		Username: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenTimeDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Myself",
		},
	}
	atoken := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenString, err := atoken.SignedString(MySignKey)
	if err != nil {
		return "", ""
	}
	// 生成refreshToken，有效期更长
	refreshClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenTimeDuration * 24 * 7)), // 7天有效期
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "Myself",
	}
	rtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := rtoken.SignedString(MySignKey)
	if err != nil {
		return "", ""
	}
	return tokenString, refreshTokenString
}

// 解析token
func ParseToken(tokenString string) (*MyClaims, error) {
	// 空结构体用于接收
	tempClaims := new(MyClaims)
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, tempClaims, func(t *jwt.Token) (any, error) {
		return MySignKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token invalid")
	}
	return tempClaims, nil
}

// 刷新Token
func RefreshToken(rToken, aToken string) (newRToken, newAToken string, err error) {
	// 解析 refresh token
	refreshClaims := new(jwt.RegisteredClaims)
	rt, err := jwt.ParseWithClaims(rToken, refreshClaims, func(t *jwt.Token) (any, error) {
		return MySignKey, nil
	})
	if err != nil || !rt.Valid {
		return "", "", errors.New("refresh token invalid")
	}

	// 解析 access token（即使过期也要解析出 claims）
	var userID int64
	var username string

	// 尝试解析 access token，获取用户信息
	aClaims := new(MyClaims)
	at, err := jwt.ParseWithClaims(aToken, aClaims, func(t *jwt.Token) (any, error) {
		return MySignKey, nil
	})

	if err != nil {
		// token 过期错误，直接生成新token
		if errors.Is(err, jwt.ErrTokenExpired) {
			userID = aClaims.UserID
			username = aClaims.Username
		} else {
			return "", "", err
		}
	} else if !at.Valid {
		return "", "", errors.New("access token invalid")
	} else {
		userID = aClaims.UserID
		username = aClaims.Username
	}

	// 生成新的 token 对
	newAToken, newRToken = GenToken(userID, username)
	return newRToken, newAToken, nil
}
