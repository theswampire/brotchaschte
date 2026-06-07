package database

import (
	"database/sql"
	"log"

	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed "sql"
var SQLFiles embed.FS

// migrate migrates the given db using goose and panics if it fails
func migrate(db *sql.DB) {
	goose.SetBaseFS(SQLFiles)
	if err := goose.SetDialect("sqlite"); err != nil {
		panic(err)
	}
	if err := goose.Up(db, "sql/migrations"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "sql/seed", goose.WithNoVersioning()); err != nil {
		log.Println("already seeded")
	}
}
