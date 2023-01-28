package authy

import (
	"errors"
	"strings"
	"time"

	"github.com/go-batteries/authy/pkg/config"
	"github.com/go-batteries/authy/src/tokens"
	"github.com/sirupsen/logrus"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func NewAuthorizer(cfg config.Config) (tokens.Service, error) {
	dialect := strings.ToLower(cfg.Dialect)

	logrus.SetFormatter(&logrus.JSONFormatter{})

	var svc tokens.TokenService

	if !strings.Contains("postgres,sqlite3", dialect) {
		return svc, errors.New("unsupported dialect")
	}

	queryBuilder := goqu.Dialect(cfg.Dialect)

	db, err := sqlx.Connect(dialect, cfg.DatabaseConfig.URL)
	if err != nil {
		return svc, err
	}

	db.SetMaxOpenConns(30)
	db.SetMaxIdleConns(10)

	expiryConfig := tokens.ExpiryConfig{
		AccessExpiresIn:  time.Duration(cfg.TokenExpiryInSec),
		RefreshExpiresIn: time.Duration(cfg.RefreshExpiryInSec),
	}

	repo := tokens.NewRepository(db, queryBuilder)

	return tokens.NewTokenService(repo, expiryConfig), nil
}
