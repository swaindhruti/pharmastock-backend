package app

import (
	"context"
	"time"

	"github.com/labstack/echo/v5"
	"go.uber.org/zap"

	"github.com/swaindhruti/pharmastock-backend/internal/auth"
	"github.com/swaindhruti/pharmastock-backend/internal/config"
	"github.com/swaindhruti/pharmastock-backend/internal/database"
	"github.com/swaindhruti/pharmastock-backend/internal/health"
	"github.com/swaindhruti/pharmastock-backend/internal/inventory"
	"github.com/swaindhruti/pharmastock-backend/internal/job"
	"github.com/swaindhruti/pharmastock-backend/internal/medicine"
	"github.com/swaindhruti/pharmastock-backend/internal/middleware"
	"github.com/swaindhruti/pharmastock-backend/internal/retailer"
	"github.com/swaindhruti/pharmastock-backend/internal/router"
	"github.com/swaindhruti/pharmastock-backend/internal/stockist"
	"github.com/swaindhruti/pharmastock-backend/internal/upload"
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

	// Stockist module
	stockistModule := stockist.NewModule(db.Pool)

	// Retailer module
	retailerModule := retailer.NewModule(db.Pool)

	// Medicine module
	medicineHandler := medicine.NewModule(db.Pool)

	// Inventory module
	inventoryHandler := inventory.NewModule(db.Pool)

	// Auth module
	authHandler := auth.NewModule(db.Pool, cfg.JWTSecret, stockistModule.Service, retailerModule.Service)

	// Upload module
	jobRepo := job.NewRepository(db.Pool)
	jobSvc := job.NewService(jobRepo, &noopProcessor{})
	uploadSvc := upload.NewService(jobSvc, cfg.UploadDir)
	uploadHandler := upload.NewHandler(uploadSvc)

	handlers := &router.Handlers{
		Auth:      authHandler,
		Health:    health.NewHandler(db),
		Stockist:  stockistModule.Handler,
		Retailer:  retailerModule.Handler,
		Medicine:  medicineHandler,
		Inventory: inventoryHandler,
		Upload:    uploadHandler,
	}

	router.RegisterRoutes(e, handlers, cfg.JWTSecret)

	// Seed admin user from env
	if err := authHandler.SeedAdmin(context.Background(), cfg.AdminUsername, cfg.AdminPassword, cfg.AdminEmail); err != nil {
		logger.Warn("failed to seed admin user", zap.Error(err))
	}

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

type noopProcessor struct{}

func (n *noopProcessor) Process(_ context.Context, _ *job.Job) error {
	return nil
}
