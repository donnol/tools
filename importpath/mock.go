package importpath

type IImportPathMock struct {
	GetByCurrentDirFunc func() (path string, err error)

	GetCurrentDirModFilePathFunc func() (modDir string, modPath string, err error)

	GetModFilePathFunc func(dir string) (modDir string, modPath string, err error)
}

var _ IImportPath = &IImportPathMock{}

func (*IImportPathMock) GetByCurrentDir() (path string, err error) {
	panic("Need to be implement!")
}

func (*IImportPathMock) GetCurrentDirModFilePath() (modDir string, modPath string, err error) {
	panic("Need to be implement!")
}

func (*IImportPathMock) GetModFilePath(dir string) (modDir string, modPath string, err error) {
	panic("Need to be implement!")
}
