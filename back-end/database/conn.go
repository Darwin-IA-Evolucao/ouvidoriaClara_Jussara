package database

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func ConnectDB() (*sqlx.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Validar variáveis obrigatórias
	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		panic(fmt.Sprintf("Variáveis de ambiente do banco não configuradas! Verifique se o .env existe e está preenchido.\nDB_HOST=%s, DB_PORT=%s, DB_USER=%s, DB_NAME=%s",
			host, port, user, dbname))
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to the database! " + dbname)
	return db, nil
}
