package tokens

import (
	"context"
	"fmt"
	"time"

	"github.com/go-batteries/authy/internal/pkg/errorrs"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Repository interface{}

type TokensRepository struct {
	db *sqlx.DB
	qb goqu.DialectWrapper
}

func NewRepository(db *sqlx.DB, qb goqu.DialectWrapper) TokensRepository {
	return TokensRepository{db: db, qb: qb}
}

func (t TokensRepository) Create(ctx context.Context, token Token) (Token, error) {
	var err error

	token.CreatedAt = time.Now().UTC()
	token.UpdatedAt = time.Now().UTC()

	sql, args, err := t.qb.Insert(token.TableName()).Rows(token).ToSQL()
	if err != nil {
		fmt.Println(err)
		logrus.WithContext(ctx).WithError(err).Error("failed to build insert query")
		return token, errorrs.InternalError(err, errorrs.CodeDatabaseQueryBuilderFailed)
	}

	logrus.WithContext(ctx).Infoln(sql)

	_, err = t.db.ExecContext(ctx, sql, args...)
	if err != nil {
		fmt.Println(err)
		return token, errorrs.InternalError(err, errorrs.CodeDatabaseInsertFailed)
	}

	logrus.WithContext(ctx).Infoln("token created")

	return token, nil
}

type FindBy map[string]interface{}

func (t TokensRepository) Find(ctx context.Context, clauses FindBy, checkExpiry bool) (Token, error) {
	token := Token{}

	var ex = goqu.Ex{}
	for key, value := range clauses {
		ex[key] = value
	}

	sql, args, err := t.qb.From(token.TableName()).Where(ex).Select(&token).ToSQL()
	if err != nil {
		logrus.WithError(err).Error("failed to build query")
		return token, err
	}

	logrus.WithContext(ctx).Infoln(sql)

	var newToken Token
	err = t.db.GetContext(ctx, &newToken, sql, args...)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).Error("failed to get token from db")
		return token, err
	}

	if checkExpiry && newToken.HasExpired() {
		logrus.WithContext(ctx).Error("access token expired")
		return token, errorrs.UnAuthorized(err, errorrs.CodeTokenExpired)
	}

	logrus.WithContext(ctx).Infoln("token found")
	return newToken, nil
}

func (t TokensRepository) Update(ctx context.Context, builder UpdateQueryBuilder, bargs ...interface{}) error {
	sql, args, err := builder.Build(ctx, t.qb, bargs...)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).Error("failed to build query for update")
		return err
	}

	_, err = t.db.ExecContext(ctx, sql, args...)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).Error("failed to block token in db")
		return errorrs.InternalError(err, errorrs.CodeUpdateQueryFailed)
	}

	return nil
}
