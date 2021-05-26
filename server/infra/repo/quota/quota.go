package quota

import (
	"database/sql"
	"github.com/poonman/entry-task/dora/log"
	"github.com/poonman/entry-task/server/domain/aggr/quota"
	"github.com/poonman/entry-task/server/infra/config"
)

type repo struct {
	db *sql.DB
}

func (r *repo) Get(uid int64) (q *quota.Quota, err error) {
	panic("implement me")
}

func NewRepo(conf *config.Config) quota.Repo {
	db, err := sql.Open("mysql", conf.MySQLConfig.SourceName)
	if err != nil {
		log.Fatal("Failed to open mysql database. err:[%v]", err)
	}

	r := &repo {
		db: db,
	}

	return r
}