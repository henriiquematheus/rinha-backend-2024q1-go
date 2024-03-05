/*// controllers/extrato_controller.go
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
*/
// controllers/extrato_controller.go
package controllers

import (
	"context"
	"errors"
	connection "rinha-backend-2024q1-go/db"
	"rinha-backend-2024q1-go/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func extrairClienteID(c *fiber.Ctx) (int, error) {
	// Extrair o ID do cliente da URL
	clienteIDParam := c.Params("id")

	// Validar se o ID do cliente é um número inteiro válido
	id, err := strconv.Atoi(clienteIDParam)
	if err != nil {
		return 0, errors.New("ID do cliente inválido")
	}

	return id, nil

}

func ObterUltimasTransacoes(tx pgx.Tx, clienteID int) ([]models.Transacao, error) {
	// Consultar o banco de dados para obter as últimas transações do cliente
	rows, err := tx.Query(context.Background(), `
        SELECT id, valor, tipo, descricao, realizada_em
        FROM transacoes
        WHERE client_id = $1
        ORDER BY realizada_em DESC
        LIMIT 5  -- Ajuste o limite conforme necessário
    `, clienteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transacoes []models.Transacao

	for rows.Next() {
		var transacao models.Transacao
		err := rows.Scan(&transacao.ClienteID, &transacao.Valor, &transacao.Tipo, &transacao.Descricao, &transacao.RealizadaEm)
		if err != nil {
			return nil, err
		}
		transacoes = append(transacoes, transacao)
	}

	return transacoes, nil
}

func ObterClientePorIDExtrato(tx pgx.Tx, clienteID int) (*models.Cliente, error) {
	// Consultar o banco de dados para obter informações do cliente
	row := tx.QueryRow(context.Background(), `
        SELECT id, nome, limite, saldo
        FROM clientes
        WHERE id = $1
    `, clienteID)

	var cliente models.Cliente
	err := row.Scan(&cliente.ID, &cliente.Nome, &cliente.Limite, &cliente.Saldo.Total)
	if err != nil {
		return nil, err
	}

	return &cliente, nil
}

func ObterExtrato(c *fiber.Ctx) error {
	clienteID, err := extrairClienteID(c)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID do cliente inválido"})
	}

	// Obter uma conexão do pool do banco de dados
	conn, err := connection.GetPool().Acquire(c.Context())
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
	cliente, err := ObterClientePorIDExtrato(tx, clienteID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "Cliente não encontrado"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao obter informações do cliente"})
	}

	transacoes, err := ObterUltimasTransacoes(tx, clienteID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao obter as transações do cliente"})
	}

	// Calcular saldo total do cliente
	saldoTotal := calcularSaldoTotal(cliente.Saldo.Total, transacoes)

	// Criar resposta do extrato
	response := fiber.Map{
		"saldo": fiber.Map{
			"total":        saldoTotal,
			"data_extrato": time.Now().Format(time.RFC3339), // Formatar a data como string RFC3339
			"limite":       cliente.Limite,
		},
		"ultimas_transacoes": transacoes,
	}

	return c.Status(200).JSON(response)
}

func calcularSaldoTotal(saldoAtual int, transacoes []models.Transacao) int {
	// Calcular o saldo total do cliente considerando todas as transações
	for _, transacao := range transacoes {
		if transacao.Tipo == "c" {
			saldoAtual += transacao.Valor
		} else {
			saldoAtual -= transacao.Valor
		}
	}
	return saldoAtual
}
