package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Okenamay/shorturl.git/internal/audit"
	"github.com/Okenamay/shorturl.git/internal/config"
	logger "github.com/Okenamay/shorturl.git/internal/logger/zap"
	pb "github.com/Okenamay/shorturl.git/internal/proto"
	grpcserver "github.com/Okenamay/shorturl.git/internal/server/grpc"
	"github.com/Okenamay/shorturl.git/internal/server/router"
	"github.com/Okenamay/shorturl.git/internal/storage/memselect"
	"github.com/Okenamay/shorturl.git/internal/worker"
	"google.golang.org/grpc"
)

// Переменные для информации о сборке (будут установлены через -ldflags)
var buildVersion string
var buildDate string
var buildCommit string

func main() {
	// Вспомогательная функция для вывода "N/A", если переменная пуста
	na := func(s string) string {
		if s == "" {
			return "N/A"
		}
		return s
	}

	// Выводим информацию о сборке в stdout
	fmt.Printf("Build version: %s\n", na(buildVersion))
	fmt.Printf("Build date: %s\n", na(buildDate))
	fmt.Printf("Build commit: %s\n", na(buildCommit))

	appLogger, err := logger.InitLogger()
	if err != nil {
		// Если логгер не стартовал, мы не можем даже это залогировать,
		// паникуем. Используем log.Fatalf, что разрешено в main.main
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer appLogger.Sync()

	conf, err := config.InitConfig()
	if err != nil {
		appLogger.Fatalw("failed to init config", "error", err)
	}
	auditor := audit.NewAuditor(conf.AuditFile, conf.AuditURL, appLogger)

	worker.Start(memselect.BatchDelete, appLogger)

	err = memselect.MemInit(conf, appLogger)
	if err != nil {
		appLogger.Errorw(err.Error(), "Main", "Initialize storage")
	}
	// MemStop отложим до graceful shutdown, чтобы не закрыть базу раньше времени

	// Создаем сервер (но не запускаем, это сделаем в горутине)
	srv := router.CreateServer(conf, appLogger, auditor)

	// Создаем gRPC сервер
	grpcSrvImpl := grpcserver.NewShortenerServer(conf, appLogger, auditor)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcSrvImpl.AuthInterceptor),
	)
	pb.RegisterShortenerServiceServer(grpcServer, grpcSrvImpl)

	// Канал для перехвата сигналов ОС
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// Запускаем HTTP сервер в отдельной горутине
	go func() {
		appLogger.Infof("Starting HTTP server on port: %s (HTTPS: %v)", conf.ServerPort, conf.EnableHTTPS)
		var err error
		if conf.EnableHTTPS {
			err = srv.ListenAndServeTLS("", "")
		} else {
			err = srv.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			appLogger.Fatalw(err.Error(), "Main", "Start HTTP server")
		}
	}()

	// Запускаем gRPC сервер в отдельной горутине
	go func() {
		if conf.GRPCAddress != "" {
			listen, err := net.Listen("tcp", conf.GRPCAddress)
			if err != nil {
				appLogger.Fatalw("Failed to listen tcp for gRPC", "error", err)
			}
			appLogger.Infof("Starting gRPC server on port: %s", conf.GRPCAddress)
			if err := grpcServer.Serve(listen); err != nil {
				appLogger.Fatalw("gRPC Serve error", "error", err)
			}
		}
	}()

	// Блокируем main и ждем сигнала
	sig := <-quit
	appLogger.Infof("Received signal: %v. Initiating graceful shutdown...", sig)

	// Создаем контекст с таймаутом для завершения текущих запросов
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Останавливаем прием новых запросов HTTP и ждем завершения текущих
	if err := srv.Shutdown(shutdownCtx); err != nil {
		appLogger.Errorf("HTTP Server forced to shutdown: %v", err)
	} else {
		appLogger.Info("HTTP Server stopped gracefully")
	}

	// 2. Останавливаем gRPC сервер
	grpcServer.GracefulStop()
	appLogger.Info("gRPC Server stopped gracefully")

	// 3. Останавливаем воркеры (они сбросят буферы в БД), делаем это после
	// остановки сервера, чтобы новые запросы на удаление не поступали
	worker.Stop(appLogger)

	// 4. Закрываем соединения с БД
	memselect.MemStop(conf, appLogger)

	appLogger.Info("Application exited properly")

	defer appLogger.Sync()
}
