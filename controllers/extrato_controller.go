// controllers/extrato_controller.go
package controllers

import (
	"context"
	db "rinha-backend-2024q1-go/db"
	"rinha-backend-2024q1-go/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

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

	// Iniciar uma transação no banco de dados
	tx, err := conn.Begin(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao iniciar uma transação no banco de dados"})
	}
	defer tx.Rollback(c.Context())

	// Obter informações do cliente e transações
	cliente, err := ObterClientePorID(tx, clienteID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Cliente não encontrado"})
	}

	transacoes, err := ObterUltimasTransacoes(tx, clienteID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao obter as transações do cliente"})
	}

	// Criar resposta do extrato
	response := fiber.Map{
		"saldo": fiber.Map{
			"total":        cliente.Saldo,
			"data_extrato": cliente.Saldo.DataExtrato,
			"limite":       cliente.Limite,
		},
		"ultimas_transacoes": transacoes,
	}

	return c.Status(200).JSON(response)
}

func ObterUltimasTransacoes(tx pgx.Tx, clienteID int) ([]models.Transacao, error) {
	rows, err := tx.Query(context.Background(), `
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

func ObterClientePorIDExtrato(tx pgx.Tx, clienteID int) (*models.Cliente, error) {
	cliente := &models.Cliente{}
	err := tx.QueryRow(context.Background(), `
    SELECT id, nome, limite, saldo, data_extrato
    FROM clientes
    WHERE id = $1
`, clienteID).Scan(&cliente.ID, &cliente.Nome, &cliente.Saldo.Total, &cliente.Saldo.Limite, &cliente.Saldo.DataExtrato)
	if err != nil {
		return nil, err
	}

	return cliente, nil
}

func extrairClienteID(c *fiber.Ctx) (int, error) {
	id := c.Params("id")
	clienteID, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}
	return clienteID, nil
}
