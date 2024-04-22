package favourites

import (
	"APL_go/database"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type FavouritesRequest struct {
	CityFrom   string `json:"city_from"`
	CityTo     string `json:"city_to"`
	DataFrom   string `json:"date_from"`
	ReturnData string `json:"return_from"`
	Price      string `json:"price"`
	User       string `json:"user"`
}

func FavouritesHandler(c FavouritesRequest) string {
	city_from := c.CityFrom
	city_to := c.CityTo
	data_from := c.DataFrom
	return_from := c.ReturnData
	price := c.Price
	user := c.User

	db, err := database.RunDB()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS favourites (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user VARCHAR(255),
		city_from VARCHAR(255),
		city_to VARCHAR(255),
		date_from VARCHAR(255),
		return_from VARCHAR(255),
		price VARCHAR(255)
    )`)

	if err != nil {
		log.Fatal("Errore durante la creazione della tabella:", err)
	}

	var count int
	checkQuery := "SELECT COUNT(*) FROM favourites WHERE user = ? AND city_from = ? AND city_to = ? AND date_from = ? AND return_from = ? AND price = ?"
	err = db.QueryRow(checkQuery, user, city_from, city_to, data_from, return_from, price).Scan(&count)

	if err != nil {
		log.Fatal("Errore durante il controllo del record:", err)
	}

	if count > 0 {
		// Se il record esiste, restituisci un messaggio che indica che è già presente
		return "Questo elemento è già nei preferiti"
	}

	query := "INSERT INTO favourites (user, city_from, city_to, date_from, return_from, price) VALUES (?, ?, ?, ?, ?, ?)"

	_, err = db.Exec(query, user, city_from, city_to, data_from, return_from, price)

	if err != nil {
		log.Fatal("Errore durante l'inserimento nella tabella:", err)
	}

	defer db.Close()

	return "Aggiunto ai preferiti"
}

func FavouritesSelect(user string) []FavouritesRequest {
	var favourites []FavouritesRequest
	db, err := database.RunDB()

	query := "SELECT city_from, city_to, date_from, return_from, price FROM favourites WHERE user = ?"

	rows, err := db.Query(query, user)
	if err != nil {
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var favourite FavouritesRequest
		err := rows.Scan(
			&favourite.CityFrom,
			&favourite.CityTo,
			&favourite.DataFrom,
			&favourite.ReturnData,
			&favourite.Price,
		)
		if err != nil {
			return nil
		}

		favourites = append(favourites, favourite)
	}

	if err := rows.Err(); err != nil {
		return nil
	}

	defer db.Close()

	return favourites
}

func FavouritesDelete(c FavouritesRequest) string {
	price := c.Price
	city_from := c.CityFrom
	city_to := c.CityTo
	data_from := c.DataFrom
	return_data := c.ReturnData
	user := c.User

	db, err := database.RunDB()

	query := "DELETE FROM favourites WHERE " +
		"user = ? AND " +
		"city_from = ? AND " +
		"city_to = ? AND " +
		"date_from = ? AND " +
		"return_from = ? AND " +
		"price = ?"

	_, err = db.Exec(query, user, city_from, city_to, data_from, return_data, price)
	if err != nil {
		log.Fatal("Errore durante l'eliminazione dalla tabella:", err)
	}

	defer db.Close()

	return "Eliminato dai preferiti"
}
