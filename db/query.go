package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/donnol/tools/funcx"
)

func WrapSQLConn(
	ctx context.Context,
	db *sql.DB,
	f func(
		ctx context.Context,
		conn *sql.Conn,
	) error,
) error {

	// 获取连接，并在返回前释放
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("close conn failed: %v", err)
		}
	}()

	if err := f(ctx, conn); err != nil {
		funcName := funcx.FuncName(2)
		return fmt.Errorf("[DB] call in %s failed: %w", funcName, err)
	}

	return nil
}

// WrapSQLQueryRows query by stmt and args, return values with dest
// only support one row
func WrapSQLQueryRows(
	ctx context.Context,
	db *sql.DB,
	stmt string,
	args []interface{},
	dest ...interface{},
) error {

	if err := WrapSQLConn(ctx, db, func(ctx context.Context, conn *sql.Conn) error {

		rows, err := conn.QueryContext(ctx, stmt, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(dest...)
			if err != nil {
				return err
			}
		}
		if err := rows.Err(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// WrapConnFindAll query by stmt and args, return values with dest
// support many rows
func WrapConnFindAll[F Finder[R], R any](
	ctx context.Context,
	db *sql.DB,
	finder F,
	inital R,
) (r []R, err error) {

	if err = WrapSQLConn(ctx, db, func(ctx context.Context, conn *sql.Conn) error {
		r, err = FindAll(conn, finder, inital)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return
	}

	return
}

func WrapTxFindAll[F Finder[R], R any](
	ctx context.Context,
	db *sql.DB,
	finder F,
	inital R,
) (r []R, err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	var success bool
	defer func() {
		if !success {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	r, err = FindAll(tx, finder, inital)
	if err != nil {
		return nil, err
	}
	success = true

	return
}
