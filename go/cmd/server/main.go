package main

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/powersjcb/monitor/go/src/lib/tracer"
	"github.com/powersjcb/monitor/go/src/server"
	"github.com/powersjcb/monitor/go/src/server/db"
	"github.com/powersjcb/monitor/go/src/server/gateway"
	"go.opentelemetry.io/otel/api/global"
	"log"
)

func main() {
	ctx := context.Background()
	c, err := server.GetConfig(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	conn, err := tracer.OpenDB(c.Database)
	if err != nil {
		log.Fatal(err.Error())
	}

	q := db.New(conn)
	tracer.InitTracer(c.HCAPIKey)
	t := global.Tracer("monitor.jacobpowers.me")
	ac := &gateway.ApplicationContext{
		Querier: q,
		Tracer:  t,
	}
	jwtConfig := gateway.JWTConfig{
		PublicKey:  c.JTWPublicKey,
		PrivateKey: c.JTWPrivateKey,
	}
	s := gateway.NewHTTPServer(ac, jwtConfig, c.Port)
	err = s.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
}
