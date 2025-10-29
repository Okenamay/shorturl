package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Okenamay/shorturl.git/internal/app/middleware/auth"
	"github.com/Okenamay/shorturl.git/internal/config"
	"github.com/Okenamay/shorturl.git/internal/storage/memselect"
	"go.uber.org/zap"
)

type RequestEntry = memselect.RequestEntry
type ResponseEntry = memselect.ResponseEntry

func BatchHandlerTransaction(conf *config.Cfg, appLogger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appLogger.Info("BatchHandlerTransaction started")
		userID, _ := r.Context().Value(auth.UserIDContextKey).(string)

		var requestBatch []RequestEntry
		if err := json.NewDecoder(r.Body).Decode(&requestBatch); err != nil {
			appLogger.Errorw("BatchHandlerTransaction stopped - invalid request body format", "error", err)
			http.Error(w, "Invalid request body format", http.StatusBadRequest)
			return
		}

		appLogger.Infof("Processing transactional batch of %d entries...", len(requestBatch))
		responseBatch, err := memselect.ProcessBatchTransaction(conf, appLogger, requestBatch, userID)
		if err != nil {
			appLogger.Errorw("BatchHandlerTransaction stopped - processing transactional batch FAIL", "error", err)
			http.Error(w, "Failed to process batch", http.StatusInternalServerError)
			return
		}
		appLogger.Info("BatchHandlerTransaction finishing - transactional batch processed OK")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(responseBatch); err != nil {
			appLogger.Errorw("BatchHandlerTransaction stopped - response write FAIL", "error", err)
		}
	}
}
