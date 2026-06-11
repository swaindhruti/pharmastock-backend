package app

import (
	"context"
	"time"

	"github.com/labstack/echo/v5"
	"go.uber.org/zap"

	"github.com/swaindhruti/pharmastock-backend/internal/config"
	"github.com/swaindhruti/pharmastock-backend/internal/database"
	"github.com/swaindhruti/pharmastock-backend/internal/health"
	"github.com/swaindhruti/pharmastock-backend/internal/middleware"
	"github.com/swaindhruti/pharmastock-backend/internal/retailer"
	"github.com/swaindhruti/pharmastock-backend/internal/router"
	"github.com/swaindhruti/pharmastock-backend/internal/stockist"
)

type App struct {
	Config   *config.Config
	Database *database.PostgresDB
	Echo     *echo.Echo
	Logger   *zap.Logger
}

func NewApp() (*App, error) {

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		return nil, err
	}

	var logger *zap.Logger
	if cfg.AppEnv == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		return nil, err
	}

	e := echo.New()

	e.Use(middleware.RequestID(), middleware.Logger(logger), middleware.Recovery())

	stockistHandler := stockist.NewModule(db.Pool)
	retailerHandler := retailer.NewModule(db.Pool)

	handlers := &router.Handlers{
		Health:   health.NewHandler(db),
		Stockist: stockistHandler,
		Retailer: retailerHandler,
	}

	router.RegisterRoutes(e, handlers)

	return &App{
		Config:   cfg,
		Database: db,
		Echo:     e,
		Logger:   logger,
	}, nil
}

func (a *App) Start(ctx context.Context) error {
	sc := echo.StartConfig{
		Address:         ":" + a.Config.AppPort,
		GracefulTimeout: 10 * time.Second,
	}
	if err := sc.Start(ctx, a.Echo); err != nil {
		a.Logger.Error("failed to start server", zap.Error(err))
		return err
	}
	return nil
}

func (a *App) Shutdown() {
	a.Logger.Sync()
	a.Database.Close()
}
