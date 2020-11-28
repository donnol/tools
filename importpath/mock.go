package importpath

type ImportPathMock struct {
	GetByCurrentDirFunc func() (path string, err error)

	GetCurrentDirModFilePathFunc func() (modDir string, modPath string, err error)

	GetModFilePathFunc func(dir string) (modDir string, modPath string, err error)
}

var _ IImportPath = &ImportPathMock{}

func (mockRecv *ImportPathMock) GetByCurrentDir() (path string, err error) {
	return mockRecv.GetByCurrentDirFunc()
}

func (mockRecv *ImportPathMock) GetCurrentDirModFilePath() (modDir string, modPath string, err error) {
	return mockRecv.GetCurrentDirModFilePathFunc()
}

func (mockRecv *ImportPathMock) GetModFilePath(dir string) (modDir string, modPath string, err error) {
	return mockRecv.GetModFilePathFunc(dir)
}
