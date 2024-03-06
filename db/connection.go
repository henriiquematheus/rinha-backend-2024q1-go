// connection.go

package connection

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func Init(connectionString string) error {
	log.Println("Iniciando a inicialização da conexão com o banco de dados...")
	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return fmt.Errorf("Erro ao fazer parse da configuração: %v", err)
	}

	// Adicione logs para imprimir as configurações de conexão
	log.Printf("Configurações de Conexão: %+v", config)

	pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("Erro ao criar o pool de conexões: %v", err)
	}

	// Adicione logs para indicar que o pool foi criado com sucesso
	log.Println("Pool de Conexões criado com sucesso")

	return nil
}

func GetPool() *pgxpool.Pool {
	if pool == nil {
		log.Println("AVISO: Pool de conexões não inicializado. Certifique-se de chamar connection.Init() primeiro.")
	}
	return pool
}
