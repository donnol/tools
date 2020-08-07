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
