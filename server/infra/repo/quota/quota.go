package quota

import (
	"context"
	"database/sql"
	"github.com/poonman/entry-task/dora/log"
	"github.com/poonman/entry-task/dora/status"
	"github.com/poonman/entry-task/server/domain/aggr/quota"
	"github.com/poonman/entry-task/server/infra/config"
	"strconv"
	"time"
)

type repo struct {
	conf *config.Config
	db   *sql.DB
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

var (
	MaxQuota = 100000
)

func (r *repo) Get(ctx context.Context, username string) (a *quota.Quota, err error) {

	defer func() {
		log.Debugf("quota:[%+v]", a)
	}()

	if r.conf.QuotaRepoConfig.FixedQuota > 0 {
		return &quota.Quota{
			Id:         0,
			Username:   username,
			ReadQuota:  r.conf.QuotaRepoConfig.FixedQuota,
			WriteQuota: r.conf.QuotaRepoConfig.FixedQuota,
		}, nil
	}
	if !r.conf.QuotaRepoConfig.UseMySQL {
		index, err := strconv.ParseInt(username, 10, 64)
		if err != nil {
			err = status.New(status.Internal, "username is invalid")
			return nil, err
		}

		uid := int(index)

		a = &quota.Quota{
			Id:         0,
			Username:   username,
			ReadQuota:  min(MaxQuota, uid*3+3),
			WriteQuota: min(MaxQuota, uid+3),
		}

	} else {

		a = &quota.Quota{}
		row := r.db.QueryRowContext(ctx, "select * from quota where username=?", username)
		if err = row.Scan(&a.Id, &a.Username, &a.ReadQuota, &a.WriteQuota); err != nil {
			err = status.New(status.InternalServerError, "query account error")
			return
		}

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
	db.SetConnMaxLifetime(time.Duration(conf.MySQLConfig.ConnMaxLifetime) * time.Second)

	r := &repo{
		db:   db,
		conf: conf,
	}

	return r
}
