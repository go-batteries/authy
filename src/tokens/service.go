package tokens

import (
	"context"
	"errors"
	"time"

	"github.com/go-batteries/authy/internal/pkg/errorrs"

	"github.com/sirupsen/logrus"
)

type Service interface {
	Create(ctx context.Context, clientID string) (Token, error)
	Revoke(ctx context.Context, accessToken string) (bool, error)
	Authenticate(ctx context.Context, accessToken string) (bool, error)
	ReAuthenticate(ctx context.Context, accessToken, refreshToken string) (Token, error)
}

type TokenService struct {
	repo TokensRepository
	ecfg ExpiryConfig
}

type ExpiryConfig struct {
	AccessExpiresIn  time.Duration
	RefreshExpiresIn time.Duration
}

func NewTokenService(repo TokensRepository, econf ExpiryConfig) TokenService {
	return TokenService{repo: repo, ecfg: econf}
}

func (s TokenService) Create(ctx context.Context, clientID string) (Token, error) {
	token := BuildToken(clientID,
		AddAccessExpiry(s.ecfg),
		AddRefreshExpiry(s.ecfg),
	)

	_, err := s.repo.Create(ctx, token)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).Error("failed to create token")
		return token, err
	}

	return token, nil
}

func (s TokenService) Revoke(ctx context.Context, accessToken string) (bool, error) {
	err := s.repo.Block(ctx, accessToken)
	return err == nil, err
}

var (
	ErrAccessBlocked = errors.New("access_blocked")
)

func (s TokenService) Authenticate(ctx context.Context, accessToken string) (bool, error) {
	token, err := s.repo.Find(ctx, FindBy{"access_token": accessToken}, true)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).Error("failed to get token")
		return false, err
	}

	if token.Blocked {
		logrus.WithContext(ctx).Infoln("token is blocked")
		return false, errorrs.UnAuthorized(ErrAccessBlocked, errorrs.CodeAccessBlocked)
	}

	return true, nil
}

var (
	ErrRefreshTokenExpired = errors.New("refresh_expired")
)

func (s TokenService) ReAuthenticate(ctx context.Context, accessToken, refreshToken string) (Token, error) {
	token, err := s.repo.Find(ctx, FindBy{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, true)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).Error("failed to find token refresh token")
		return token, err
	}

	if token.HasRefreshExpired() {
		logrus.WithContext(ctx).Error("refresh token has expired")
		return token, errorrs.UnAuthorized(ErrRefreshTokenExpired, errorrs.CodeRefreshTokenExpired)
	}

	newToken := BuildToken(token.ClientID,
		WithBlocked(token.Blocked),
		WithCreatedAt(token.CreatedAt),
		WithRefreshToken(token.RefreshToken),
		AddAccessExpiry(s.ecfg),
		WithRefreshExpiry(token.RefreshExpiresAt),
	)

	_, err = s.repo.Create(ctx, newToken)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).Error("failed to create new access token while reauthenticating")
		return token, err
	}

	logrus.WithContext(ctx).Infoln("new access token generated ", token.RefreshToken == newToken.RefreshToken)

	return newToken, nil
}
