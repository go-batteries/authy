package tokens

import (
	"strings"
	"time"

	"github.com/go-batteries/authy/internal/pkg/bob"

	"github.com/google/uuid"
)

type Token struct {
	ClientID         string          `db:"client_id" validate:"required"`
	AccessToken      string          `db:"access_token" validate:"required"`
	RefreshToken     string          `db:"refresh_token" validate:"required"`
	Scopes           bob.StringArray `db:"scopes"`
	Blocked          bool            `db:"blocked"`
	ExpiresAt        time.Time       `db:"expires_at" validate:"required"`
	RefreshExpiresAt time.Time       `db:"refresh_expires_at" validate:"required"`

	CreatedAt time.Time `db:"created_at" validate:"required"`
	UpdatedAt time.Time `db:"updated_at" validate:"required"`
}

func (Token) TableName() string {
	return "tokens"
}

type BuildOpts func(t Token) Token

func WithRefreshToken(refreshToken string) BuildOpts {
	return func(t Token) Token {
		t.RefreshToken = refreshToken
		return t
	}
}

func WithCreatedAt(createdAt time.Time) BuildOpts {
	return func(t Token) Token {
		t.CreatedAt = createdAt
		return t
	}
}

func WithBlocked(blocked bool) BuildOpts {
	return func(t Token) Token {
		t.Blocked = blocked
		return t
	}
}

func AddAccessExpiry(conf ExpiryConfig) BuildOpts {
	return func(t Token) Token {
		now := time.Now().UTC()

		t.ExpiresAt = now.Add(conf.AccessExpiresIn)
		return t
	}
}

func AddRefreshExpiry(conf ExpiryConfig) BuildOpts {
	return func(t Token) Token {
		now := time.Now().UTC()

		t.RefreshExpiresAt = now.Add(conf.RefreshExpiresIn)
		return t
	}
}

func WithRefreshExpiry(refreshExpiresAt time.Time) BuildOpts {
	return func(t Token) Token {
		t.RefreshExpiresAt = refreshExpiresAt
		return t
	}
}

func BuildToken(clientID string, opts ...BuildOpts) Token {
	token := Token{
		ClientID:     clientID,
		AccessToken:  UUID(),
		RefreshToken: UUID(),
		Scopes:       bob.StringArray{"all"},
		Blocked:      false,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	for _, opt := range opts {
		token = opt(token)
	}

	return token
}

func (t Token) HasExpired() bool {
	currentTime := time.Now().UTC().Add(-30 * time.Second) // expire 30 seconds before
	return t.ExpiresAt.UTC().Before(currentTime)
}

func (t Token) HasRefreshExpired() bool {
	currentTime := time.Now().UTC().Add(-30 * time.Second) // expire 30 seconds before
	return t.RefreshExpiresAt.UTC().Before(currentTime)
}

func UUID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}
