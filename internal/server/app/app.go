package app

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/codelif/hostbin/internal/clock"
	"github.com/codelif/hostbin/internal/server/adminauth"
	"github.com/codelif/hostbin/internal/server/adminhttp"
	serverconfig "github.com/codelif/hostbin/internal/server/config"
	"github.com/codelif/hostbin/internal/server/dispatch"
	"github.com/codelif/hostbin/internal/server/documentsvc"
	"github.com/codelif/hostbin/internal/server/logging"
	"github.com/codelif/hostbin/internal/server/middleware"
	"github.com/codelif/hostbin/internal/server/nonce"
	"github.com/codelif/hostbin/internal/server/publichttp"
	"github.com/codelif/hostbin/internal/server/store/sqlite"
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

func New(cfg serverconfig.Config, opts Options) (*App, error) {
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
	documentService := documentsvc.New(documentStore, appClock)
	publicHandler := publichttp.NewHandler(documentService)

	adminHandler := adminhttp.NewHandler(documentService, cfg.BaseDomain, cfg.ReservedSet)

	nonceStore := nonce.NewMemoryStore(cfg.NonceTTL)
	authVerifier := adminauth.NewVerifier(cfg.AdminHost, []byte(cfg.PresharedKey), appClock, cfg.AuthTimestampSkew, nonceStore)

	publicEngine := publichttp.NewEngine(publicHandler)
	adminEngine := adminhttp.NewEngine(adminHandler, cfg.MaxDocSize, authVerifier.Middleware())
	dispatcher := dispatch.NewHandler(cfg.BaseDomain, cfg.AdminHost, cfg.ReservedSet, adminEngine, publicEngine)

	handler := middleware.RequestID(
		logging.Middleware(logger, cfg.TrustProxyHeaders, cfg.TrustedProxyNets)(
			middleware.Recovery(logger, dispatcher),
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
