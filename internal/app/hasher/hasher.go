package hasher

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
	"io"
	"sync"

	"github.com/Okenamay/shorturl.git/internal/config"
)

// Изменение: Создаем пул для md5.Hash
// Это экономит аллокацию md5.New() на каждый вызов
var md5Pool = sync.Pool{
	New: func() interface{} {
		return md5.New()
	},
}

// Кодирование строки с URL в md5-сумму с обрезанием до ShortIDLen символов:
func ShortenURL(conf *config.Cfg, fullURL string) string {

	// Изменение: Берем хешер из пула
	hash := md5Pool.Get().(hash.Hash)
	hash.Reset() // Сбрасываем состояние хешера

	io.WriteString(hash, fullURL)

	// Изменение: Оптимизируем аллокации
	// Выделяем 16 байт (размер md5) на стеке
	var md5sum [16]byte
	// hash.Sum пишет результат в срез (md5sum[:0]), не выделяя новую память
	sum := hash.Sum(md5sum[:0])

	md5Pool.Put(hash)

	// Выделяем 32 байта (md5 в hex) на стеке
	var hexBuf [32]byte
	// Кодируем md5sum в hex-буфер
	hex.Encode(hexBuf[:], sum)

	// Возвращаем строку, сделанную из среза нужной длины
	return string(hexBuf[:conf.ShortIDLen])
}
