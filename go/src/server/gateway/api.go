package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/powersjcb/monitor/go/src/server/db"
	"net/http"
)

type Account struct {
	ID int64  `json:"id"`
	APIKey    string `json:"api_key"`
}

func (s HTTPServer) ShowAPIKey(rw http.ResponseWriter, r *http.Request) {
	accountID, err := AccountIDFromContext(r.Context())
	if err != nil {
		fmt.Println("accountID not available on context")
		rw.WriteHeader(500)
		return
	}
	account, err := s.appContext.Querier.GetAccountByID(r.Context(), accountID)
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

func (s HTTPServer) MetricStats(rw http.ResponseWriter, r *http.Request) {
	accountID, err := AccountIDFromContext(r.Context())
	if err != nil {
		fmt.Println("accountID not available on context")
		rw.WriteHeader(500)
		return
	}
	var p db.GetMetricStatsPerPeriodParams
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		fmt.Println(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
	}
	p.AccountID = accountID

	stats, err := s.appContext.Querier.GetMetricStatsPerPeriod(r.Context(), p)
	if err != nil {
		fmt.Println(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(rw).Encode(&stats)
	if err != nil {
		fmt.Println(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
