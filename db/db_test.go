package db

import (
	"database/sql"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	createSQL = `
	create table if not exists user (
		id integer not null,
		name varchar(255) not null,
		primary key(id)
	);
	INSERT OR IGNORE INTO user (id, name) values (1, 'jd');
	`
)

var (
	tdb = func() *sql.DB {
		gdb, err := gorm.Open(sqlite.Open("./testdata/test.db"))
		if err != nil {
			panic(err)
		}

		sqldb, err := gdb.DB()
		if err != nil {
			panic(err)
		}

		_, err = sqldb.Exec(createSQL)
		if err != nil {
			panic(err)
		}

		return sqldb
	}()
)
