package importpath

import (
	"os"
	"testing"
)

func TestGetCurrentDirModFilePath(t *testing.T) {
	ip := &ImportPath{}
	modDir, modPath, err := ip.GetCurrentDirModFilePath()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("modDir: %s, modPath: %s\n", modDir, modPath)

	// 用os.Stat校验目录是否存在
	info, err := os.Stat(modDir)
	if os.IsNotExist(err) {
		t.Fatalf("modDir is not exist: %+v\n", info)
	}

	info, err = os.Stat(modPath)
	if os.IsNotExist(err) {
		t.Fatalf("modDir is not exist: %+v\n", info)
	}
}

func TestGetModFilePath(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		return
	}

	ip := &ImportPath{}
	modDir, modPath, err := ip.GetModFilePath(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("modDir: %s, modPath: %s\n", modDir, modPath)

	// 用os.Stat校验目录是否存在
	info, err := os.Stat(modDir)
	if os.IsNotExist(err) {
		t.Fatalf("modDir is not exist: %+v\n", info)
	}

	info, err = os.Stat(modPath)
	if os.IsNotExist(err) {
		t.Fatalf("modDir is not exist: %+v\n", info)
	}
}

func TestGetCurrentDirPath(t *testing.T) {
	ip := &ImportPath{}
	path, err := ip.GetByCurrentDir()
	if err != nil {
		t.Fatal(err)
	}
	want := "github.com/donnol/tools/importpath"
	if path != want {
		t.Fatalf("Bad result: %v != %v\n", path, want)
	}
}

func TestImportPath_FindAllModule(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name     string
		p        *ImportPath
		args     args
		wantMods []Module
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			p:    &ImportPath{},
			args: args{
				dir: "../../",
			},
			wantMods: nil,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ImportPath{}
			gotMods, err := p.FindAllModule(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImportPath.FindAllModule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(gotMods) == 0 {
				t.Errorf("ImportPath.FindAllModule() = %v", gotMods)
			}
		})
	}
}
