package search

import (
	"APL_go/database"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

var apiKey string

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Errore nel caricamento delle variabili d'ambiente: %v", err)
	}
	apiKey = os.Getenv("API_KEY")
}

type SearchRequest struct {
	CityFrom   string `json:"city_from"`
	CityTo     string `json:"city_to"`
	DataFrom   string `json:"date_from"`
	DataTo     string `json:"date_to"`
	ReturnFrom string `json:"return_from"`
	ReturnTo   string `json:"return_to"`
	PriceMin   string `json:"price_min"`
	PriceMax   string `json:"price_max"`
}

func SearchHandler(c SearchRequest) (map[string]interface{}, error) {
	//user, errore := check.CheckToken(c.Token)
	Init()
	city_from := c.CityFrom
	city_to := c.CityTo
	data_from := c.DataFrom
	data_to := c.DataTo
	return_from := c.ReturnFrom
	return_to := c.ReturnTo
	price_min := c.PriceMin
	price_max := c.PriceMax

	db, err := database.RunDB()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS research (
		id INT AUTO_INCREMENT PRIMARY KEY,
		city_from VARCHAR(255),
		city_to VARCHAR(255)
    )`)

	if err != nil {
		log.Fatal("Errore durante la creazione della tabella:", err)
	}

	query := "INSERT INTO research (city_from, city_to) VALUES (?, ?)"

	_, err = db.Exec(query, city_from, city_to)

	if err != nil {
		log.Fatal("Errore durante l'inserimento nella tabella:", err)
	}

	iata_from, err := getIata(city_from)
	if err != nil {
		return nil, err
	}
	iata_to, err := getIata(city_to)
	if err != nil {
		return nil, err
	}
	flights, err := getFlights(iata_from, iata_to, data_from, data_to, return_from, return_to, price_min, price_max)

	defer db.Close()

	return flights, nil
}

func getIata(city string) (string, error) {
	const apiURL = "https://api.tequila.kiwi.com/locations/query"

	params := url.Values{}
	params.Set("term", city)
	params.Set("locale", "it-IT")
	params.Set("location_types", "city")
	params.Set("limit", "10")
	params.Set("active_only", "true")

	req, err := http.NewRequest("GET", apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("apikey", apiKey) // Assicurati di definire apiKey

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("errore nella richiesta. Codice di stato: %d", resp.StatusCode)
	}

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", fmt.Errorf("errore nella decodifica della risposta JSON: %v", err)
	}

	locations, ok := data["locations"].([]interface{})
	if !ok || len(locations) == 0 {
		return "", fmt.Errorf("array 'locations' vuoto o assente nella risposta dell'API")
	}

	iata, ok := locations[0].(map[string]interface{})["code"].(string)
	if !ok {
		return "", fmt.Errorf("impossibile estrarre il codice IATA dalla risposta")
	}

	return iata, nil
}

func getFlights(iataFrom, iataTo, dateFrom, dateTo, returnFrom, returnTo, priceFrom, priceTo string) (map[string]interface{}, error) {
	const apiURL = "https://api.tequila.kiwi.com/v2/search"

	params := url.Values{}
	params.Set("fly_from", iataFrom)
	params.Set("fly_to", iataTo)
	params.Set("date_from", dateFrom)
	params.Set("date_to", dateTo)
	params.Set("return_from", returnFrom)
	params.Set("return_to", returnTo)
	params.Set("adults", "1")
	params.Set("adult_hand_bag", "1")
	params.Set("partner_market", "it")
	params.Set("price_from", priceFrom)
	params.Set("price_to", priceTo)
	params.Set("vehicle_type", "aircraft")
	params.Set("sort", "price")
	params.Set("limit", "2")
	params.Set("locale", "it")

	req, err := http.NewRequest("GET", apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("apikey", apiKey) // Assicurati di definire apiKey

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("errore nella richiesta. Codice di stato: %d", resp.StatusCode)
	}

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("errore nella decodifica della risposta JSON: %v", err)
	}

	return data, nil
}

func ShowPopular() ([]map[string]interface{}, error) {
	Init()
	db, err := database.RunDB()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS research (
		id INT AUTO_INCREMENT PRIMARY KEY,
		city_from VARCHAR(255),
		city_to VARCHAR(255)
    )`)

	if err != nil {
		log.Fatal("Errore durante la creazione della tabella:", err)
	}

	rows, err := db.Query("SELECT city_from, city_to FROM research")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Creare una mappa per contare quante volte ogni coppia città appare nelle ricerche
	counts := make(map[string]int)

	// Scorrere i risultati della query e contare le occorrenze di ogni coppia città
	for rows.Next() {
		var cityFrom, cityTo string
		if err := rows.Scan(&cityFrom, &cityTo); err != nil {
			return nil, err
		}
		// Costruire una stringa univoca per rappresentare la coppia città
		key := cityFrom + "|" + cityTo
		// Incrementare il conteggio per questa coppia città
		counts[key]++
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Creare una slice per memorizzare le ricerche più popolari
	var popularSearches []map[string]interface{}

	// Scorrere la mappa dei conteggi e aggiungere le ricerche più popolari alla slice
	for key, count := range counts {
		cityFrom, cityTo := parseKey(key)
		popularSearch := map[string]interface{}{
			"city_from": cityFrom,
			"city_to":   cityTo,
			"count":     count,
		}
		popularSearches = append(popularSearches, popularSearch)
	}

	// Ordinare le ricerche più popolari per conteggio (in ordine decrescente)
	sort.Slice(popularSearches, func(i, j int) bool {
		return popularSearches[i]["count"].(int) > popularSearches[j]["count"].(int)
	})

	// Limitare l'elenco alle prime 10 ricerche più popolari
	if len(popularSearches) > 4 {
		popularSearches = popularSearches[:4]
	}

	defer db.Close()

	return popularSearches, nil
}

// Funzione di utilità per estrarre le città da una chiave composta
func parseKey(key string) (string, string) {
	parts := strings.Split(key, "|")
	return parts[0], parts[1]
}
