package main

import (
	"APL_go/check"
	"APL_go/favourites"
	"APL_go/search"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()

	e.POST("/favourites", func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token mancante")
		}

		// Verifico che il token sia nel formato corretto (Bearer <token>)
		if len(token) < 7 || token[:7] != "Bearer " {
			return echo.NewHTTPError(http.StatusUnauthorized, "Formato del token non valido")
		}

		//rimuovo il prefisso "Bearer ")
		token = token[7:]
		fmt.Println(token)
		user, err := check.CheckToken(token)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		var req favourites.FavouritesRequest
		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Errore durante la lettura del corpo della richiesta JSON")
		}
		req.User = user

		return c.JSON(http.StatusOK, favourites.FavouritesHandler(req))
	})

	e.POST("/selectfav", func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token mancante")
		}

		// Verifico che il token sia nel formato corretto (Bearer <token>)
		if len(token) < 7 || token[:7] != "Bearer " {
			return echo.NewHTTPError(http.StatusUnauthorized, "Formato del token non valido")
		}

		//rimuovo il prefisso "Bearer ")
		token = token[7:]

		user, err := check.CheckToken(token)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		//var req favourites.FavouritesRequest

		//if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		//	return echo.NewHTTPError(http.StatusBadRequest, "Errore durante la lettura del corpo della richiesta JSON")
		//}
		//user := req.User

		return c.JSON(http.StatusOK, favourites.FavouritesSelect(user))
	})

	e.POST("/deletefav", func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token mancante")
		}

		// Verifico che il token sia nel formato corretto (Bearer <token>)
		if len(token) < 7 || token[:7] != "Bearer " {
			return echo.NewHTTPError(http.StatusUnauthorized, "Formato del token non valido")
		}

		//rimuovo il prefisso "Bearer ")
		token = token[7:]

		user, err := check.CheckToken(token)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
		var req favourites.FavouritesRequest

		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Errore durante la lettura del corpo della richiesta JSON")
		}
		req.User = user
		return c.JSON(http.StatusOK, favourites.FavouritesDelete(req))
	})

	e.POST("/search", func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token mancante")
		}

		// Verifico che il token sia nel formato corretto (Bearer <token>)
		if len(token) < 7 || token[:7] != "Bearer " {
			return echo.NewHTTPError(http.StatusUnauthorized, "Formato del token non valido")
		}

		//rimuovo il prefisso "Bearer ")
		token = token[7:]

		user, err := check.CheckToken(token)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		fmt.Println(user)

		var req search.SearchRequest
		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Errore durante la lettura del corpo della richiesta JSON")
		}

		flights, err := search.SearchHandler(req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Errore durante il recupero dei voli")
		}

		return c.JSON(http.StatusOK, flights)
	})

	e.POST("/show_pop", func(c echo.Context) error {
		popular, err := search.ShowPopular()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Errore durante il recupero dei voli popolari")
		}

		return c.JSON(http.StatusOK, popular)
	})

	// Avvio del server
	e.Logger.Fatal(e.Start(":8080"))
}
