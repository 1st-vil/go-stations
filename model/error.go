package model

import (
	"fmt"
)

type (
	ErrNotFound struct{
		RowID int64
	}
)

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("The row with id %v was not found", e.RowID)
}