package tokens

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/go-batteries/authy/internal/pkg/errorrs"
	"github.com/sirupsen/logrus"
)

type UpdateQueryBuilder interface {
	Build(context.Context, goqu.DialectWrapper, ...interface{}) (string, []interface{}, error)
}

type BlockedUpdater struct {
	Blocked bool `db:"blocked"`
}

func (b BlockedUpdater) Build(ctx context.Context, qb goqu.DialectWrapper, args ...interface{}) (string, []interface{}, error) {
	var (
		sql   string
		dargs []interface{}
		err   error
	)

	if len(args) < 1 {
		return sql, args, errorrs.InternalError(err, errorrs.CodeDatabaseQueryBuilderFailed)
	}

	accessToken := args[0].(string)

	sql, dargs, err = qb.Update("tokens").
		Set(BlockedUpdater{Blocked: true}).
		Where(goqu.Ex{
			"access_token": accessToken,
		}).ToSQL()

	if err != nil {
		logrus.WithContext(ctx).WithError(err).Error("failed to build query for update")
		err = errorrs.InternalError(err, errorrs.CodeDatabaseQueryBuilderFailed)
	}

	return sql, dargs, err
}

type ExpireUpdate struct {
	ExpiresAt time.Time `db:"expires_at"`
}

func (e ExpireUpdate) Build(ctx context.Context, qb goqu.DialectWrapper, args ...interface{}) (sql string, dargs []interface{}, err error) {
	if len(args) < 1 {
		return sql, args, errorrs.InternalError(err, errorrs.CodeDatabaseQueryBuilderFailed)
	}

	accessToken := args[0].(string)

	sql, dargs, err = qb.Update((Token{}).TableName()).
		Set(ExpireUpdate{ExpiresAt: time.Now().UTC()}).
		Where(goqu.Ex{
			"access_token": accessToken,
		}).ToSQL()

	if err != nil {
		logrus.WithContext(ctx).WithError(err).Error("failed to build query for update")
		err = errorrs.InternalError(err, errorrs.CodeDatabaseQueryBuilderFailed)
	}

	return sql, dargs, err
}
