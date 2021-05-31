package account

import (
	"context"
	"database/sql"
	"github.com/poonman/entry-task/dora/misc/log"
	"github.com/poonman/entry-task/dora/status"
	"github.com/poonman/entry-task/server/domain/aggr/account"
	"github.com/poonman/entry-task/server/infra/config"
	"time"
)

type repo struct {
	db *sql.DB
}

func (r *repo) Get(ctx context.Context, username string) (a *account.Account, err error) {

	a = &account.Account{}
	row := r.db.QueryRowContext(ctx, "select * from account where username=?", username)
	if err = row.Scan(&a.Id, &a.Username, &a.Password); err != nil {
		err = status.New(status.InternalServerError, "query account error. "+err.Error())
		return
	}

	return
}

func NewRepo(conf *config.Config) account.Repo {
	db, err := sql.Open("mysql", conf.MySQLConfig.SourceName)
	if err != nil {
		log.Fatal("Failed to open mysql database. err:[%v]", err)
	}

	db.SetMaxOpenConns(conf.MySQLConfig.MaxOpenConn)
	db.SetMaxIdleConns(conf.MySQLConfig.MaxIdleConn)
	db.SetConnMaxLifetime(time.Duration(conf.MySQLConfig.ConnMaxLifetime) * time.Second)

	r := &repo{
		db: db,
	}

	return r
}
