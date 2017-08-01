package main

import (
	"database/sql"
	"fmt"
	"strconv"

	"errors"

	_ "github.com/lib/pq"
)

const (
	dbUser          = "postgres"
	dbPassword      = "postgres"
	dbName          = "hilltop_rain"
	dbInsertSitemap = "INSERT INTO sitemap (name, lat, lng) VALUES($1, $2, $3);"
	dbGetSitemapID  = "SELECT id FROM sitemap WHERE name = $1;"
)

var db = openDatabase()

func openDatabase() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	return db
}

func dbAddSite(site sitemap) (int, error) {
	_, err1 := strconv.ParseFloat(site.Lat, 64)
	_, err2 := strconv.ParseFloat(site.Lng, 64)
	if err1 != nil || err2 != nil {
		return 0, errors.New("Lat and Lng are not set")
	}

	Tx, err := db.Begin()
	checkErr(err)
	defer Tx.Rollback()

	stmt, err := Tx.Prepare(dbInsertSitemap)
	checkErr(err)

	res, err := stmt.Exec(site.Name, site.Lat, site.Lng)
	checkErr(err)

	_, err = res.RowsAffected()
	checkErr(err)

	err = Tx.Commit()
	checkErr(err)

	rows, err := db.Query(dbGetSitemapID, site.Name)
	checkErr(err)
	defer rows.Close()

	id := 0
	if rows.Next() {
		err = rows.Scan(&id)
		checkErr(err)
	}

	return id, nil
}

func dbAddRainRecord() {

}
