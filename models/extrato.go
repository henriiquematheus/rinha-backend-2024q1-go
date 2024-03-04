package models

import "time"

type Extrato struct {
	Saldo             Saldo       `json:"saldo"`
	UltimasTransacoes []Transacao `json:"ultimas_transacoes"`
}

type Saldo struct {
	Total       int       `json:"total"`
	DataExtrato time.Time `json:"data_extrato"`
	Limite      int       `json:"limite"`
}
