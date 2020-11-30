package importpath

type ImportPathMockMock struct {
	GetByCurrentDirFunc func() (path string, err error)

	GetCurrentDirModFilePathFunc func() (modDir string, modPath string, err error)

	GetModFilePathFunc func(dir string) (modDir string, modPath string, err error)
}

var _ IImportPathMock = &ImportPathMockMock{}

func (mockRecv *ImportPathMockMock) GetByCurrentDir() (path string, err error) {
	return mockRecv.GetByCurrentDirFunc()
}

func (mockRecv *ImportPathMockMock) GetCurrentDirModFilePath() (modDir string, modPath string, err error) {
	return mockRecv.GetCurrentDirModFilePathFunc()
}

func (mockRecv *ImportPathMockMock) GetModFilePath(dir string) (modDir string, modPath string, err error) {
	return mockRecv.GetModFilePathFunc(dir)
}

type ImportPathMock struct {
	GetByCurrentDirFunc func() (path string, err error)

	GetCurrentDirModFilePathFunc func() (modDir string, modPath string, err error)

	GetModFilePathFunc func(dir string) (modDir string, modPath string, err error)

	SplitImportPathWithTypeFunc func(importPathWithType string) (string, string)
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

func (mockRecv *ImportPathMock) SplitImportPathWithType(importPathWithType string) (string, string) {
	return mockRecv.SplitImportPathWithTypeFunc(importPathWithType)
}
