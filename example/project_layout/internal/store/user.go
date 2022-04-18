package store

import (
	"project_layout/model/db/user"

	"github.com/gocraft/dbr/v2"
	"github.com/pkg/errors"
)

type UserStore interface {
	ModName(id uint, name string) error
}

func NewUserStore(
	conn *dbr.Connection,
) UserStore {
	return &userImpl{
		conn: conn,
	}
}

type userImpl struct {
	conn *dbr.Connection
}

func (impl *userImpl) ModName(id uint, name string) error {
	sess := impl.conn.NewSession(nil)
	_, err := sess.Update(user.TableName).
		Where(
			"id = ?", id,
		).
		Set("name", name).
		Exec()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
