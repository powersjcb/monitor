package gateway

import (
	"github.com/powersjcb/monitor/src/client"
	"google.golang.org/appengine/log"
	"net/http"
)

func RunPings(w http.ResponseWriter, r *http.Request) {
	err := client.RunPings(client.DefaultPingConfigs, true)
	if err != nil {
		log.Errorf(r.Context(), err.Error())
	}
}