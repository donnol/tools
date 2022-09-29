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
