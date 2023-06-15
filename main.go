package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type City struct {
	ID          int    `json:"id,omitempty"  db:"ID"`
	Name        string `json:"name,omitempty"  db:"Name"`
	CountryCode string `json:"countryCode,omitempty"  db:"CountryCode"`
	District    string `json:"district,omitempty"  db:"District"`
	Population  int    `json:"population,omitempty"  db:"Population"`
}

var (
	db *sqlx.DB
)

func main() {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}

	conf := mysql.Config{
		User:      os.Getenv("DB_USERNAME"),
		Passwd:    os.Getenv("DB_PASSWORD"),
		Net:       "tcp",
		Addr:      os.Getenv("DB_HOSTNAME") + ":" + os.Getenv("DB_PORT"),
		DBName:    os.Getenv("DB_DATABASE"),
		ParseTime: true,
		Collation: "utf8mb4_unicode_ci",
		Loc:       jst,
	}

	_db, err := sqlx.Open("mysql", conf.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("conntected")
	db = _db

	e := echo.New()

	e.GET("/cities/:cityName", getCityInfoHandler)
	e.POST("/cities", postCityHandler)

	e.Start(":8080")
}

func getCityInfoHandler(c echo.Context) error {
	cityName := c.Param("cityName")
	fmt.Println(cityName)

	var city City
	if err := db.Get(&city, "SELECT * FROM city WHERE Name=?", cityName); errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("No such city Name = %s", cityName))
	} else if err != nil {
		log.Fatalf("DB Error: %s", err)
	}

	return c.JSON(http.StatusOK, city)
}

type postCityRequest struct {
	Name        string `json:"name,omitempty" db:"Name"`
	CountryCode string `json:"countryCode,omitempty" db:"CountryCode"`
	Population  int    `json:"population,omitempty" db:"Population"`
	District    string `json:"district,omitempty" db:"District"`
}

type postCityResponse struct {
	Id          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	CountryCode string `json:"countryCode,omitempty"`
	Population  int    `json:"population,omitempty"`
	District    string `json:"district,omitempty"`
}

func postCityHandler(c echo.Context) error {
	var req postCityRequest
	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	r, err := db.Exec("INSERT INTO city (Name, CountryCode, Population, District) VALUES (?, ?, ?, ?)", req.Name, req.CountryCode, req.Population, req.District)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	id, err := r.LastInsertId()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, postCityResponse{
		Id:          int(id),
		Name:        req.Name,
		District:    req.District,
		CountryCode: req.CountryCode,
		Population:  req.Population,
	})
}
