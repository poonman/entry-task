package quota

import (
	"context"
	"database/sql"
	"github.com/poonman/entry-task/dora/log"
	"github.com/poonman/entry-task/dora/status"
	"github.com/poonman/entry-task/server/domain/aggr/quota"
	"github.com/poonman/entry-task/server/infra/config"
	"time"
)

type repo struct {
	db *sql.DB
}

func (r *repo) Get(ctx context.Context, username string) (a *quota.Quota, err error) {

	a = &quota.Quota{}
	row := r.db.QueryRowContext(ctx, "select * from quota where username=?", username)
	if err = row.Scan(&a.Id, &a.Username, &a.ReadQuota, &a.WriteQuota); err != nil {
		err = status.New(status.InternalServerError, "query account error")
		return
	}

	return
}

func NewRepo(conf *config.Config) quota.Repo {
	db, err := sql.Open("mysql", conf.MySQLConfig.SourceName)
	if err != nil {
		log.Fatal("Failed to open mysql database. err:[%v]", err)
	}

	db.SetMaxOpenConns(conf.MySQLConfig.MaxOpenConn)
	db.SetMaxIdleConns(conf.MySQLConfig.MaxIdleConn)
	db.SetConnMaxLifetime(time.Duration(conf.MySQLConfig.ConnMaxLifetime)*time.Second)

	r := &repo {
		db: db,
	}

	return r
}