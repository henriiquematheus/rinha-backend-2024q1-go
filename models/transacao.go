package models

import (
	"time"
)

// Transacao representa a estrutura de dados para uma transação
type Transacao struct {
	Valor      int    `json:"valor"`
	Tipo       string `json:"tipo"`
	Descricao  string `json:"descricao"`
	RealizadaEm time.Time `json:"realizada_em"`
}
