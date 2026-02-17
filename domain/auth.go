package domain

import (
	"context"
	"fin-auth/dto"
	"fin-auth/models"
)

type AuthService interface {
	RegisterClient(ctx context.Context, req *dto.RegisterClientReq) (*dto.RegisterClientRes, error)
	Login(ctx context.Context, req *dto.LoginReq, ipAddress, userAgent string) (*ClientWithSecrets, *dto.TokenResponse, error)
	RefreshToken(ctx context.Context, req *dto.RefreshTokenReq) (*dto.RefreshTokenRes, error)
	Logout(ctx context.Context, token string) error
	ListSessions(ctx context.Context, clientId string) ([]dto.SessionResponse, error)
	RevokeSession(ctx context.Context, clientId, token string) error
}

type AuthRepository interface {
	CreateAuthClient(ctx context.Context, user *models.User, secret *models.Secret) (*dto.RegisterClientRes, error)
	FindClientWithSecrets(ctx context.Context, clientId string) (*ClientWithSecrets, error)
	CreateAccessToken(ctx context.Context, token *models.AccessToken) error
	CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error
	FindValidRefreshToken(ctx context.Context, tokenStr string) (*models.RefreshToken, error)
	FindValidAccessToken(ctx context.Context, tokenStr string) (*models.AccessToken, error)
}

type ClientWithSecrets struct {
	User   *models.User
	Secret *models.Secret
}
