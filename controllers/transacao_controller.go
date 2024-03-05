// controllers/transacao_controller.go
package controllers

import (
	"context"
	db "rinha-backend-2024q1-go/db"
	"rinha-backend-2024q1-go/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func CriarTransacao(c *fiber.Ctx) error {
	// 1. Parse da Requisição
	clienteID, err := extrairClienteID(c)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID do cliente inválido"})
	}

	var transacao models.Transacao
	if err := c.BodyParser(&transacao); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Erro ao parsear o corpo da requisição"})
	}

	// 2. Validações de Campos
	errMsg := validarCamposTransacao(transacao)
	if errMsg != "" {
		return c.Status(422).JSON(fiber.Map{"error": errMsg})
	}

	// 3. Verificação de Cliente
	conn, err := db.GetPool().Acquire(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao obter uma conexão do pool do banco de dados"})
	}
	defer conn.Release()

	tx, err := conn.Begin(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao iniciar uma transação no banco de dados"})
	}
	defer tx.Rollback(c.Context())

	cliente, err := ObterClientePorIDExtrato(tx, clienteID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "Cliente não encontrado"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao obter informações do cliente"})
	}

	// 4. Realização da Transação
	novoSaldo := calcularNovoSaldo(cliente.Saldo.Total, transacao)
	if novoSaldo < -cliente.Saldo.Limite {
		return c.Status(422).JSON(fiber.Map{"error": "Transação resulta em saldo inconsistente"})
	}

	// Atualizar saldo do cliente no banco de dados
	if err := AtualizarSaldoCliente(tx, clienteID, novoSaldo); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao atualizar saldo do cliente"})
	}

	// Registrar a transação no banco de dados
	if err := RegistrarTransacao(tx, clienteID, transacao); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao registrar a transação"})
	}

	tx.Commit(c.Context())

	// 5. Resposta
	response := fiber.Map{
		"limite": cliente.Saldo.Limite,
		"saldo":  novoSaldo,
	}

	return c.Status(200).JSON(response)
}

func validarCamposTransacao(transacao models.Transacao) string {
	// Validar campos conforme as regras fornecidas na checklist
	// Retornar string de mensagem de erro em caso de campos inválidos
	var errMsg string

	// 1. Verificar se o campo valor é um número inteiro positivo.
	if transacao.Valor <= 0 {
		errMsg = "O valor deve ser um número inteiro positivo"
	}

	// 2. Validar se o campo tipo contém apenas "c" ou "d".
	if transacao.Tipo != "c" && transacao.Tipo != "d" {
		errMsg = "O tipo deve ser 'c' para crédito ou 'd' para débito"
	}

	// 3. Validar se o campo descricao tem entre 1 e 10 caracteres.
	descricaoLen := len(transacao.Descricao)
	if descricaoLen < 1 || descricaoLen > 10 {
		errMsg = "A descrição deve ter entre 1 e 10 caracteres"
	}

	return errMsg
}

func calcularNovoSaldo(saldoAtual int, transacao models.Transacao) int {
	// Calcular o novo saldo com base na transação (crédito ou débito)
	if transacao.Tipo == "c" {
		return saldoAtual + transacao.Valor
	} else {
		return saldoAtual - transacao.Valor
	}
}

func AtualizarSaldoCliente(tx pgx.Tx, clienteID int, novoSaldo int) error {
	// Atualizar o saldo do cliente no banco de dados
	_, err := tx.Exec(context.Background(), `
		UPDATE clientes
		SET saldo = $1
		WHERE id = $2
	`, novoSaldo, clienteID)

	return err
}

func RegistrarTransacao(tx pgx.Tx, clienteID int, transacao models.Transacao) error {
	// Registrar a transação no banco de dados
	_, err := tx.Exec(context.Background(), `
		INSERT INTO transacoes (client_id, valor, tipo, descricao, realizada_em)
		VALUES ($1, $2, $3, $4, $5)
	`, clienteID, transacao.Valor, transacao.Tipo, transacao.Descricao, time.Now())

	return err
}
