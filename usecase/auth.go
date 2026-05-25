package usecase

import (
	"progas-wms-be/config"
	"progas-wms-be/constant"
	"progas-wms-be/dto"
	"progas-wms-be/global"
	"progas-wms-be/helper"
	"progas-wms-be/repository"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthUseCase interface {
	Login(req *dto.LoginRequest) (*dto.LoginResponse, global.ErrorResponse)
	RefreshToken(req *dto.RefreshTokenRequest) (*dto.LoginResponse, global.ErrorResponse)
	Logout(userId string) global.ErrorResponse
}

type authUseCase struct {
	userRepo     repository.UserRepository
	auditLogRepo repository.AuditLogRepository
}

func NewAuthUseCase(userRepo repository.UserRepository, auditLogRepo repository.AuditLogRepository) AuthUseCase {
	return &authUseCase{
		userRepo:     userRepo,
		auditLogRepo: auditLogRepo,
	}
}

func (u *authUseCase) Login(req *dto.LoginRequest) (*dto.LoginResponse, global.ErrorResponse) {
	// 1. Get user by email
	user, errRes := u.userRepo.FindByEmail(req.Email)
	if errRes != nil {
		return nil, global.BadRequestError("invalid email or password")
	}

	// 2. Check if user is active
	if !user.IsActive {
		return nil, global.BadRequestError("user is not active")
	}

	// 3. Compare password
	if !helper.CheckPasswordHash(req.Password, user.Password) {
		return nil, global.BadRequestError("invalid email or password")
	}

	// 4. Generate token
	accessToken, refreshToken, err := helper.GenerateAuthToken(user.Id, user.RoleId)
	if err != nil {
		return nil, global.InternalServerError(err)
	}

	// 5. Update last logged in
	_ = u.userRepo.UpdateLastLogin(nil, user.Id)

	// 6. Build response
	accessTokenExpiredInMinutes, _ := strconv.Atoi(config.GetEnv(constant.AuthTokenExpiredInMinutes))
	expiredAt := time.Now().Add(time.Duration(accessTokenExpiredInMinutes) * time.Minute)

	roleName := ""
	if user.Role.Id != "" {
		roleName = user.Role.Name
	}

	res := &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiredAt:    expiredAt,
		User: dto.UserResponse{
			Id:       user.Id,
			Name:     user.Name,
			Email:    user.Email,
			Phone:    user.Phone,
			RoleId:   user.RoleId,
			RoleName: roleName,
		},
	}

	_ = u.auditLogRepo.Log(user.Id, constant.AuditUserLogin, constant.AuditObjectUser, user.Id, map[string]string{
		"email": user.Email,
	})

	return res, nil
}

func (u *authUseCase) RefreshToken(req *dto.RefreshTokenRequest) (*dto.LoginResponse, global.ErrorResponse) {
	// 1. Parse token
	token, err := jwt.ParseWithClaims(req.RefreshToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetEnv(constant.RefreshTokenSecretKey)), nil
	})

	if err != nil || !token.Valid {
		return nil, global.BadRequestError("invalid refresh token")
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, global.BadRequestError("invalid refresh token claims")
	}

	userId := claims.Subject

	// 2. Get user by id
	user, errRes := u.userRepo.FindById(userId)
	if errRes != nil {
		return nil, global.BadRequestError("user not found")
	}

	if !user.IsActive {
		return nil, global.BadRequestError("user is not active")
	}

	// 3. Generate new tokens
	newAccessToken, newRefreshToken, err := helper.GenerateAuthToken(user.Id, user.RoleId)
	if err != nil {
		return nil, global.InternalServerError(err)
	}

	// 4. Update last logged in
	_ = u.userRepo.UpdateLastLogin(nil, user.Id)

	// 5. Build response
	accessTokenExpiredInMinutes, _ := strconv.Atoi(config.GetEnv(constant.AuthTokenExpiredInMinutes))
	expiredAt := time.Now().Add(time.Duration(accessTokenExpiredInMinutes) * time.Minute)

	roleName := ""
	if user.Role.Id != "" {
		roleName = user.Role.Name
	}

	res := &dto.LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiredAt:    expiredAt,
		User: dto.UserResponse{
			Id:       user.Id,
			Name:     user.Name,
			Email:    user.Email,
			Phone:    user.Phone,
			RoleId:   user.RoleId,
			RoleName: roleName,
		},
	}

	return res, nil
}

func (u *authUseCase) Logout(userId string) global.ErrorResponse {
	// Stateless logout, so just return nil
	return nil
}
