package usecases

import (
	"context"
	"database/sql"
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

	account, err := q.GetAccountByProviderID(ctx, db.GetAccountByProviderIDParams{AuthProvider: provider, AuthProviderID: providerID})
	if err != nil {
		if err != sql.ErrNoRows {
			return account, err
		}
	} else {
		return account, nil
	}

	return q.InsertAccount(ctx, db.InsertAccountParams{
		AuthProviderID: providerID,
		AuthProvider:   provider,
		ApiKey:         crypto.GetToken(64),
	})
}
