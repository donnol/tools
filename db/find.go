package db

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
)

type Finder[R any] interface {
	Query() (query string, args []any) // 返回查询语句，参数

	// 新建结果类型对象，不要使用同一个，用来接收结果
	// r需要是指针类型
	// fields需要是字段指针类型，需要与表的列保持一一对应
	NewScanObjAndFields(colTypes []*sql.ColumnType) (r *R, fields []any)
}

type Storer interface {
	*sql.DB | *sql.Tx | *sql.Conn
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func FindList[S Storer, F Finder[R], R any](db S, finder F, res *[]R) (err error) {
	query, args := finder.Query()
	rows, err := db.QueryContext(context.TODO(), query, args...) // sql里select了n列
	if err != nil {
		return
	}
	defer rows.Close()

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return
	}
	for rows.Next() {
		obj, fields := finder.NewScanObjAndFields(colTypes) // fields也必须有n个元素
		if err = rows.Scan(fields...); err != nil {
			return
		}
		// PrintFields(fields)

		*res = append(*res, *obj)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func FindFirst[S Storer, F Finder[R], R any](db S, finder F, res *R) (err error) {
	var r []R
	err = FindList(db, finder, &r)
	if err != nil {
		return
	}
	if len(r) > 0 {
		*res = r[0]
	}
	return
}

// FindAll
// sql里select的字段数量必须与R的ScanFields方法返回的数组元素数量一致
func FindAll[S Storer, F Finder[R], R any](db S, finder F, inital R) (r []R, err error) {
	query, args := finder.Query()
	rows, err := db.QueryContext(context.TODO(), query, args...) // sql里select了n列
	if err != nil {
		return
	}
	defer rows.Close()

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return
	}
	for rows.Next() {
		obj, fields := finder.NewScanObjAndFields(colTypes) // fields也必须有n个元素
		if err = rows.Scan(fields...); err != nil {
			return
		}
		// PrintFields(fields)

		r = append(r, *obj)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func FindOne[S Storer, F Finder[R], R any](db S, finder F, inital R) (r R, err error) {
	res, err := FindAll(db, finder, inital)
	if err != nil {
		return r, err
	}
	if len(res) > 0 {
		r = res[0]
	}
	return
}

func PrintFields(fields []any) {
	fmt.Println("=== begin print fields")
	for i := range fields {
		vall := reflect.ValueOf(fields[i])
		if vall.Kind() == reflect.Pointer {
			vall = vall.Elem()
		}
		fmt.Printf("field: %v\n", vall.Interface())
	}
	fmt.Println("=== end print fields")
}
