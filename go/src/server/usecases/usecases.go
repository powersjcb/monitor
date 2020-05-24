package usecases

import (
	"context"
	"errors"
	"github.com/powersjcb/monitor/go/src/lib/crypto"
	"github.com/powersjcb/monitor/go/src/server/db"
)

func GetOrCreateAccount(ctx context.Context, q db.Querier, provider, providerID string) (db.Account, error) {
	if provider == "" {
		return db.Account{}, errors.New("provider cannot be empty")
	}
	if providerID == "" {
		return db.Account{}, errors.New("poviderID cannot be empty")
	}
	return q.GetOrCreateAccount(ctx, db.GetOrCreateAccountParams{
		AuthProviderID: providerID,
		AuthProvider:   provider,
		ApiKey:         crypto.GetToken(64),
	})
}

