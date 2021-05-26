package quota

import (
	"database/sql"
	"github.com/poonman/entry-task/dora/log"
	"github.com/poonman/entry-task/server/domain/aggr/session"
	"github.com/poonman/entry-task/server/infra/config"
)

type repo struct {
	db *sql.DB
}

func NewRepo(conf *config.Config) session.Repo {
	db, err := sql.Open("mysql", conf.MySQLConfig.SourceName)
	if err != nil {
		log.Fatal("Failed to open mysql database. err:[%v]", err)
	}

	r := &repo {
		db: db,
	}

	return r
}