package models

import (
	"time"
)

type Transacao struct {
	ClienteID   int       `json:"cliente_id"`
	Valor       int       `json:"valor"`
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	RealizadaEm time.Time `json:"realizada_em"`
}
