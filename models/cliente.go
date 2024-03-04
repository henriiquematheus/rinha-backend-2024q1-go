package models

// Cliente representa a estrutura de dados para um cliente
type Cliente struct {
	ID     int    `json:"id"`
	Limite int    `json:"limite"`
	Saldo  int    `json:"saldo"`
	Nome   string `json:"nome"`