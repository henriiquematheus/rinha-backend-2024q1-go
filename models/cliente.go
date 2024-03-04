package models

// Cliente representa a estrutura de dados para um cliente
type Cliente struct {
    ID     int    `json:"id"`
    Limite int    `json:"limite"`
    Saldo  Saldo  `json:"saldo"` // Altere esta linha
    Nome   string `json:"nome"`
}