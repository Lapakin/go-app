package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

const (
	dbURL  = "postgres://admin:qwerty@mydatabase:5432/postgres?sslmode=disable"
	migDir = "file:///usr/src/app/migrations"
)

func main() {
	//two connections to the BD, because different types - .pool and .db
	dbpool, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatal("pgxpool: ", err)
	}
	defer dbpool.Close()

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("sqlopen: ", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("driver: ", err)
	}

	m, err := migrate.NewWithDatabaseInstance(migDir, "postgres", driver)
	if err != nil {
		log.Fatal("migrate: ", err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("migrate up: ", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			name := r.FormValue("field1")
			salary := r.FormValue("field2")
			time := time.Now()

			query := "INSERT INTO list_db (full_name, salary, date_receive) VALUES ($1, $2, $3)"

			if _, err = dbpool.Exec(context.Background(), query, name, salary, time); err != nil {
				fmt.Fprintf(os.Stderr, "Exec failed: %v\n", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			fmt.Fprintln(w, "Data has been successfully added to the database!")
			return
		}

		http.ServeFile(w, r, "index.html")
	})

	if err = http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal("host: ", err)
	}
}
