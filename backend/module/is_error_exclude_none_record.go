package module

import (
	"errors"

	"gorm.io/gorm"
)

func IsErrorExcludeNoneRecord(err error) bool {
	return err == nil && errors.Is(err, gorm.ErrRecordNotFound)
}
