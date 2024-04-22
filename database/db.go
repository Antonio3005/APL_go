package database

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var DB_HOST string
var DB_NAME string
var DB_PASSWORD string
var DB_USER string

func Init() {
	// Carica la chiave segreta dalle variabili d'ambiente
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Errore nel caricamento delle variabili d'ambiente: %v", err)
	}
	DB_HOST = os.Getenv("DB_HOST")
	DB_NAME = os.Getenv("DB_NAME")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_USER = os.Getenv("DB_USER")
}

func RunDB() (*sql.DB, error) {
	Init()
	connString := fmt.Sprintf("%s:%s@tcp(%s)/%s", DB_USER, DB_PASSWORD, DB_HOST, DB_NAME)
	db, err := sql.Open("mysql", connString) //"root@tcp(localhost)/fly_scan")
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	// Controlla la connessione al database
	err = db.Ping()
	if err != nil {
		log.Fatal("Could not connect to the database:", err)
	}

	return db, err
}
