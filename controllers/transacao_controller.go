// controllers/transacao_controller.go
package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"rinha-backend-2024q1-go/models"
)

// ProcessarTransacao processa uma transação
func ProcessarTransacao(c *fiber.Ctx) error {
	var transacao models.Transacao

	// Parse do corpo da requisição para a estrutura de dados Transacao
	if err := c.BodyParser(&transacao); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Erro ao analisar o corpo da requisição",
		})
	}

	// Obter uma conexão do pool do banco de dados
	conn, err := db.GetPool().Acquire(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Erro ao obter uma conexão do pool do banco de dados",
		})
	}
	defer conn.Release()

	// Iniciar uma transação no banco de dados
	tx, err := conn.Begin(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Erro ao iniciar uma transação no banco de dados",
		})
	}
	defer tx.Rollback(c.Context())

	// Executar a lógica de transação
	cliente, err := AtualizarSaldo(tx, transacao.ClienteID, transacao.Valor, transacao.Tipo)
	if err != nil {
		// Tratar erros específicos aqui, se necessário
		switch err {
		case ErrClienteNaoEncontrado:
			return c.Status(404).JSON(fiber.Map{"error": "Cliente não encontrado"})
		case ErrLimiteExcedido:
			return c.Status(422).JSON(fiber.Map{"error": "Transação de débito excede o limite do cliente"})
		default:
			return c.Status(500).JSON(fiber.Map{"error": "Erro ao processar a transação"})
		}
	}

	// Commit da transação
	if err := tx.Commit(c.Context()); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Erro ao fazer commit da transação no banco de dados",
		})
	}

	// Retorne a resposta com os dados atualizados
	return c.Status(200).JSON(fiber.Map{
		"limite": cliente.Limite,
		"saldo":  cliente.Saldo,
	})
}

var (
	// ErrClienteNaoEncontrado é um erro indicando que o cliente não foi encontrado.
	ErrClienteNaoEncontrado = fmt.Errorf("Cliente não encontrado")
	
	// ErrLimiteExcedido é um erro indicando que o limite do cliente foi excedido durante uma transação de débito.
	ErrLimiteExcedido = fmt.Errorf("Transação de débito excede o limite do cliente")
)

// AtualizarSaldo atualiza o saldo e o limite do cliente no banco de dados
func AtualizarSaldo(tx pgx.Tx, clienteID int, valor int, tipo string) (*models.Cliente, error) {
	// Obter informações atuais do cliente
	cliente, err := ObterClientePorID(tx, clienteID)
	if err != nil {
		return nil, ErrClienteNaoEncontrado
	}

	// Validar a transação com base no tipo (crédito ou débito)
	if tipo == "d" && cliente.Saldo-valor < -cliente.Limite {
		return nil, ErrLimiteExcedido
	}

	// Atualizar o saldo e limite com base no tipo de transação
	if tipo == "c" {
		cliente.Saldo += valor
	} else if tipo == "d" {
		cliente.Saldo -= valor
		cliente.Limite -= valor
	}

	// Atualizar o cliente no banco de dados
	_, err = tx.Exec(tx.Context(), `
		UPDATE clientes
		SET saldo = $1, limite = $2
		WHERE id = $3
	`, cliente.Saldo, cliente.Limite, clienteID)

	if err != nil {
		return nil, err
	}

	return cliente, nil
}

// ObterClientePorID obtém as informações do cliente do banco de dados
func ObterClientePorID(tx pgx.Tx, clienteID int) (*models.Cliente, error) {
	cliente := &models.Cliente{}
	err := tx.QueryRow(tx.Context(), `
		SELECT id, nome, limite, saldo
		FROM clientes
		WHERE id = $1
	`, clienteID).Scan(&cliente.ID, &cliente.Nome, &cliente.Limite, &cliente.Saldo)

	if err != nil {
		return nil, err
	}

	return cliente, nil
}
