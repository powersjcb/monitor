package server

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type Config struct {
	Port     string
	Database string
	HCAPIKey string
}

func GetConfig(ctx context.Context) (Config, error) {
	c := Config{}
	sc, err := initClient(ctx)
	if err != nil {
		return c, err
	}

	port, err := getPort(nil)
	c.Port = port
	if err != nil {
		return c, err
	}

	db, err := getDBConnectionString(ctx, sc)
	c.Database = db
	if err != nil {
		return c, err
	}

	apiKey, err := getHoneycombKey(ctx, sc)
	c.HCAPIKey = apiKey
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

func getHoneycombKey(ctx context.Context, sc *secretmanager.Client) (string, error) {
	val := os.Getenv("HC_API_KEY")
	if val != "" {
		return val, nil
	}
	return getSecretValue(ctx, sc, "monitor_hc_api_key")
}

func getDBConnectionString(ctx context.Context, sc *secretmanager.Client) (string, error) {
	val := os.Getenv("DATABASE")
	if val != "" {
		return val, nil
	}
	return getSecretValue(ctx, sc, "monitor_db_connection")
}

// helper funcs

func getSecretValue(ctx context.Context, sc *secretmanager.Client, key string) (string, error) {
	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if project == "" {
		return "", errors.New("GOOGLE_CLOUD_PROJECT unset")
	}

	// projects/*/secrets/*/versions/latest is an alias to the latest
	res, err := sc.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", project, key),
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

func initClient(ctx context.Context) (*secretmanager.Client, error) {
	creds, err := google.FindDefaultCredentials(ctx, secretmanager.DefaultAuthScopes()...)
	if err != nil {
		return nil, err
	}
	sc, err := secretmanager.NewClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	return sc, err
}
