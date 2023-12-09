package erro

import (
	"errors"
)

var (
	ErrNotFound = errors.New("Item não encontrado")
	ErrUnmarshal = errors.New("Erro na conversão para JSON")
)
