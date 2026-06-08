package app

import (
	"github.com/labstack/echo/v5"

	"github.com/swaindhruti/pharmastock-backend/internal/config"
	"github.com/swaindhruti/pharmastock-backend/internal/database"
	"github.com/swaindhruti/pharmastock-backend/internal/health"
	"github.com/swaindhruti/pharmastock-backend/internal/middleware"
	"github.com/swaindhruti/pharmastock-backend/internal/router"
	"github.com/swaindhruti/pharmastock-backend/internal/stockist"
)

type App struct {
	Config   *config.Config
	Database *database.PostgresDB
	Echo     *echo.Echo
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

	e := echo.New()

	e.Use(middleware.RequestID(), middleware.Logger("api"), middleware.Recovery())

	stockistHandler := stockist.NewModule(db.Pool)

	handlers := &router.Handlers{
		Health:   health.NewHandler(db),
		Stockist: stockistHandler,
	}

	router.RegisterRoutes(e, handlers)

	return &App{
		Config:   cfg,
		Database: db,
		Echo:     e,
	}, nil
}

func (a *App) Start() error {
	return a.Echo.Start(":" + a.Config.AppPort)
}

func (a *App) Shutdown() {
	a.Database.Close()
}
