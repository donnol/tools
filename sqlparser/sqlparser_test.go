package sqlparser

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/andreyvit/diff"

	_ "github.com/pingcap/tidb/types/parser_driver"
)

func Test_parse(t *testing.T) {
	type args struct {
		sql string
	}
	tests := []struct {
		name     string
		args     args
		wantCols []string
		wantErr  bool
	}{
		{
			name: "1",
			args: args{
				sql: "SELECT a, b FROM t",
			},
			wantCols: []string{"a", "b"},
		},
		{
			name: "2",
			args: args{
				sql: "SELECT a, b, c FROM t",
			},
			wantCols: []string{"a", "b", "c"},
		},
		{
			name: "3",
			args: args{
				sql: `create table user (
					id integer not null,
					name varchar(255) not null, 
					created_at datetime not null, 
					updated_at timestamp not null
				)`,
			},
			wantCols: []string{"id", "name", "created_at", "updated_at"},
		},
		{
			name: "4",
			args: args{
				sql: `update user set name = 'jd' where id = 1`,
			},
			wantCols: []string{"name", "id"},
		},
		{
			name: "5",
			args: args{
				sql: `insert into user (name) values ('jd')`,
			},
			wantCols: []string{"name"},
		},
		{
			name: "6",
			args: args{
				sql: `delete from user where id = 1`,
			},
			wantCols: []string{"id"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parse(tt.args.sql)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			r := extract(got)
			if !reflect.DeepEqual(r, tt.wantCols) {
				t.Errorf("parse() = %v, want %v", r, tt.wantCols)
			}
		})
	}
}

func TestStruct_Gen(t *testing.T) {
	type fields struct {
		Name    string
		Comment string
		Fields  []Field
	}
	type args struct {
		opt Option
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				Name:    "user_table",
				Comment: "用户表",
				Fields: []Field{
					{
						Name:    "id",
						DBField: "id",
						Type:    "UNSIGNED BIGINT",
						Tag:     "",
						Comment: "主键id",
					},
					{
						Name:    "name",
						DBField: "name",
						Type:    "varchar",
						Tag:     "",
						Comment: "名称",
					},
					{
						Name:    "created_at",
						DBField: "created_at",
						Type:    "datetime",
						Tag:     "",
						Comment: "创建时间",
					},
					{
						Name:    "updated_at",
						DBField: "updated_at",
						Type:    "timestamp",
						Tag:     "",
						Comment: "更新时间",
					},
				},
			},
			args: args{
				opt: Option{},
			},
			wantW: "import \"github.com/donnol/do\"\n\n" +
				"	// UserTable 用户表" + "\n" +
				"	type UserTable struct {" + "\n" +
				"		Id uint64 `json:\"id\" db:\"id\"` // 主键id" + "\n" +
				"		Name string `json:\"name\" db:\"name\"` // 名称" + "\n" +
				"		CreatedAt time.Time `json:\"createdAt\" db:\"created_at\"` // 创建时间" + "\n" +
				"		UpdatedAt time.Time `json:\"updatedAt\" db:\"updated_at\"` // 更新时间" + "\n" +
				"	}" + "\n" +
				`	func (UserTable) TableName() string {
		return "user_table"
	}
	
	func (s UserTable) Columns() []string {
		return s.NameHelper().Columns()
	}
	
	func (s UserTable) Values() []any {
		return []any{
	s.Id,
	s.Name,
	s.CreatedAt,
	s.UpdatedAt,
	
		}
	}
	
	func (s *UserTable) ValuePtrs() []any {
		return []any{
	&s.Id,
	&s.Name,
	&s.CreatedAt,
	&s.UpdatedAt,
	
		}
	}
	type _UserTableNameHelper struct {
		Id string // field: id
		Name string // field: name
		CreatedAt string // field: created_at
		UpdatedAt string // field: updated_at
	}
	// FuzzWrap make v become %v%
	func (_UserTableNameHelper) FuzzWrap(v string) string {
		return "%" + v + "%"
	}
	
	func (_UserTableNameHelper) Columns() []string {
		return []string{
	"id",
	"name",
	"created_at",
	"updated_at",
	
			}
		}
	
	func (UserTable) NameHelper() _UserTableNameHelper {
		return _UserTableNameHelper{
	Id: "id",
	Name: "name",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	
		}
	}
	`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Struct{
				Name:    tt.fields.Name,
				Comment: tt.fields.Comment,
				Fields:  tt.fields.Fields,
			}
			w := &bytes.Buffer{}
			if err := s.Gen(w, tt.args.opt); (err != nil) != tt.wantErr {
				t.Errorf("Struct.Gen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Struct.Gen() = %v, want %v, diff: %s", gotW, tt.wantW, diff.LineDiff(gotW, tt.wantW))
			}
		})
	}
}

