// controllers/extrato_controller.go
package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/henriiquematheus/rinha-backend-2024q1-go/db"
	"github.com/henriiquematheus/rinha-backend-2024q1-go/models"
)

// ObterExtrato obtém o extrato do cliente
func ObterExtrato(c *fiber.Ctx) error {
	clienteID, err := extrairClienteID(c)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID do cliente inválido"})
	}

	// Obter uma conexão do pool do banco de dados
	conn, err := db.GetPool().Acquire(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao obter uma conexão do pool do banco de dados"})
	}
	defer conn.Release()

	// Obter informações do cliente e transações
	cliente, err := ObterClientePorID(conn, clienteID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Cliente não encontrado"})
	}

	transacoes, err := ObterUltimasTransacoes(conn, clienteID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao obter as transações do cliente"})
	}

	// Criar resposta do extrato
	response := fiber.Map{
		"saldo": fiber.Map{
			"total":       cliente.Saldo,
			"data_extrato": time.Now().UTC(),
			"limite":      cliente.Limite,
		},
		"ultimas_transacoes": transacoes,
	}

	return c.Status(200).JSON(response)
}

// ObterUltimasTransacoes obtém as últimas transações do cliente
func ObterUltimasTransacoes(tx pgx.Tx, clienteID int) ([]models.Transacao, error) {
	rows, err := tx.Query(tx.Context(), `
		SELECT valor, tipo, descricao, realizada_em
		FROM transacoes
		WHERE client_id = $1
		ORDER BY realizada_em DESC
		LIMIT 10
	`, clienteID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transacoes []models.Transacao
	for rows.Next() {
		var transacao models.Transacao
		err := rows.Scan(&transacao.Valor, &transacao.Tipo, &transacao.Descricao, &transacao.RealizadaEm)
		if err != nil {
			return nil, err
		}
		transacoes = append(transacoes, transacao)
	}

	return transacoes, nil
}
