// main.go

package main

import (
	"log"

	connection "rinha-backend-2024q1-go/db"
)

func main() {
	// Substitua a string de conexão conforme necessário
	connectionString := "postgres://rinha:123@db:5432/rinha"

	err := connection.Init(connectionString)
	if err != nil {
		log.Fatal(err)
	}

	defer connection.GetPool().Close()
}