func TestParseCreateSQL(t *testing.T) {
	type args struct {
		sql string
		opt Option
	}
	tests := []struct {
		name  string
		args  args
		want  *Struct
		wantW string
	}{
		{
			name: "1",
			args: args{
				sql: `create table user (
id integer unsigned not null comment 'id',
name varchar(255) not null comment '名称', 
created_at datetime not null comment '创建时间', 
updated_at timestamp not null comment '更新时间'
) comment '用户表'`,
			},
			want: &Struct{
				Name:    "user",
				Comment: "用户表",
				Fields: []Field{
					{Name: "id", DBField: "id", Type: "int unsigned", Tag: "", Comment: "id"},
					{Name: "name", DBField: "name", Type: "varchar", Tag: "", Comment: "名称"},
					{Name: "created_at", DBField: "created_at", Type: "datetime", Tag: "", Comment: "创建时间"},
					{Name: "updated_at", DBField: "updated_at", Type: "timestamp", Tag: "", Comment: "更新时间"},
				},
			},
			wantW: "import \"github.com/donnol/do\"\n\n" +
				"	// User 用户表" + "\n" +
				"	type User struct {" + "\n" +
				"		Id uint `json:\"id\" db:\"id\"` // id" + "\n" +
				"		Name string `json:\"name\" db:\"name\"` // 名称" + "\n" +
				"		CreatedAt time.Time `json:\"createdAt\" db:\"created_at\"` // 创建时间" + "\n" +
				"		UpdatedAt time.Time `json:\"updatedAt\" db:\"updated_at\"` // 更新时间" + "\n" +
				"	}" + "\n" +
				`	func (User) TableName() string {
		return "user"
	}
	
	func (s User) Columns() []string {
		return s.NameHelper().Columns()
	}
	
	func (s User) Values() []any {
		return []any{
	s.Id,
	s.Name,
	s.CreatedAt,
	s.UpdatedAt,
	
		}
	}
	
	func (s *User) ValuePtrs() []any {
		return []any{
	&s.Id,
	&s.Name,
	&s.CreatedAt,
	&s.UpdatedAt,
	
		}
	}
	type _UserNameHelper struct {
		Id string // field: id
		Name string // field: name
		CreatedAt string // field: created_at
		UpdatedAt string // field: updated_at
	}
	// FuzzWrap make v become %v%
	func (_UserNameHelper) FuzzWrap(v string) string {
		return "%" + v + "%"
	}
	
	func (_UserNameHelper) Columns() []string {
		return []string{
	"id",
	"name",
	"created_at",
	"updated_at",
	
			}
		}
	
	func (User) NameHelper() _UserNameHelper {
		return _UserNameHelper{
	Id: "id",
	Name: "name",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	
		}
	}
	`,
		},
		{
			name: "ignoreField",
			args: args{
				sql: `create table user (
id integer unsigned not null comment 'id',
name varchar(255) not null comment '名称', 
created_at datetime not null comment '创建时间', 
updated_at timestamp not null comment '更新时间'
) comment '用户表'`,
				opt: Option{
					IgnoreField: []string{"updated_at"},
				},
			},
			want: &Struct{
				Name:    "user",
				Comment: "用户表",
				Fields: []Field{
					{Name: "id", DBField: "id", Type: "int unsigned", Tag: "", Comment: "id"},
					{Name: "name", DBField: "name", Type: "varchar", Tag: "", Comment: "名称"},
					{Name: "created_at", DBField: "created_at", Type: "datetime", Tag: "", Comment: "创建时间"},
					{Name: "updated_at", DBField: "updated_at", Type: "timestamp", Tag: "", Comment: "更新时间"},
				},
			},
			wantW: "import \"github.com/donnol/do\"\n\n" +
				"	// User 用户表" + "\n" +
				"	type User struct {" + "\n" +
				"		Id uint `json:\"id\" db:\"id\"` // id" + "\n" +
				"		Name string `json:\"name\" db:\"name\"` // 名称" + "\n" +
				"		CreatedAt time.Time `json:\"createdAt\" db:\"created_at\"` // 创建时间" + "\n" +
				"	}" + "\n" +
				`	func (User) TableName() string {
		return "user"
	}
	
	func (s User) Columns() []string {
		return s.NameHelper().Columns()
	}
	
	func (s User) Values() []any {
		return []any{
	s.Id,
	s.Name,
	s.CreatedAt,
	
		}
	}
	
	func (s *User) ValuePtrs() []any {
		return []any{
	&s.Id,
	&s.Name,
	&s.CreatedAt,
	
		}
	}
	type _UserNameHelper struct {
		Id string // field: id
		Name string // field: name
		CreatedAt string // field: created_at
	}
	// FuzzWrap make v become %v%
	func (_UserNameHelper) FuzzWrap(v string) string {
		return "%" + v + "%"
	}
	
	func (_UserNameHelper) Columns() []string {
		return []string{
	"id",
	"name",
	"created_at",
	
			}
		}
	
	func (User) NameHelper() _UserNameHelper {
		return _UserNameHelper{
	Id: "id",
	Name: "name",
	CreatedAt: "created_at",
	
		}
	}
	`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseCreateSQL(tt.args.sql); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCreateSQL() = %+v, want %+v", got, tt.want)
			} else {
				buf := new(bytes.Buffer)
				if err := got.Gen(buf, tt.args.opt); err != nil {
					t.Error(err)
				}
				if buf.String() != tt.wantW {
					t.Errorf("Struct.Gen() = %v, want %v, diff: %s", buf.String(), tt.wantW, diff.LineDiff(buf.String(), tt.wantW))
				}
			}
		})
	}
}

