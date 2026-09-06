package api

import "net/http"

type healthResponse struct {
	Status string `json:"status"`
}

// Health reports that the process is up and serving. It deliberately performs
// no upstream checks, so it stays cheap enough to poll as a liveness probe and
// never reports this service as unhealthy because GitHub is having a bad day.
func Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	writeJSON(w, http.StatusOK, healthResponse{Status: "ok"})
}
