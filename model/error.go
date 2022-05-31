package model

import (
	"fmt"
)

type (
	ErrNotFound struct {
		RowIDs []int64
	}
)

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("The row with id(s) %v was not found", e.RowIDs)
}
