package gateway_test

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"database/sql"
	"encoding/json"
	"github.com/powersjcb/monitor/go/src/server/db"
	"github.com/powersjcb/monitor/go/src/server/gateway"
	"github.com/powersjcb/monitor/go/src/server/gateway/mocks"
	"github.com/stretchr/testify/mock"
	"net/http/httptest"
	"testing"
	"time"
)

type mockRandomReader struct {
}

func (mockRandomReader) Read(b []byte) (int, error) {
	return len(b), nil
}

func server() (gateway.HTTPServer, *mocks.Querier) {
	q := &mocks.Querier{}
	appContext := gateway.ApplicationContext{
		Querier: q,
	}

	key, err := ecdsa.GenerateKey(elliptic.P256(), mockRandomReader{})
	if err != nil {
		panic(err)
	}
	jwtConfig := gateway.JWTConfig{
		PublicKey:  key.PublicKey,
		PrivateKey: *key,
	}
	return gateway.NewHTTPServer(&appContext, jwtConfig, gateway.OAUTHConfig{}, "9999", ""), q
}

func TestHTTPServer_Metric_EmptyPost(t *testing.T) {
	s, _ := server()

	r := httptest.NewRequest("POST", "/metric", nil)
	w := httptest.NewRecorder()

	s.Metric(w, r)

	if w.Code == 200 {
		t.Errorf("returned a 200 response code without a request body: %d", w.Code)
	}
}

func TestHTTPServer_Metric_Valid(t *testing.T) {
	s, q := server()
	metricParams := db.InsertMetricParams{
		Ts: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Source: "ping",
		Name:   "test",
		Target: "test-target",
		Value: sql.NullFloat64{
			Float64: 1.0001,
			Valid:   true,
		},
	}

	metric := db.Metric{
		Source:     metricParams.Source,
		Ts:         metricParams.Ts,
		InsertedAt: time.Now(),
		Name:       metricParams.Name,
		Target:     metricParams.Target,
		Value:      metricParams.Value,
	}

	data, err := json.Marshal(metricParams)
	if err != nil {
		t.Errorf("unable to marshal test data: %s", err.Error())
	}

	r := httptest.NewRequest("POST", "/metric", bytes.NewReader(data))
	w := httptest.NewRecorder()

	q.On("InsertMetric", mock.Anything, mock.AnythingOfType("db.InsertMetricParams")).Return(metric, nil)
	s.Metric(w, r)

	if w.Code != 200 {
		t.Errorf("unable to insert metric: %d, %s", w.Code, w.Body.String())
	}
}
