package modules

import (
	"database/sql"
	"fmt"

	"github.com/Nosolok/vt/vtApi"
	_ "github.com/mattn/go-sqlite3"
)

type FileInfo struct {
	Id        int
	Sha256    string
	Scan_date string
	Positives int
	Total     int
}

/*
	create database structure
*/
func InitDatabase(db *sql.DB) {
	sqlStmt := `
		create table file_checks (
			id integer not null primary key autoincrement,
			sha256 string unique,
			scan_date string,
			positives int,
			total int
			);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		fmt.Println("%q: %s\n", err, sqlStmt)
	}
}

/*
	check hash in database
*/
func CheckHash(db *sql.DB, hash string) FileInfo {
	stmt := `
		select *
		from file_checks
		where sha256 = ?
	`

	var row FileInfo
	err := db.QueryRow(stmt, hash).Scan(&row.Id, &row.Sha256, &row.Scan_date, &row.Positives, &row.Total)

	switch err {
	case sql.ErrNoRows:
		return row
	case nil:
		return row
	default:

	}

	return row
}

func StoreCheck(db *sql.DB, report vtApi.FileReportResponse) {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("%q: %s\n", err, tx)
	}

	stmt, err := tx.Prepare(`
		insert into file_checks(sha256, scan_date, positives, total)
		values(?, ?, ?, ?)
	`)
	if err != nil {
		fmt.Println("%q: %s\n", err, stmt)
	}
	defer stmt.Close()

	result, err := stmt.Exec(report.Sha256, report.ScanDate, report.Positives, report.Total)
	if err != nil {
		fmt.Println("%q: %s\n", err, stmt)
	}
	if result != nil {
	}

	tx.Commit()

}
