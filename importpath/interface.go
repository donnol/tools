package importpath

type IImportPath interface {
	GetByCurrentDir() (path string, err error)
	GetCurrentDirModFilePath() (modDir string, modPath string, err error)
	GetModFilePath(dir string) (modDir string, modPath string, err error)
}
