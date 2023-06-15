package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// #region city
type City struct {
	ID          int    `json:"ID,omitempty" db:"ID"`
	Name        string `json:"name,omitempty" db:"Name"`
	CountryCode string `json:"countryCode,omitempty"  db:"CountryCode"`
	District    string `json:"district,omitempty"  db:"District"`
	Population  int    `json:"population,omitempty"  db:"Population"`
}

type PopulationInfo struct {
	Population int     `db:"Population"`
	Percent    float64 `db:"Percent"`
}

// #endregion city
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

	db, err := sqlx.Open("mysql", conf.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("conntected")
	// #region get

	cityName := os.Args[1]

	// var city City
	var populationInfo PopulationInfo
	// err = db.Get(&city, "SELECT * FROM city WHERE Name = ?", "Tokyo")
	err = db.Get(&populationInfo, "SELECT city.Population AS Population, city.Population * 100 / country.Population AS Percent FROM city JOIN country ON city.CountryCode = country.Code WHERE city.Name = ?", cityName)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("no such city Name = '%s'\n", cityName)
	} else if err != nil {
		log.Fatalf("DB Error: %s\n", err)
	}
	// #endregion get
	fmt.Printf("%sの人口は%d人でその国の%f%%です\n", cityName, populationInfo.Population, populationInfo.Percent)
}
