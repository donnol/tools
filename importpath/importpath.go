package importpath

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/mod/modfile"
)

type ImportPath struct {
}

const (
	gomodFileName = "go.mod"
)

func (p *ImportPath) GetCurrentDirModFilePath() (modDir, modPath string, err error) {
	// 获取目录
	dir, err := os.Getwd()
	if err != nil {
		return
	}

	modDir, modPath, err = getModFilePath(dir)
	if err != nil {
		return
	}

	return
}

func (p *ImportPath) GetModFilePath(dir string) (modDir, modPath string, err error) {
	return getModFilePath(dir)
}

func getModFilePath(dir string) (modDir, modPath string, err error) {
	// 找到go.mod所在目录
	var modfilePath string
	var modfileDir = dir
	for {
		modfilePath = filepath.Join(modfileDir, gomodFileName)

		_, err := os.Stat(modfilePath)
		if os.IsNotExist(err) {
			// 不存在，则继续往上层目录找
			tmpDir, file := filepath.Split(modfileDir)
			if file == "" {
				modfilePath = ""
				break
			}
			modfileDir = filepath.Clean(tmpDir)
			continue
		}

		break
	}
	if modfilePath == "" {
		return "", "", errors.Errorf("Can't find go mod file")
	}
	modDir = modfileDir
	modPath = modfilePath

	return
}

// GetCurrentDirPath 获取当前目录的包导入路径
func (p *ImportPath) GetByCurrentDir() (path string, err error) {
	return getPkgPathFromDir()
}

func getPkgPathFromDir() (pkgPath string, err error) {
	// 获取目录
	dir, err := os.Getwd()
	if err != nil {
		return
	}

	modfileDir, modfilePath, err := getModFilePath(dir)
	if err != nil {
		return
	}

	// 解析目录里的go.mod文件，获取模块名
	content, err := os.ReadFile(modfilePath)
	if err != nil {
		return
	}
	modPath := modfile.ModulePath(content)

	// 拿到go.mod所在目录和模块名，再结合当前目录信息，得到当前包路径
	// modPath + (modfileDir - dir)
	relPart := strings.ReplaceAll(dir, modfileDir, "")
	pkgPath = filepath.Join(modPath, relPart)

	return
}

func (p *ImportPath) SplitImportPathWithType(importPathWithType string) (string, string) {
	lastDotIndex := strings.LastIndex(importPathWithType, ".")
	if lastDotIndex == -1 || lastDotIndex == len(importPathWithType)-1 || lastDotIndex == 0 {
		panic(fmt.Errorf("路径有问题: %s", importPathWithType))
	}
	return importPathWithType[:lastDotIndex], importPathWithType[lastDotIndex+1:]
}

type Module struct {
	modfile.Module

	RelDir string
}

// FindAllModule 寻找目录下的所有module
func (p *ImportPath) FindAllModule(dir string) (mods []Module, err error) {
	if err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if errors.Is(err, os.ErrPermission) {
				return filepath.SkipDir
			}
			return err
		}

		if d.IsDir() {
			return nil
		}
		if d.Name() != "go.mod" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		f, err := modfile.Parse(path, content, nil)
		if err != nil {
			return err
		}

		mod := Module{
			Module: *f.Module,
			RelDir: filepath.Dir(path),
		}
		mods = append(mods, mod)

		return nil
	}); err != nil {
		return
	}
	return
}
