package memstorage

var (
	URLStore = make(map[string]string) // Мапа для хранения пар ID – URL
)

// Сохранение пары fullURL-shortID в URLStore:
func StoreURLIDPair(shortID, fullURL string) {
	URLStore[shortID] = fullURL
}
