package memselect

import (
	"os"
	"testing"

	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	"github.com/Okenamay/shorturl.git/internal/storage/memstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var Conf *config.Cfg
var TestLogger *zap.SugaredLogger

func TestMain(m *testing.M) {
	var err error

	TestLogger, err = logger.InitLogger()
	if err != nil {
		TestLogger.Fatalw("Tests stopped - start logger FAIL", "error", err)
	}
	defer TestLogger.Sync()

	Conf = config.InitConfig()
	Conf.MemMode = "memstore"

	os.Exit(m.Run())
}

func TestStoreAndCheckPair(t *testing.T) {
	memstorage.Store = memstorage.NewURLMap()

	userID := "test-user-123"
	shortID := "abcdef123"
	fullURL := "https://example.com/full"

	exists, err := StorePair(Conf, TestLogger, userID, shortID, fullURL)
	require.NoError(t, err)
	assert.False(t, exists)

	retrievedInfo, err := CheckPair(Conf, shortID)
	require.NoError(t, err)
	assert.Equal(t, fullURL, retrievedInfo.OriginalURL)
	assert.False(t, retrievedInfo.IsDeleted)
}

func TestProcessBatchTransaction(t *testing.T) {
	memstorage.Store = memstorage.NewURLMap()

	userID := "test-user-456"
	requestBatch := []RequestEntry{
		{CorrelationID: "a", OriginalURL: "https://test1.com"},
		{CorrelationID: "b", OriginalURL: "https://test2.com"},
	}

	responseBatch, err := ProcessBatchTransaction(Conf, TestLogger, requestBatch, userID)
	require.NoError(t, err)
	require.Len(t, responseBatch, 2)

	assert.Equal(t, "a", responseBatch[0].CorrelationID)
	assert.Contains(t, responseBatch[0].ShortURL, Conf.ShortIDAddress)

	allStored := memstorage.Store.GetAll()
	assert.Len(t, allStored, 2)
}

func TestGetUserURLs(t *testing.T) {
	t.Run("GetUserURLs_memstore_returns_nil", func(t *testing.T) {
		urls, err := GetUserURLs(Conf, TestLogger, "any-user-id")
		require.NoError(t, err)
		assert.Nil(t, urls)
	})
}
