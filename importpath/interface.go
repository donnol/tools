package importpath

type IImportPath interface {
	GetByCurrentDir() (path string, err error)
	GetCurrentDirModFilePath() (modDir string, modPath string, err error)
	GetModFilePath(dir string) (modDir string, modPath string, err error)
	SplitImportPathWithType(importPathWithType string) (string, string)
}

type IImportPathMock interface {
	GetByCurrentDir() (path string, err error)
	GetCurrentDirModFilePath() (modDir string, modPath string, err error)
	GetModFilePath(dir string) (modDir string, modPath string, err error)
}
