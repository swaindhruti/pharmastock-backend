package app

import (
	"github.com/labstack/echo/v5"

	"github.com/swaindhruti/pharmastock-backend/internal/config"
	"github.com/swaindhruti/pharmastock-backend/internal/database"
	"github.com/swaindhruti/pharmastock-backend/internal/health"
	"github.com/swaindhruti/pharmastock-backend/internal/router"
)

type App struct {
	Config   *config.Config
	Database *database.PostgresDB
	Echo     *echo.Echo
}

func NewApp(cfg *config.Config) (*App, error) {

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		return nil, err
	}

	e := echo.New()

	handlers := &router.Handlers{
		Health: health.NewHandler(db),
	}

	router.RegisterRoutes(e, handlers)

	return &App{
		Config:   cfg,
		Database: db,
		Echo:     e,
	}, nil
}
