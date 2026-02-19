package service

import (
	"context"
	"fin-auth/cache"
	"fin-auth/domain"
	"fin-auth/dto"
	"fin-auth/models"
	"fin-auth/utils"
	"strings"
	"time"
)

type Auth struct {
	Repo  domain.AuthRepository
	Cache *cache.RedisCache
}

func NewAuthService(repo domain.AuthRepository, cache *cache.RedisCache) *Auth {
	return &Auth{
		Repo:  repo,
		Cache: cache,
	}
}

func (auth *Auth) RegisterClient(ctx context.Context, req *dto.RegisterClientReq) (*dto.RegisterClientRes, error) {

	clientId := utils.GenerateRandomString(50)
	secret := utils.GenerateRandomString(50)
	secondarySecret := utils.GenerateRandomString(50)

	user := &models.User{
		ClientId:    clientId,
		Name:        req.Name,
		Email:       req.Email,
		IsActive:    req.IsActive,
		Description: req.Description,
	}

	secretModel := &models.Secret{
		ClientId:        clientId,
		Secret:          secret,
		SecondarySecret: secondarySecret,
	}

	res, err := auth.Repo.CreateAuthClient(ctx, user, secretModel)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (auth *Auth) Login(ctx context.Context, req *dto.LoginReq, ipAddress, userAgent string) (*domain.ClientWithSecrets, *dto.TokenResponse, error) {

	clientData, err := auth.Repo.FindClientWithSecrets(ctx, req.ClientId)
	if err != nil {
		return nil, nil, err
	}

	accessToken := utils.GenerateRandomString(50)
	refreshToken := utils.GenerateRandomString(50)

	accessExpiresAt := time.Now().UTC().Add(24 * time.Hour)
	refreshExpiresAt := time.Now().UTC().Add(7 * 24 * time.Hour)

	accessTokenModel := &models.AccessToken{
		ClientId:  req.ClientId,
		Token:     accessToken,
		ExpiredAt: accessExpiresAt,
	}

	err = auth.Repo.CreateAccessToken(ctx, accessTokenModel)
	if err != nil {
		return nil, nil, err
	}

	refreshTokenModel := &models.RefreshToken{
		ClientId:      req.ClientId,
		Token:         refreshToken,
		AccessTokenId: accessTokenModel.ID,
		ExpiredAt:     refreshExpiresAt,
	}

	err = auth.Repo.CreateRefreshToken(ctx, refreshTokenModel)
	if err != nil {
		return nil, nil, err
	}

	if auth.Cache != nil {
		sessionData := &models.SessionData{
			ClientId:     req.ClientId,
			Token:        accessToken,
			LoginTime:    time.Now().UTC(),
			IPAddress:    ipAddress,
			UserAgent:    userAgent,
			DeviceType:   detectDeviceType(userAgent),
			LastActivity: time.Now().UTC(),
		}
		auth.Cache.CreateSession(ctx, req.ClientId, accessToken, sessionData)
	}

	response := &dto.TokenResponse{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresAt:  accessExpiresAt,
		RefreshExpiresAt: refreshExpiresAt,
	}

	return clientData, response, nil
}

func (auth *Auth) Logout(ctx context.Context, token string) error {
	if auth.Cache == nil {
		return nil
	}
	accessToken, err := auth.Repo.FindValidAccessToken(ctx, token)
	if err != nil {
		return nil
	}

	ttl := time.Until(accessToken.ExpiredAt)
	if ttl <= 0 {
		return nil
	}
	if err := auth.Cache.BlacklistToken(ctx, token, ttl); err != nil {
		return err
	}
	if err := auth.Cache.DeleteTokenFromCache(ctx, token); err != nil {
		// Log but don't fail
	}

	if err := auth.Cache.DeleteSession(ctx, accessToken.ClientId, token); err != nil {
		// Log but don't fail
	}

	return nil
}

func (auth *Auth) ListSessions(ctx context.Context, clientId string) ([]dto.SessionResponse, error) {
	if auth.Cache == nil {
		return []dto.SessionResponse{}, nil
	}

	sessions, err := auth.Cache.ListSessions(ctx, clientId)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.SessionResponse, 0, len(sessions))
	for _, session := range sessions {
		responses = append(responses, dto.SessionResponse{
			ClientId:     session.ClientId,
			Token:        session.Token,
			LoginTime:    session.LoginTime,
			IPAddress:    session.IPAddress,
			UserAgent:    session.UserAgent,
			DeviceType:   session.DeviceType,
			LastActivity: session.LastActivity,
		})
	}

	return responses, nil
}

func (auth *Auth) RevokeSession(ctx context.Context, clientId, token string) error {
	return auth.Logout(ctx, token)
}

// Helper function to detect device type from user agent
func detectDeviceType(userAgent string) string {
	ua := strings.ToLower(userAgent)

	if strings.Contains(ua, "mobile") || strings.Contains(ua, "android") || strings.Contains(ua, "iphone") {
		return "mobile"
	}

	if strings.Contains(ua, "tablet") || strings.Contains(ua, "ipad") {
		return "tablet"
	}

	return "desktop"
}

func (auth *Auth) RefreshToken(ctx context.Context, req *dto.RefreshTokenReq) (*dto.RefreshTokenRes, error) {

	refreshToken, err := auth.Repo.FindValidRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	accessToken := utils.GenerateRandomString(50)
	accessExpiresAt := time.Now().UTC().Add(24 * time.Hour)

	accessTokenModel := &models.AccessToken{
		ClientId:  refreshToken.ClientId,
		Token:     accessToken,
		ExpiredAt: accessExpiresAt,
	}

	err = auth.Repo.CreateAccessToken(ctx, accessTokenModel)
	if err != nil {
		return nil, err
	}

	response := &dto.RefreshTokenRes{
		AccessToken:     accessToken,
		AccessExpiresAt: accessExpiresAt,
	}

	return response, nil
}
