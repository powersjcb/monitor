package main

import (
	"database/sql"
	"github.com/powersjcb/monitor/src/gateway"
	"github.com/powersjcb/monitor/src/server/db"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	conn, err := sql.Open("postgres", "host=127.0.0.1 dbname=monitor sslmode=disable")
	if err != nil {
		log.Fatal(err.Error())
	}
	q := db.New(conn)

	port := os.Getenv("PORT")
	s := gateway.NewHTTPServer(q, port)
	err = s.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
}
