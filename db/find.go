package db

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
)

type Finder interface {
	Query() (query string, args []any) // 返回查询语句，参数

	// 新建同类型对象，不要使用同一个，用来接收结果
	// 	这里返回的结果必须是本类型对象，但这种写法其实不能保证(可以返回其它实现了本接口的对象)，使用时注意
	// fields指定对象字段与列的对应关系
	NewScanObjAndFields(colTypes []*sql.ColumnType) (finder Finder, fields []any)
}

type Storer interface {
	*sql.DB | *sql.Tx | *sql.Conn
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// FindAll
// sql里select的字段数量必须与R的ScanFields方法返回的数组元素数量一致
func FindAll[S Storer, R Finder](db S, initial R) (r []R, err error) {
	query, args := initial.Query()
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
		t, fields := initial.NewScanObjAndFields(colTypes) // fields也必须有n个元素
		if err = rows.Scan(fields...); err != nil {
			return
		}
		// PrintFields(fields)

		r = append(r, t.(R))
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func FindOne[S Storer, R Finder](db S, ir R) (r R, err error) {
	res, err := FindAll(db, ir)
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
