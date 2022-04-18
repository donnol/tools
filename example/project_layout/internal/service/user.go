package service

import (
	"project_layout/internal/store"
)

type UserSrv interface {
	ModName(id uint, name string) error
}

func NewUserSrv(
	userStore store.UserStore,
) UserSrv {
	return &userImpl{
		userStore: userStore,
	}
}

type userImpl struct {
	userStore store.UserStore
}

func (impl *userImpl) ModName(id uint, name string) error {
	// 一些检查

	return impl.userStore.ModName(id, name)
}