func TestStruct_GenData(t *testing.T) {
	type args struct {
		sql string
		n   int64
		opt Option
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name: "full",
			args: args{
				sql: `create table user (
					id integer unsigned not null comment 'id',
					name varchar(255) not null comment '名称', 
					created_at datetime not null comment '创建时间', 
					updated_at timestamp not null comment '更新时间'
					) comment '用户表'`,
				n:   0,
				opt: doption,
			},
			wantW: "INSERT IGNORE INTO `user` (" + "\n" +
				"`id`," + "\n" +
				"`name`," + "\n" +
				"`created_at`," + "\n" +
				"`updated_at`" + "\n" +
				") VALUES (",
			wantErr: false,
		},
		{
			name: "ignore",
			args: args{
				sql: `create table user (
					id integer unsigned not null comment 'id',
					name varchar(255) not null comment '名称', 
					created_at datetime not null comment '创建时间', 
					updated_at timestamp not null comment '更新时间'
					) comment '用户表'`,
				n: 0,
				opt: Option{
					IgnoreField: []string{"updated_at"},
				},
			},
			wantW: "INSERT IGNORE INTO `user` (" + "\n" +
				"`id`," + "\n" +
				"`name`," + "\n" +
				"`created_at`" + "\n" +
				") VALUES (",
			wantErr: false,
		},
		{
			name: "ignore",
			args: args{
				sql: `create table user (
					id integer unsigned not null comment 'id',
					name varchar(255) not null comment '名称', 
					created_at datetime not null comment '创建时间', 
					updated_at timestamp not null comment '更新时间'
					) comment '用户表'`,
				n: 2,
				opt: Option{
					IgnoreField: []string{"updated_at"},
				},
			},
			wantW: "INSERT IGNORE INTO `user` (" + "\n" +
				"`id`," + "\n" +
				"`name`," + "\n" +
				"`created_at`" + "\n" +
				") VALUES (",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ParseCreateSQL(tt.args.sql)
			w := &bytes.Buffer{}
			if err := s.GenData(w, tt.args.n, tt.args.opt); (err != nil) != tt.wantErr {
				t.Errorf("Struct.GenData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// 因为值是随机生成的，所以只比较前面部分
			if gotW := w.String(); len(gotW) <= len(tt.wantW) {
				t.Errorf("Struct.GenData() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestParseCreateSQLBatch(t *testing.T) {
	type args struct {
		sql string
	}
	tests := []struct {
		name string
		args args
		want []*Struct
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				sql: `create table user (
					id integer not null,
					name varchar(255) not null, 
					created_at datetime not null, 
					updated_at timestamp not null
				);
				create table role (
					id integer not null,
					name varchar(255) not null, 
					created_at datetime not null, 
					updated_at timestamp not null
				);
				`,
			},
			want: []*Struct{
				{Name: "user", Fields: []Field{
					{
						Name:    "id",
						Type:    "int",
						DBField: "id",
					},
					{
						Name:    "name",
						Type:    "varchar",
						DBField: "name",
					},
					{
						Name:    "created_at",
						Type:    "datetime",
						DBField: "created_at",
					},
					{
						Name:    "updated_at",
						Type:    "timestamp",
						DBField: "updated_at",
					},
				}},
				{Name: "role", Fields: []Field{
					{
						Name:    "id",
						Type:    "int",
						DBField: "id",
					},
					{
						Name:    "name",
						Type:    "varchar",
						DBField: "name",
					},
					{
						Name:    "created_at",
						Type:    "datetime",
						DBField: "created_at",
					},
					{
						Name:    "updated_at",
						Type:    "timestamp",
						DBField: "updated_at",
					},
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseCreateSQLBatch(tt.args.sql)
			for i, one := range got {
				if !reflect.DeepEqual(*one, *tt.want[i]) {
					t.Errorf("ParseCreateSQLBatch() = %v, want %v", *one, *tt.want[i])
				}
			}
		})
	}
}

func Test_processFieldType(t *testing.T) {
	type args struct {
		fieldType string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				fieldType: "integer",
			},
			want: "integer",
		},
		{
			name: "2",
			args: args{
				fieldType: "varchar(255)",
			},
			want: "varchar",
		},
		{
			name: "3",
			args: args{
				fieldType: "double(10,2)",
			},
			want: "double",
		},
		{
			name: "4",
			args: args{
				fieldType: "BIGINT UNSIGNED",
			},
			want: "BIGINT UNSIGNED",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processFieldType(tt.args.fieldType); got != tt.want {
				t.Errorf("processFieldType() = %v, want %v", got, tt.want)
			}
		})
	}
}
