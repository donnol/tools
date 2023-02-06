package db

import (
	"database/sql"
	"testing"
)

// === 对应表结构 ===

type User struct {
	Id   uint64
	Name string
}

// === 查询 ===

type finderOfUser struct {
	id uint64
}

func (f *finderOfUser) Query() (query string, args []any) {
	query = `select * from user where id = ?`
	args = append(args, f.id)
	return
}

func (f *finderOfUser) NewScanObjAndFields(colTypes []*sql.ColumnType) (r *User, fields []any) {
	r = &User{}
	fields = append(fields,
		&r.Id,
		&r.Name,
	)
	return
}

func TestFindAll(t *testing.T) {
	finder := &finderOfUser{
		id: 1,
	}
	r, err := FindAll(tdb, finder, (*User)(nil))
	if err != nil {
		t.Error(err)
	}
	if len(r) != 1 {
		t.Errorf("bad result len, len(r) %v != %v", len(r), 1)
	}
	for _, r := range r {
		if r.Id != 1 {
			t.Errorf("bad case of id, %v != %v", r.Id, 1)
		}
		if r.Name != "jd" {
			t.Errorf("bad case of name, %v != %v", r.Name, "jd")
		}
	}
}

type finderOfUserOnlyId struct {
	id uint64
}

func (f *finderOfUserOnlyId) Query() (query string, args []any) {
	query = `select id from user where id = ?`
	args = append(args, f.id)
	return
}

func (f *finderOfUserOnlyId) NewScanObjAndFields(colTypes []*sql.ColumnType) (r *User, fields []any) {
	r = &User{}
	fields = append(fields,
		&r.Id,
	)
	return
}

func TestFindAll_Id(t *testing.T) {
	finder := &finderOfUserOnlyId{
		id: 1,
	}
	r, err := FindAll(tdb, finder, (*User)(nil))
	if err != nil {
		t.Error(err)
	}
	if len(r) != 1 {
		t.Errorf("bad result len, len(r) %v != %v", len(r), 1)
	}
	for _, r := range r {
		if r.Id != 1 {
			t.Errorf("bad case of id, %v != %v", r.Id, 1)
		}
		if r.Name != "" {
			t.Errorf("bad case of name, %v != %v", r.Name, "")
		}
	}
}

type finderOfUserOnlyTotal struct {
	id uint64
}

func (f *finderOfUserOnlyTotal) Query() (query string, args []any) {
	query = `select count(id) as count from user where id = ?`
	args = append(args, f.id)
	return
}

func (f *finderOfUserOnlyTotal) NewScanObjAndFields(colTypes []*sql.ColumnType) (r *int64, fields []any) {
	var total int64
	r = &total
	fields = append(fields,
		r,
	)
	return
}

func TestFindAll_Total(t *testing.T) {
	finder := &finderOfUserOnlyTotal{
		id: 1,
	}
	r, err := FindAll(tdb, finder, (*int64)(nil))
	if err != nil {
		t.Error(err)
	}
	if len(r) != 1 {
		t.Errorf("bad result len1, len(r) %v != %v", len(r), 1)
	}
	for _, r := range r {
		if *r != 1 {
			t.Errorf("bad case of id, %v != %v", r, 1)
		}
	}
}
