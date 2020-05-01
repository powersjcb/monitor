package main

import (
	"database/sql"
	"github.com/powersjcb/monitor/src/gateway"
	"github.com/powersjcb/monitor/src/server/db"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	conn, err := sql.Open("postgres", "host=127.0.0.1 dbname=monitor sslmode=disable")
	if err != nil {
		log.Fatal(err.Error())
	}
	q := db.New(conn)

	s := gateway.NewHTTPServer(q)
	err = s.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
}
