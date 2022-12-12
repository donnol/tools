package sqlparser

import (
	"bytes"
	"reflect"
	"testing"

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
						Type:    "UNSIGNED BIGINT",
						Tag:     "",
						Comment: "主键id",
					},
					{
						Name:    "name",
						Type:    "varchar",
						Tag:     "",
						Comment: "名称",
					},
					{
						Name:    "created_at",
						Type:    "datetime",
						Tag:     "",
						Comment: "创建时间",
					},
					{
						Name:    "updated_at",
						Type:    "timestamp",
						Tag:     "",
						Comment: "更新时间",
					},
				},
			},
			args: args{
				opt: Option{},
			},
			wantW: "\n" +
				"	// UserTable 用户表" + "\n" +
				"	type UserTable struct {" + "\n" +
				"		Id uint64 `json:\"id\" db:\"id\"` // 主键id" + "\n" +
				"		Name string `json:\"name\" db:\"name\"` // 名称" + "\n" +
				"		CreatedAt time.Time `json:\"createdAt\" db:\"created_at\"` // 创建时间" + "\n" +
				"		UpdatedAt time.Time `json:\"updatedAt\" db:\"updated_at\"` // 更新时间" + "\n" +
				"	}" + "\n" +
				`	func (UserTable) TableName() string {
		return "user_table"
	}`,
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
				t.Errorf("Struct.Gen() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestParseCreateSQL(t *testing.T) {
	type args struct {
		sql string
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
					{Name: "id", Type: "int unsigned", Tag: "", Comment: "id"},
					{Name: "name", Type: "varchar", Tag: "", Comment: "名称"},
					{Name: "created_at", Type: "datetime", Tag: "", Comment: "创建时间"},
					{Name: "updated_at", Type: "timestamp", Tag: "", Comment: "更新时间"},
				},
			},
			wantW: "\n" +
				"	// User 用户表" + "\n" +
				"	type User struct {" + "\n" +
				"		Id uint `json:\"id\" db:\"id\"` // id" + "\n" +
				"		Name string `json:\"name\" db:\"name\"` // 名称" + "\n" +
				"		CreatedAt time.Time `json:\"createdAt\" db:\"created_at\"` // 创建时间" + "\n" +
				"		UpdatedAt time.Time `json:\"updatedAt\" db:\"updated_at\"` // 更新时间" + "\n" +
				"	}" + "\n" +
				`	func (User) TableName() string {
		return "user"
	}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseCreateSQL(tt.args.sql); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCreateSQL() = %+v, want %+v", got, tt.want)
			} else {
				buf := new(bytes.Buffer)
				if err := got.Gen(buf, Option{}); err != nil {
					t.Error(err)
				}
				if buf.String() != tt.wantW {
					t.Errorf("Struct.Gen() = %v, want %v", buf.String(), tt.wantW)
				}
			}
		})
	}
}
