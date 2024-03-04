package models

type Cliente struct {
	ID     int    `json:"id"`
	Limite int    `json:"limite"`
	Saldo  Saldo  `json:"saldo"`
	Nome   string `json:"nome"`
}
