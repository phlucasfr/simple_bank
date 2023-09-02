package main

import (
	"database/sql"
	"log"
	"picpay_simplificado/api"
	db "picpay_simplificado/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/picpay_simplificado?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.StartServer(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
