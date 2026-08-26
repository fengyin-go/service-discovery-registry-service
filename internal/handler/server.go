// Package handler 实现 HTTP 处理器层。
package handler

import (
	"errors"
	"net/http"
	"runtime/debug"
	"time"

	"serviceregistry/internal/config"
	"serviceregistry/internal/model"
	"serviceregistry/internal/service"
	"serviceregistry/internal/store"
	"serviceregistry/pkg/httpx"
	"serviceregistry/pkg/logger"
)

// Server 聚合服务与配置，提供路由。
type Server struct {
	svc       *service.Service
	log       *logger.Logger
	cfg       *config.Config
	authToken string
	limiter   *rateLimiter
}

// NewServer 构造 HTTP 服务器。
func NewServer(svc *service.Service, log *logger.Logger, cfg *config.Config) *Server {
	return &Server{
		svc:       svc,
		log:       log,
		cfg:       cfg,
		authToken: getenv("AUTH_TOKEN", ""),
		limiter:   newRateLimiter(getenvInt("RATE_LIMIT", 200), time.Minute),
	}
}

// Routes 组装全部路由并包裹中间件。
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	s.registerServiceRoutes(mux)
	s.registerInstanceRoutes(mux)
	s.registerHealthCheckRoutes(mux)
	s.registerHeartbeatRoutes(mux)
	s.registerLoadBalanceRoutes(mux)
	s.registerDependencyRoutes(mux)
	s.registerEventRoutes(mux)
	s.registerStatsRoutes(mux)
	s.registerExportRoutes(mux)
	s.registerHealthRoutes(mux)
	s.registerLabelRoutes(mux)
	s.registerConfigRoutes(mux)
	s.registerSearchRoutes(mux)
	s.registerReportRoutes(mux)
	s.registerTopologyRoutes(mux)
	s.registerMetricsRoutes(mux)
	s.registerOpsRoutes(mux)
	s.registerProbeRoutes(mux)

	var h http.Handler = mux
	h = s.authMiddleware(h)
	h = s.rateLimitMiddleware(h)
	h = s.corsMiddleware(h)
	h = s.requestIDMiddleware(h)
	h = s.recoveryMiddleware(h)
	h = s.loggingMiddleware(h)
	return h
}

func (s *Server) maxPageSize() int {
	if s.cfg != nil && s.cfg.MaxPageSize > 0 {
		return s.cfg.MaxPageSize
	}
	return 100
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		s.log.Infof("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				s.log.Errorf("panic: %v\n%s", rec, debug.Stack())
				httpx.InternalError(w, "服务器内部错误")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case model.IsValidationError(err):
		httpx.BadRequest(w, err.Error())
	case errors.Is(err, store.ErrNotFound):
		httpx.NotFound(w, err.Error())
	case errors.Is(err, store.ErrConflict):
		httpx.Conflict(w, err.Error())
	default:
		httpx.InternalError(w, err.Error())
	}
}
