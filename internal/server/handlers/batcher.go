package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Okenamay/shorturl.git/internal/app/middleware/auth"
	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/Okenamay/shorturl.git/internal/storage/memselect"
)

type RequestEntry = memselect.RequestEntry
type ResponseEntry = memselect.ResponseEntry

func BatchHandlerTransaction(conf *config.Cfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Zap.Info("BatchHandlerTransaction. Start")
		userID, _ := r.Context().Value(auth.UserIDContextKey).(string)

		var requestBatch []RequestEntry
		if err := json.NewDecoder(r.Body).Decode(&requestBatch); err != nil {
			logger.Zap.Errorw("BatchHandlerTransaction. Invalid request body format", "error", err)
			http.Error(w, "Invalid request body format", http.StatusBadRequest)
			return
		}

		logger.Zap.Infof("Processing transactional batch of %d entries...", len(requestBatch))
		responseBatch, err := memselect.ProcessBatchTransaction(conf, requestBatch, userID)
		if err != nil {
			logger.Zap.Errorw("Error processing transactional batch", "error", err)
			http.Error(w, "Failed to process batch", http.StatusInternalServerError)
			return
		}
		logger.Zap.Info("Transactional batch processed successfully.")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(responseBatch); err != nil {
			logger.Zap.Errorw("BatchHandlerTransaction. Failed to write response", "error", err)
		}
	}
}
