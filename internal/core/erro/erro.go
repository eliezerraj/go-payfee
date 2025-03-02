package erro

import (
	"errors"
)

var (
	ErrNotFound 		= errors.New("item not found")
	ErrInsert 			= errors.New("insert data error")
	ErrUnmarshal 		= errors.New("unmarshal json error")
	ErrServer		 	= errors.New("server identified error")
)