package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
)

// RequestEntry defines the structure of each object in the incoming JSON array.
type RequestEntry struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ResponseEntry struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type EntryToProcess struct {
	CorrelationID string
	OriginalURL   string
}

// ProcessedEntry represents a single processed URL.
// This is the data structure our ProcessBatch function will return.
type ProcessedEntry struct {
	CorrelationID string
	ShortURL      string
}

func BatchHandler(conf *config.Cfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sugar, _ := logger.InitLogger()
		sugar.Info("JSONHandler. Start")

		// 2. Decode the request body into a slice of RequestEntry
		var requestBatch []RequestEntry
		err := json.NewDecoder(r.Body).Decode(&requestBatch)
		if err != nil {
			http.Error(w, "Invalid request body format", http.StatusBadRequest)
			return
		}

		sugar.Info(requestBatch)
		for v := range requestBatch {
			temp := requestBatch[v]
			sugar.Infof("В батче запись номер: %d содержит: %s", v, temp)
		}

		// 3. Convert the request data into the format expected by ProcessBatch.
		// This decouples our API contract from our internal processing logic.
		rBatchStore := make([]EntryToProcess, len(requestBatch))
		for i, entry := range requestBatch {
			rBatchStore[i] = EntryToProcess{
				CorrelationID: entry.CorrelationID,
				OriginalURL:   entry.OriginalURL,
			}
		}

		sugar.Infof("%+v\n", rBatchStore)

		// // 4. Call the provided processing function
		// log.Printf("Processing batch of %d entries...", len(rBatchStore))
		// wBatchStore, err := memselect.ProcessBatch(conf, rBatchStore)
		// if err != nil {
		// 	// If the processing fails, return an internal server error.
		// 	log.Printf("Error processing batch: %v", err)
		// 	http.Error(w, "Failed to process batch", http.StatusInternalServerError)
		// 	return
		// // }
		// log.Println("Batch processed successfully.")

		// // 5. Convert the processed data into the response format.
		// responseBatch := make([]ResponseEntry, len(wBatchStore))
		// for i, entry := range wBatchStore {
		// 	responseBatch[i] = ResponseEntry{
		// 		CorrelationID: entry.CorrelationID,
		// 		ShortURL:      entry.ShortURL,
		// 	}
		// }

		// 6. Send the response
		w.Header().Set("Content-Type", "application/json")
		// Use StatusCreated (201) as we are "creating" new short URLs.
		w.WriteHeader(http.StatusCreated)
		// json.NewEncoder(w).Encode(responseBatch)
	}
}
