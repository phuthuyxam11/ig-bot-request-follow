package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func Connection(dataSource string) (error, *sql.DB) {
	database, err := sql.Open("sqlite3", dataSource)
	if err != nil {
		return err, nil
	}

	return nil, database
}

func CreateTable(db *sql.DB) error {
	createTableSql := `
		CREATE TABLE IF NOT EXISTS "users"
		(
		  "pk_id" TEXT NOT NULL,
		  "user_name" TEXT,
		  "Is_private" integer,
		  "follow_flag" TEXT,
		  "created_at" TEXT,
		  "modified_at" TEXT,
		  PRIMARY KEY ("pk_id")
		);
		`
	statement, err := db.Prepare(createTableSql)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}

	// clear data
	deleteStatement, err := db.Prepare("DELETE FROM users")
	if err != nil {
		return err
	}
	_, err = deleteStatement.Exec()
	if err != nil {
		return err
	}

	return nil
}
