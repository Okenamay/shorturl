package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/Okenamay/shorturl.git/internal/app/middleware/auth"
	"github.com/Okenamay/shorturl.git/internal/config"
	"github.com/Okenamay/shorturl.git/internal/server/handlers"
	"github.com/Okenamay/shorturl.git/internal/storage/memselect"
	"github.com/Okenamay/shorturl.git/internal/storage/memstorage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// setupExampleServer инициализирует роутер и тестовый сервер.
func setupExampleServer() (*httptest.Server, *config.Cfg) {
	testConf := config.InitConfig()
	testConf.MemMode = "memstore" // Используем in-memory хранилище
	testLogger := zap.NewNop().Sugar()

	// Инициализируем хранилище
	memstorage.Store = memstorage.NewURLMap()
	if err := memselect.MemInit(testConf, testLogger); err != nil {
		fmt.Printf("Failed to init memselect: %v\n", err)
		return nil, nil
	}

	r := chi.NewRouter()
	// Подключаем middleware и хендлеры
	r.Use(auth.Authenticator(testConf))
	r.Post("/", handlers.ShortenHandler(testConf, testLogger, nil))
	r.Get("/{id}", handlers.RedirectHandler(testConf, testLogger, nil))
	r.Post("/api/shorten", handlers.JSONHandler(testConf, testLogger, nil))
	r.Get("/api/user/urls", handlers.UserURLsHandler(testConf, testLogger))
	r.Post("/api/shorten/batch", handlers.BatchHandlerTransaction(testConf, testLogger))
	r.Delete("/api/user/urls", handlers.BatchDeleter(testConf))
	r.Get("/ping", handlers.PingHandler(testConf, testLogger))

	server := httptest.NewServer(r)

	// Обновляем конфиг, чтобы он указывал на тестовый сервер
	testConf.ShortIDAddress = server.URL

	return server, testConf
}

