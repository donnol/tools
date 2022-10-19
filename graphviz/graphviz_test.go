package graphviz

import (
	"os"
	"testing"
	"time"
)

type Article struct {
	Id       uint
	Author   uint
	Title    string
	Content  string
	Created  time.Time
	Modified time.Time
}

type User struct {
	Id       uint
	Name     string
	Age      int
	Created  time.Time
	Modified time.Time
}

type Role struct {
	Id   uint
	Name string
}

type UserRole struct {
	Id     uint
	UserId uint
	RoleId uint
}

type Question struct {
	Id     uint
	Title  string
	Desc   string
	UserId uint
}

type Answer struct {
	Id         uint
	QuestionId uint
	Content    string
	UserId     uint
}

func TestViz(t *testing.T) {
	tf, err := os.OpenFile("test.png", os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	defer tf.Close()

	artEnt := ToEntity(&Article{})
	userEnt := ToEntity(&User{})
	roleEnt := ToEntity(&Role{})
	userRoleEnt := ToEntity(&UserRole{})

	type args struct {
		edges map[string]map[string]Edge
		vs    []Entity
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				edges: map[string]map[string]Edge{
					artEnt.Name(): {
						userEnt.Name(): MakeEdge(
							artEnt,
							userEnt,
							"Author",
							"Id",
						),
					},
					userRoleEnt.Name(): {
						userEnt.Name(): MakeEdge(
							userRoleEnt,
							userEnt,
							"UserId",
							"Id",
						),
						roleEnt.Name(): MakeEdge(
							userRoleEnt,
							roleEnt,
							"RoleId",
							"Id",
						),
					},
				},
				vs: []Entity{
					artEnt,
					userEnt,
					roleEnt,
					userRoleEnt,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Viz(tf, tt.args.edges, tt.args.vs...); (err != nil) != tt.wantErr {
				t.Errorf("Viz() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGraphViz(t *testing.T) {
	tf, err := os.OpenFile("graphs.png", os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	defer tf.Close()

	artEnt := ToEntity(&Article{})
	userEnt := ToEntity(&User{})
	roleEnt := ToEntity(&Role{})
	userRoleEnt := ToEntity(&UserRole{})
	quesEnt := ToEntity(&Question{})
	ansEnt := ToEntity(&Answer{})

	type args struct {
		graphs []Graph
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				graphs: []Graph{
					{
						Edges: map[string]map[string]Edge{
							artEnt.Name(): {
								userEnt.Name(): MakeEdge(
									artEnt,
									userEnt,
									"Author",
									"Id",
								),
							},
							userRoleEnt.Name(): {
								userEnt.Name(): MakeEdge(
									userRoleEnt,
									userEnt,
									"UserId",
									"Id",
								),
								roleEnt.Name(): MakeEdge(
									userRoleEnt,
									roleEnt,
									"RoleId",
									"Id",
								),
							},
						},
						Entitys: []Entity{
							artEnt,
							userEnt,
							roleEnt,
							userRoleEnt,
						},
					},
					{
						Edges: map[string]map[string]Edge{
							ansEnt.Name(): {
								quesEnt.Name(): MakeEdge(
									ansEnt,
									quesEnt,
									"QuestionId",
									"Id",
								),
							},
						},
						Entitys: []Entity{
							quesEnt,
							ansEnt,
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GraphViz(tf, tt.args.graphs); (err != nil) != tt.wantErr {
				t.Errorf("Viz() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
