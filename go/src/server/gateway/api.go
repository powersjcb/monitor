package gateway

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Account struct {
	ID int64  `json:"id"`
	APIKey    string `json:"api_key"`
}

func (s HTTPServer) ShowAPIKey(rw http.ResponseWriter, r *http.Request) {
	userID, err := UserIDFromContext(r.Context())
	if err != nil {
		fmt.Println("userID not available on context")
		rw.WriteHeader(500)
	}
	account, err := s.appContext.Querier.GetAccountByID(r.Context(), userID)
	if err != nil {
		fmt.Println(err.Error())
		rw.WriteHeader(500)
		return
	}

	encoder := json.NewEncoder(rw)
	err = encoder.Encode(&Account{
		ID: account.ID,
		APIKey:    account.ApiKey,
	})
	if err != nil {
		_, err = rw.Write([]byte(err.Error()))
		if err != nil {
			fmt.Printf("failed to write error response: %s", err.Error())
		}
		rw.WriteHeader(500)
		return
	}
}

