package tests

import (
	"context"
	"testing"

	"github.com/go-batteries/authy"
	"github.com/go-batteries/authy/pkg/config"
	"github.com/go-batteries/authy/src/tokens"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TokenAuthenticationSuite struct {
	suite.Suite
	db  *sqlx.DB
	svc tokens.Service
}

func (s *TokenAuthenticationSuite) SetupSuite() {
	cfg := config.Config{
		DatabaseConfig: &config.DatabaseConfig{
			Dialect: "sqlite3",
			URL:     "../oauth.db",
		},
		AppConfig: &config.AppConfig{
			TokenExpiryInSec:   60,
			RefreshExpiryInSec: 120,
		},
	}

	db, err := sqlx.Connect(cfg.Dialect, cfg.DatabaseConfig.URL)
	require.NoError(s.T(), err, "should not have failed to get database connection")

	svc, err := authy.NewAuthorizer(cfg)
	require.NoError(s.T(), err, "should not have failed to initialize")

	s.db = db
	s.svc = svc
}

func (s *TokenAuthenticationSuite) TearDownTest() {
	_, err := s.db.Exec("DELETE FROM tokens")
	require.NoError(s.T(), err, "should not have failed to delete all tokens")
}

func (s *TokenAuthenticationSuite) TestCreateToken() {
	ctx := context.Background()

	token, err := s.svc.Create(ctx, "client_id")
	require.NoError(s.T(), err, "should not have failed to create token")

	require.NotEmpty(s.T(), token.AccessToken)
	require.NotEmpty(s.T(), token.RefreshToken)

	require.False(s.T(), token.HasExpired())
	require.False(s.T(), token.HasRefreshExpired())
	require.False(s.T(), token.Blocked)
}

func (s *TokenAuthenticationSuite) TestValidateToken() {
	ctx := context.Background()

	token, err := s.svc.Create(ctx, "client_id")
	require.NoError(s.T(), err, "should not have failed to create token")

	valid, err := s.svc.Authenticate(ctx, token.AccessToken)
	require.NoError(s.T(), err, "should not have failed to authenticate")

	require.True(s.T(), valid)

	success, err := s.svc.Revoke(ctx, token.AccessToken)
	require.NoError(s.T(), err, "should not have failed to revoke token")

	require.True(s.T(), success)

	valid, err = s.svc.Authenticate(ctx, token.AccessToken)
	require.Error(s.T(), err, "should have failed to authenticate revoked token")

	require.False(s.T(), valid)
}

func (s *TokenAuthenticationSuite) TestReAuthenticate() {
	ctx := context.Background()

	token, err := s.svc.Create(ctx, "client_id")
	require.NoError(s.T(), err, "should not have failed to create token")

	valid, err := s.svc.Authenticate(ctx, token.AccessToken)
	require.NoError(s.T(), err, "should not have failed to authenticate")
	require.True(s.T(), valid)

	success, err := s.svc.Revoke(ctx, token.AccessToken)
	require.NoError(s.T(), err, "should not have failed to revoke token")
	require.True(s.T(), success)

	newToken, err := s.svc.ReAuthenticate(ctx, token.AccessToken, token.RefreshToken)
	require.NoError(s.T(), err, "should not have failed to reauthenticate")

	require.Equal(s.T(), newToken.RefreshToken, token.RefreshToken)
	require.NotEqual(s.T(), newToken.AccessToken, token.AccessToken)
}

func TestTokenAuthenticationSuite(t *testing.T) {
	suite.Run(t, new(TokenAuthenticationSuite))
}