// ExampleShortenHandler демонстрирует использование эндпоинта POST /
func ExampleShortenHandler() {
	server, _ := setupExampleServer()
	if server == nil {
		fmt.Println("Failed to setup server")
		return
	}
	defer server.Close()

	client := server.Client()

	// 1. Создаем запрос на сокращение URL
	resp, err := client.Post(server.URL+"/", "text/plain", strings.NewReader("https://yandex.ru"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	fmt.Println("Status:", resp.StatusCode)
	fmt.Println("Body contains server URL:", strings.Contains(string(body), server.URL))
	// Output:
	// Status: 201
	// Body contains server URL: true
}

// ExampleJSONHandler демонстрирует использование эндпоинта POST /api/shorten
func ExampleJSONHandler() {
	server, _ := setupExampleServer()
	if server == nil {
		fmt.Println("Failed to setup server")
		return
	}
	defer server.Close()

	client := server.Client()

	// 1. Готовим JSON-тело
	reqBody, _ := json.Marshal(map[string]string{
		"url": "https://google.com",
	})

	// 2. Отправляем запрос
	resp, err := client.Post(server.URL+"/api/shorten", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// 3. Читаем и декодируем ответ
	var jsonResponse struct {
		Result string `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&jsonResponse); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Status:", resp.StatusCode)
	fmt.Println("Response contains server URL:", strings.Contains(jsonResponse.Result, server.URL))
	// Output:
	// Status: 201
	// Response contains server URL: true
}

// ExampleRedirectHandler демонстрирует редирект с короткого URL
func ExampleRedirectHandler() {
	server, _ := setupExampleServer()
	if server == nil {
		fmt.Println("Failed to setup server")
		return
	}
	defer server.Close()

	client := server.Client()

	// 1. Сначала создадим короткий URL
	originalURL := "https://practicum.yandex.ru"
	resp, err := client.Post(server.URL+"/", "text/plain", strings.NewReader(originalURL))
	if err != nil {
		fmt.Println(err)
		return
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	shortURL := string(body) // e.g., http://127.0.0.1:54321/f47c4cAB

	// 2. Делаем запрос на редирект
	// Мы используем кастомный CheckRedirect, чтобы остановить редирект
	// и проинспектировать 307 ответ.
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err = client.Get(shortURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.StatusCode)
	fmt.Println("Location:", resp.Header.Get("Location"))
	// Output:
	// Status: 307
	// Location: https://practicum.yandex.ru
}

// ExamplePingHandler демонстрирует проверку /ping
func ExamplePingHandler() {
	server, conf := setupExampleServer()
	if server == nil {
		fmt.Println("Failed to setup server")
		return
	}
	defer server.Close()

	// Для /ping в режиме "memstore" ожидается ошибка,
	// т.к. он работает только с "postgres".
	conf.MemMode = "memstore"

	resp, err := server.Client().Get(server.URL + "/ping")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.StatusCode)
	// Output:
	// Status: 500
}

// ExampleUserURLsHandler демонстрирует получение URL пользователя
// (Это пример без проверки вывода, т.к. cookie непредсказуемы)
func ExampleUserURLsHandler() {
	server, _ := setupExampleServer()
	if server == nil {
		fmt.Println("Failed to setup server")
		return
	}
	defer server.Close()
	client := server.Client()

	// 1. Создаем URL, чтобы получить аутентификационный cookie
	resp, err := client.Post(server.URL+"/", "text/plain", strings.NewReader("https://google.com"))
	if err != nil {
		fmt.Println(err)
		return
	}
	resp.Body.Close()

	// Ищем cookie "token"
	var userCookie *http.Cookie
	for _, c := range resp.Cookies() {
		if c.Name == "token" {
			userCookie = c
			break
		}
	}
	if userCookie == nil {
		fmt.Println("Token cookie not found")
		return
	}

	// 2. Делаем запрос к /api/user/urls с этим cookie
	req, _ := http.NewRequest("GET", server.URL+"/api/user/urls", nil)
	req.AddCookie(userCookie)

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Println("Status:", resp.StatusCode)
	log.Println("Body:", string(body))
	// Ожидаемый вывод будет JSON-массив, например:
	// [{"short_url":"...","original_url":"https://google.com"}]
}

// ExampleBatchHandlerTransaction демонстрирует пакетное сокращение URL
func ExampleBatchHandlerTransaction() {
	server, _ := setupExampleServer()
	if server == nil {
		fmt.Println("Failed to setup server")
		return
	}
	defer server.Close()
	client := server.Client()

	batch := []map[string]string{
		{"correlation_id": "1", "original_url": "https://apple.com"},
		{"correlation_id": "2", "original_url": "https://microsoft.com"},
	}
	reqBody, _ := json.Marshal(batch)

	resp, err := client.Post(server.URL+"/api/shorten/batch", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	var respBody []map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&respBody)

	fmt.Println("Status:", resp.StatusCode)
	fmt.Println("Response count:", len(respBody))
	fmt.Println("Correlation ID:", respBody[0]["correlation_id"])
	// Output:
	// Status: 201
	// Response count: 2
	// Correlation ID: 1
}

// ExampleBatchDeleter демонстрирует асинхронное удаление
// (Это пример без проверки вывода)
func ExampleBatchDeleter() {
	server, _ := setupExampleServer()
	if server == nil {
		fmt.Println("Failed to setup server")
		return
	}
	defer server.Close()
	client := server.Client()

	// 1. Получаем cookie
	resp, err := client.Post(server.URL+"/", "text/plain", strings.NewReader("https://to-be-deleted.com"))
	if err != nil {
		fmt.Println(err)
		return
	}
	resp.Body.Close()

	var userCookie *http.Cookie
	for _, c := range resp.Cookies() {
		if c.Name == "token" {
			userCookie = c
			break
		}
	}
	if userCookie == nil {
		fmt.Println("Token cookie not found")
		return
	}

	// 2. Готовим ID для удаления (для примера)
	idsToDelete := []string{"id1", "id2", "id3"}
	reqBody, _ := json.Marshal(idsToDelete)

	// 3. Отправляем запрос на удаление
	req, _ := http.NewRequest("DELETE", server.URL+"/api/user/urls", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(userCookie)

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	log.Println("Status:", resp.StatusCode)
	// Ожидаемый статус: 202 Accepted
}
