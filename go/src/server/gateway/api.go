package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/powersjcb/monitor/go/src/server/db"
	"io/ioutil"
	"net/http"
)

type Account struct {
	ID     int64  `json:"id"`
	APIKey string `json:"api_key"`
}

func (s HTTPServer) ShowAPIKey(rw http.ResponseWriter, r *http.Request) {
	accountID, err := AccountIDFromContext(r.Context())
	if err != nil {
		fmt.Println(err.Error())
		rw.WriteHeader(500)
		return
	}
	account, err := s.appContext.Querier.GetAccountByID(r.Context(), accountID)
	if err != nil {
		fmt.Println("GetAccountByID: ", err.Error())
		rw.WriteHeader(500)
		return
	}

	encoder := json.NewEncoder(rw)
	err = encoder.Encode(&Account{
		ID:     account.ID,
		APIKey: account.ApiKey,
	})
	if err != nil {
		fmt.Println("encoding and sending response: ", err.Error())
		rw.WriteHeader(500)
		return
	}
}

func (s HTTPServer) MetricStats(rw http.ResponseWriter, r *http.Request) {
	accountID, err := AccountIDFromContext(r.Context())
	if err != nil {
		fmt.Println(err.Error())
		rw.WriteHeader(500)
		return
	}
	var p db.GetMetricStatsPerPeriodParams
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("failed to read body")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &p)
	if err != nil {
		fmt.Printf("decoding failed: %s, request.Body: '%s'\n", err.Error(), data)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	p.AccountID = accountID

	stats, err := s.appContext.Querier.GetMetricStatsPerPeriod(r.Context(), p)
	if err != nil {
		fmt.Println("db query failed: " + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(&stats)
	if err != nil {
		fmt.Println("encoding failed: " + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(200)
}
