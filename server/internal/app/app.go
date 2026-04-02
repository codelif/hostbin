package app

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"

	"hostbin/internal/admin"
	"hostbin/internal/auth"
	"hostbin/internal/clock"
	"hostbin/internal/config"
	"hostbin/internal/httpserver"
	"hostbin/internal/logging"
	"hostbin/internal/nonce"
	publicpkg "hostbin/internal/public"
	"hostbin/internal/router"
	"hostbin/internal/storage/sqlite"
)

type Options struct {
	Clock  clock.Clock
	Logger *zap.Logger
}

type App struct {
	Server *http.Server
	logger *zap.Logger
	db     *sql.DB
}

func New(cfg config.Config, opts Options) (*App, error) {
	logger := opts.Logger
	if logger == nil {
		var err error
		logger, err = logging.New(cfg.LogLevel)
		if err != nil {
			return nil, err
		}
	}

	appClock := opts.Clock
	if appClock == nil {
		appClock = clock.System{}
	}

	db, err := sqlite.Open(cfg.DBPath)
	if err != nil {
		return nil, err
	}

	documentStore := sqlite.NewDocumentStore(db)
	publicService := publicpkg.NewService(documentStore)
	publicHandler := publicpkg.NewHandler(publicService)

	adminService := admin.NewService(documentStore, cfg.BaseDomain, appClock)
	adminHandler := admin.NewHandler(adminService, cfg.ReservedSet)

	nonceStore := nonce.NewMemoryStore(cfg.NonceTTL)
	authVerifier := auth.NewVerifier(cfg.AdminHost, []byte(cfg.PresharedKey), appClock, cfg.AuthTimestampSkew, nonceStore)

	publicEngine := httpserver.NewPublicEngine(publicHandler)
	adminEngine := httpserver.NewAdminEngine(adminHandler, cfg.MaxDocSize, authVerifier.Middleware())
	dispatcher := router.NewDispatcher(cfg.BaseDomain, cfg.AdminHost, cfg.ReservedSet, adminEngine, publicEngine)

	handler := httpserver.RequestID(
		logging.Middleware(logger, cfg.TrustProxyHeaders, cfg.TrustedProxyNets)(
			httpserver.Recovery(logger, dispatcher),
		),
	)

	server := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &App{
		Server: server,
		logger: logger,
		db:     db,
	}, nil
}

func (a *App) Close() error {
	var errs []error
	if a.db != nil {
		errs = append(errs, a.db.Close())
	}
	if a.logger != nil {
		errs = append(errs, a.logger.Sync())
	}
	return errors.Join(errs...)
}
