package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("403 Forbidden"))
		return
	}

	if err := cfg.db.ResetUsers(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error reseting users"))
		return
	}

	
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and users reset"))
}
