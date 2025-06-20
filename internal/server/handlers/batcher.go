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
		sugar, _ := logger.InitLogger()
		sugar.Info("BatchHandler. Start")

		var requestBatch []RequestEntry
		err := json.NewDecoder(r.Body).Decode(&requestBatch)
		if err != nil {
			sugar.Error("BatchHandler. Invalid request body format")
			http.Error(w, "Invalid request body format", http.StatusBadRequest)
			return
		}

		sugar.Infof("Processing batch of %d entries...", len(requestBatch))
		responseBatch, err := memselect.ProcessBatch(conf, requestBatch)
		if err != nil {
			sugar.Infof("Error processing batch: %v", err)
			http.Error(w, "Failed to process batch", http.StatusInternalServerError)
			return
		}
		sugar.Info("Batch processed successfully.")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(responseBatch)
	}
}
