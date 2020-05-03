package server

import (
	"context"
	"errors"
	"fmt"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)
type Config struct {
	Port string
	Database string
}

func GetConfig(ctx context.Context) (Config, error) {
	c := Config{}

	port, err := getPort(nil)
	c.Port = port
	if err != nil {
		return c, err
	}

	db, err := getDBConnectionString(ctx)
	c.Database = db
	if err != nil {
		return c, err
	}

	return c, nil
}

func getPort(_ context.Context) (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return port, errors.New("invalid port: '" + port + "'")
	}
	return port, nil
}

func getDBConnectionString(ctx context.Context) (string, error) {
	val := os.Getenv("DATABASE")
	if val != "" {
		return val, nil
	}
	return getSecretValue(ctx, "monitor_database_connection_string")
}

// helper funcs

func getSecretValue(ctx context.Context, key string) (string, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", err
	}
	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if project == "" {
		return "", errors.New("GOOGLE_CLOUD_PROJECT unset")
	}

	// projects/*/secrets/*/versions/latest is an alias to the latest
	res, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest"),
	})
	if err != nil {
		return "", err
	}
	data := res.Payload.GetData()
	if data == nil || len(data) == 0 {
		return "", errors.New("no secret found for key: " + key)
	}
	return string(data), nil
}
