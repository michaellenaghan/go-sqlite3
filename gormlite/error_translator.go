package gormlite

import (
	"errors"

	"gorm.io/gorm"

	"github.com/ncruces/go-sqlite3"
)

// Translate it will translate the error to native gorm errors.
func (_Dialector) Translate(err error) error {
	switch {
	case
		errors.Is(err, sqlite3.CONSTRAINT_UNIQUE),
		errors.Is(err, sqlite3.CONSTRAINT_PRIMARYKEY):
		return gorm.ErrDuplicatedKey
	case
		errors.Is(err, sqlite3.CONSTRAINT_FOREIGNKEY):
		return gorm.ErrForeignKeyViolated
	}
	return err
}
