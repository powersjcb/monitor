package server

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type Config struct {
	Port          string
	Database      string
	HCAPIKey      string
	JTWPublicKey  ecdsa.PublicKey
	JTWPrivateKey ecdsa.PrivateKey
	OAuthClientID     string
	OAuthClientSecret string
	OAuthRedirectURL  string
	APIKey			  string
}

func GetConfig(ctx context.Context) (Config, error) {
	c := Config{}
	secretClient, err := initClient(ctx)
	if err != nil {
		return c, err
	}

	port, err := getPort(ctx)
	c.Port = port
	if err != nil {
		return c, err
	}

	db, err := getDBConnectionString(ctx, secretClient)
	c.Database = db
	if err != nil {
		return c, err
	}

	apiKey, err := getHoneycombKey(ctx, secretClient)
	c.HCAPIKey = apiKey
	if err != nil {
		return c, err
	}

	publicKeyData, err := getPublicKey(ctx, secretClient)
	if err != nil {
		return c, err
	}
	publicKey, err := jwt.ParseECPublicKeyFromPEM([]byte(publicKeyData))
	if err != nil || publicKey == nil {
		return c, err
	}
	c.JTWPublicKey = *publicKey

	privateKeyData, err := getPrivateKey(ctx, secretClient)
	if err != nil {
		return c, err
	}
	privateKey, err := jwt.ParseECPrivateKeyFromPEM([]byte(privateKeyData))
	if err != nil || privateKey == nil {
		return c, err
	}
	c.JTWPrivateKey = *privateKey

	clientID, err := getGoogleClientID(ctx, secretClient)
	c.OAuthClientID = clientID
	if err != nil {
		return c, err
	}

	cs, err := getGoogleClientSecret(ctx, secretClient)
	c.OAuthClientSecret = cs
	if err != nil {
		return c, err
	}

	r, err := getGoogleRedirectURL(ctx, secretClient)
	c.OAuthRedirectURL = r
	if err != nil {
		return c, err
	}

	key, err := getAPIKey(ctx, secretClient)
	c.APIKey = key
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

func getHoneycombKey(ctx context.Context, secretClient *secretmanager.Client) (string, error) {
	val := os.Getenv("HC_API_KEY")
	if val != "" {
		return val, nil
	}
	return getSecretValue(ctx, secretClient, "monitor_hc_api_key")
}

func getDBConnectionString(ctx context.Context, secretClient *secretmanager.Client) (string, error) {
	val := os.Getenv("DATABASE")
	if val != "" {
		return val, nil
	}
	return getSecretValue(ctx, secretClient, "monitor_db_connection")
}

func getPublicKey(ctx context.Context, secretClient *secretmanager.Client) (string, error) {
	val := os.Getenv("JWT_EC_PUBLIC_KEY")
	if val != "" {
		return val, nil
	}
	return getSecretValue(ctx, secretClient, "monitor_jwt_ec_public_key")
}

func getPrivateKey(ctx context.Context, secretClient *secretmanager.Client) (string, error) {
	val := os.Getenv("JWT_EC_PRIVATE_KEY")
	if val != "" {
		return val, nil
	}
	return getSecretValue(ctx, secretClient, "monitor_jwt_ec_private_key")
}

func getGoogleClientID(ctx context.Context, secretClient *secretmanager.Client) (string, error) {
	val := os.Getenv("GOOGLE_CLIENT_ID")
	if val != "" {
		return val, nil
	}
	return getSecretValue(ctx, secretClient, "monitor_google_client_id")
}

func getGoogleClientSecret(ctx context.Context, secretClient *secretmanager.Client) (string, error) {
	val := os.Getenv("GOOGLE_CLIENT_SECRET")
	if val != "" {
		return val, nil
	}
	return getSecretValue(ctx, secretClient, "monitor_google_client_secret")
}

func getGoogleRedirectURL(ctx context.Context, secretClient *secretmanager.Client) (string, error) {
	val := os.Getenv("GOOGLE_CLIENT_REDIRECT_URL")
	if val != "" {
		return val, nil
	}
	return getSecretValue(ctx, secretClient, "monitor_google_client_redirect_url")
}

func getAPIKey(ctx context.Context, secretClient *secretmanager.Client) (string, error) {
	val := os.Getenv("MONITOR_API_KEY")
	if val != "" {
		return val, nil
	}
	return getSecretValue(ctx, secretClient, "monitor_api_key")
}

// helper funcs

func getSecretValue(ctx context.Context, secretClient *secretmanager.Client, key string) (string, error) {
	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if project == "" {
		return "", errors.New("GOOGLE_CLOUD_PROJECT unset")
	}

	// projects/*/secrets/*/versions/latest is an alias to the latest
	res, err := secretClient.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", project, key),
	})
	if err != nil {
		return "", err
	}
	data := res.Payload.GetData()
	if len(data) == 0 {
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
