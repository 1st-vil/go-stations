package model

import (
	"fmt"
)

type (
	ErrNotFound struct{
		FileName string
	}
)

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("The file %s was not found")
}