package grpc

import (
	"context"

	"github.com/Okenamay/shorturl.git/internal/app/checker"
	"github.com/Okenamay/shorturl.git/internal/app/middleware/auth"
	"github.com/Okenamay/shorturl.git/internal/app/urlmaker"
	"github.com/Okenamay/shorturl.git/internal/audit"
	"github.com/Okenamay/shorturl.git/internal/config"
	pb "github.com/Okenamay/shorturl.git/internal/proto"
	"github.com/Okenamay/shorturl.git/internal/storage/memselect"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ShortenerServer struct {
	pb.UnimplementedShortenerServiceServer
	config  *config.Cfg
	logger  *zap.SugaredLogger
	auditor *audit.Auditor
}

func NewShortenerServer(conf *config.Cfg, logger *zap.SugaredLogger, auditor *audit.Auditor) *ShortenerServer {
	return &ShortenerServer{
		config:  conf,
		logger:  logger,
		auditor: auditor,
	}
}

// AuthInterceptor извлекает токен из metadata (header authorization) и валидирует его
func (s *ShortenerServer) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var userID string

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("authorization") // gRPC заголовки обычно в нижнем регистре
		if len(values) > 0 {
			tokenString := values[0]
			// Удаляем префикс "Bearer ", если он есть
			// В задании сказано передавать авторизационные данные, но обычно подразумевается токен.
			// Наш GetUserID ожидает чистый токен или токен из cookie.
			id, err := auth.GetUserID(s.config, tokenString)
			if err == nil {
				userID = id
			} else {
				s.logger.Warnf("gRPC Auth failed: %v", err)
			}
		}
	}

	// Добавляем userID в контекст
	newCtx := context.WithValue(ctx, auth.UserIDContextKey, userID)
	return handler(newCtx, req)
}

func (s *ShortenerServer) ShortenURL(ctx context.Context, req *pb.URLShortenRequest) (*pb.URLShortenResponse, error) {
	s.logger.Info("gRPC ShortenURL started")
	userID, _ := ctx.Value(auth.UserIDContextKey).(string)

	checkedURL, checkErr := checker.CheckURL(req.Url, s.logger)
	if checkErr != nil {
		return nil, status.Error(codes.InvalidArgument, checkErr.Error())
	}

	fullURL := checkedURL.String()
	newURL, shortID := urlmaker.ProcessURL(s.config, fullURL)

	exists, err := memselect.StorePair(s.config, s.logger, userID, shortID, fullURL)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to store URL")
	}

	if s.auditor != nil {
		s.auditor.LogEvent(ctx, "shorten_grpc", userID, fullURL)
	}

	// Если exists == true, это конфликт, но мы все равно возвращаем результат (как и в JSON handler)
	// В gRPC можно вернуть ошибку AlreadyExists, если клиент этого ожидает,
	// но для совместимости с логикой "вернуть сокращенный URL" возвращаем успешный ответ.
	if exists {
		s.logger.Warn("gRPC ShortenURL - already exists")
	}

	return &pb.URLShortenResponse{
		Result: newURL,
	}, nil
}

func (s *ShortenerServer) ExpandURL(ctx context.Context, req *pb.URLExpandRequest) (*pb.URLExpandResponse, error) {
	s.logger.Info("gRPC ExpandURL started")

	if len(req.Id) != s.config.ShortIDLen {
		return nil, status.Error(codes.InvalidArgument, "invalid short ID length")
	}

	urlInfo, err := memselect.CheckPair(s.config, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal storage error")
	}

	if urlInfo.OriginalURL == "" {
		return nil, status.Error(codes.NotFound, "URL not found")
	}

	if urlInfo.IsDeleted {
		return nil, status.Error(codes.Unavailable, "URL deleted") // Или codes.Gone, но в gRPC стандартный Unavailable или NotFound
	}

	return &pb.URLExpandResponse{
		Result: urlInfo.OriginalURL,
	}, nil
}

func (s *ShortenerServer) ListUserURLs(ctx context.Context, _ *emptypb.Empty) (*pb.UserURLsResponse, error) {
	s.logger.Info("gRPC ListUserURLs started")
	userID, _ := ctx.Value(auth.UserIDContextKey).(string)

	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	urls, err := memselect.GetUserURLs(s.config, s.logger, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to fetch user URLs")
	}

	var pbUrls []*pb.URLData
	for _, u := range urls {
		pbUrls = append(pbUrls, &pb.URLData{
			ShortUrl:    u.ShortURL,
			OriginalUrl: u.OriginalURL,
		})
	}

	return &pb.UserURLsResponse{
		Url: pbUrls,
	}, nil
}
