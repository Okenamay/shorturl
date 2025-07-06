package memselect

import (
	"os"
	"testing"

	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/Okenamay/shorturl.git/internal/storage/memstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testConf *config.Cfg

func TestMain(m *testing.M) {
	if err := logger.InitLogger(); err != nil {
		panic("failed to initialize logger for tests: " + err.Error())
	}
	testConf = config.InitConfig()
	testConf.MemMode = "memstore"

	os.Exit(m.Run())
}

func TestStoreAndCheckPair(t *testing.T) {
	memstorage.Store = memstorage.NewURLMap()

	shortID := "abcdef123"
	fullURL := "https://example.com/full"

	exists, err := StorePair(testConf, shortID, fullURL)
	require.NoError(t, err)
	assert.False(t, exists)

	retrievedURL, err := CheckPair(testConf, shortID)
	require.NoError(t, err)
	assert.Equal(t, fullURL, retrievedURL)
}

func TestProcessBatchTransaction(t *testing.T) {
	memstorage.Store = memstorage.NewURLMap()

	requestBatch := []RequestEntry{
		{CorrelationID: "a", OriginalURL: "https://test1.com"},
		{CorrelationID: "b", OriginalURL: "https://test2.com"},
	}

	responseBatch, err := ProcessBatchTransaction(testConf, requestBatch)
	require.NoError(t, err)
	require.Len(t, responseBatch, 2)

	assert.Equal(t, "a", responseBatch[0].CorrelationID)
	assert.Contains(t, responseBatch[0].ShortURL, testConf.ShortIDServerPort)

	allStored := memstorage.Store.GetAll()
	assert.Len(t, allStored, 2)
}
