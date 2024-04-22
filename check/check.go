package check

import (
	"APL_go/database"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Token struct {
	Username       string  `json:"username"`
	ExpirationTime float64 `json:"expirationTime"`
	jwt.StandardClaims
}

var SECRET_KEY string //= "aplproject" //os.Getenv("SECRET_KEY")

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Errore nel caricamento delle variabili d'ambiente: %v", err)
	}
	SECRET_KEY = os.Getenv("SECRET_KEY")
}

func CheckToken(tokenString string) (string, error) {
	// Parse del token JWT utilizzando la chiave segreta
	Init()

	db, err := database.RunDB()
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS token_revocated (
		id INT AUTO_INCREMENT PRIMARY KEY,
		token VARCHAR(255)
    )`)

	if err != nil {
		log.Fatal("Errore durante la creazione della tabella:", err)
	}
	query := "SELECT token FROM token_revoked WHERE token=?"
	row := db.QueryRow(query, tokenString)

	var revokedToken string
	err = row.Scan(&revokedToken)
	if err == nil {
		// Se il token è stato trovato nella tabella dei token revocati, restituisci un errore
		return "", errors.New("Il token è stato revocato")
	} else if err != nil && err != sql.ErrNoRows {
		// Se si verifica un errore diverso da "nessuna riga restituita", restituisci l'errore
		log.Printf("Errore durante l'esecuzione della query: %v", err)
		return "", err
	}

	token, err := jwt.ParseWithClaims(tokenString, &Token{}, func(token *jwt.Token) (interface{}, error) {
		// Verifica il metodo di firma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Metodo di firma non valido")
		}
		return []byte(SECRET_KEY), nil
	})
	fmt.Println(token)
	if err != nil {
		fmt.Println(err.Error())
		return "", errors.New("Token non valido")
	}

	// Verifica se il token è valido e non scaduto
	if claims, ok := token.Claims.(*Token); ok && token.Valid {
		// Verifica se il token è scaduto
		fmt.Println(claims.Username, claims.ExpirationTime)
		if time.Now().Unix() > (int64(claims.ExpirationTime)) {
			fmt.Println(time.Now().Unix())
			return "", errors.New("Token scaduto")
		}
		// Token valido, restituisci l'username
		return claims.Username, nil
	}
	// Token non valido
	return "", errors.New("Token non valido")
}
