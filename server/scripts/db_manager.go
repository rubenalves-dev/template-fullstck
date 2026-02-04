package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Expected 'create' or 'drop' command")
	}

	command := os.Args[1]
	dbName := "template-fullstack"
	if envDbName := os.Getenv("DB_NAME"); envDbName != "" {
		dbName = envDbName
	}

	// Connect to default 'postgres' database to perform administrative tasks
	connStr := os.Getenv("DB_ADMIN_CONN_STRING")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close()

	switch command {
	case "create":
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE \"%s\"", dbName))
		if err != nil {
			log.Fatalf("Failed to create database: %v", err)
		}
		fmt.Printf("Database %s created successfully\n", dbName)
	case "drop":
		_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS \"%s\"", dbName))
		if err != nil {
			log.Fatalf("Failed to drop database: %v", err)
		}
		fmt.Printf("Database %s dropped successfully\n", dbName)
	default:
		log.Fatalf("Unknown command: %s", command)
	}
}
