package beater

import (
	"fmt"
	"github.com/elastic/apm-server/agentcfg"
	"net/http"
)

func agentConfigHandler(lookupCfg agentcfg.LookupFunc, secretToken string) http.Handler {

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		params := r.URL.Query()
		serviceName := params.Get("service_name")
		serviceEnv := params.Get("service_env")
		cfg, etag := lookupCfg(serviceName, serviceEnv)

		if len(cfg) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", int64(agentcfg.PollInterval.Seconds())))
		if clientEtag := r.Header.Get("If-None-Match"); clientEtag != "" && clientEtag == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		if etag != "" {
			w.Header().Set("Etag", etag)
		}

		send(w, r, cfg, http.StatusOK)
	})

	return authHandler(secretToken, logHandler(handler))
}
