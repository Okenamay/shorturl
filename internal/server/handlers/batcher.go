package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/Okenamay/shorturl.git/internal/storage/memselect"
)

type RequestEntry = memselect.RequestEntry
type ResponseEntry = memselect.ResponseEntry

func BatchHandler(conf *config.Cfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Zap.Info("BatchHandler. Start")

		var requestBatch []RequestEntry
		if err := json.NewDecoder(r.Body).Decode(&requestBatch); err != nil {
			logger.Zap.Errorw("BatchHandler. Invalid request body format", "error", err)
			http.Error(w, "Invalid request body format", http.StatusBadRequest)
			return
		}

		logger.Zap.Infof("Processing batch of %d entries...", len(requestBatch))
		responseBatch, err := memselect.ProcessBatch(conf, requestBatch)
		if err != nil {
			logger.Zap.Errorw("Error processing batch", "error", err)
			http.Error(w, "Failed to process batch", http.StatusInternalServerError)
			return
		}
		logger.Zap.Info("Batch processed successfully.")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(responseBatch); err != nil {
			logger.Zap.Errorw("BatchHandler. Failed to write response", "error", err)
		}
	}
}

func BatchHandlerTransaction(conf *config.Cfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Zap.Info("BatchHandlerTransaction. Start")

		var requestBatch []RequestEntry
		if err := json.NewDecoder(r.Body).Decode(&requestBatch); err != nil {
			logger.Zap.Errorw("BatchHandlerTransaction. Invalid request body format", "error", err)
			http.Error(w, "Invalid request body format", http.StatusBadRequest)
			return
		}

		logger.Zap.Infof("Processing transactional batch of %d entries...", len(requestBatch))
		responseBatch, err := memselect.ProcessBatchTransaction(conf, requestBatch)
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
