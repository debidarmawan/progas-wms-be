package helper

import (
	"progas-wms-be/config"
	"progas-wms-be/constant"
	"progas-wms-be/dto"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAuthToken(userId, roleId string) (string, string, error) {
	accessTokenExpiredInMinutes, _ := strconv.Atoi(config.GetEnv(constant.AuthTokenExpiredInMinutes))
	accessTokenClaims := &dto.JWTClaims{
		UserId: userId,
		RoleId: roleId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(accessTokenExpiredInMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	tokenString, err := accessToken.SignedString([]byte(config.GetEnv(constant.AuthTokenSecretKey)))
	if err != nil {
		return "", "", err
	}

	refreshTokenExpiredInDays, _ := strconv.Atoi(config.GetEnv(constant.RefreshTokenExpiredInDays))
	refreshTokenClaims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(refreshTokenExpiredInDays) * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   string(userId),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(config.GetEnv(constant.RefreshTokenSecretKey)))
	if err != nil {
		return "", "", err
	}
	return tokenString, refreshTokenString, nil
}
