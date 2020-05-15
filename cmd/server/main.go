package main

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/powersjcb/monitor/src/gateway"
	"github.com/powersjcb/monitor/src/lib/tracer"
	"github.com/powersjcb/monitor/src/server"
	"github.com/powersjcb/monitor/src/server/db"
	"go.opentelemetry.io/otel/api/global"
	"log"
)

func main() {
	ctx := context.Background()
	c, err := server.GetConfig(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	conn, err := sql.Open("postgres", c.Database)
	if err != nil {
		log.Fatal(err.Error())
	}

	q := db.New(conn)
	tracer.InitTracer(c.HCAPIKey)
	t := global.Tracer("monitor.jacobpowers.me")
	ac := &gateway.ApplicationContext{
		Tracer: t,
	}
	s := gateway.NewHTTPServer(ac, q, c.Port)
	err = s.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
}
