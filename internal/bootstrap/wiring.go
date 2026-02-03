package bootstrap

import (
	"database/sql"

	"synthema/internal/config"
	"synthema/internal/http"
	"synthema/internal/observability"
	"synthema/internal/repository"
	"synthema/internal/service"

	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type APIApp struct {
	Config config.Config
	Logger *observability.Logger
	App    *fiber.App
}

type CaptureApp struct {
	Config config.Config
	Logger *observability.Logger
}

type WorkerApp struct {
	Config config.Config
	Logger *observability.Logger
}

func BootstrapAPI() (APIApp, error) {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		return APIApp{}, err
	}
	logger := observability.NewLogger(cfg)

	db, err := sql.Open("pgx", cfg.Postgres.DSN)
	if err != nil {
		return APIApp{}, err
	}

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	authService := service.NewAuthService(userRepo, sessionRepo)

	authHandler := http.NewAuthHandler(authService)

	app := fiber.New()

	app.Post("/auth/login", authHandler.Login)
	app.Post("/auth/logout", authHandler.Logout)

	api := app.Group("/api", http.AuthMiddleware(sessionRepo))
	api.Get("/protected", func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(int64)
		return c.JSON(fiber.Map{"message": "welcome to the protected area", "user_id": userID})
	})

	return APIApp{Config: cfg, Logger: logger, App: app}, nil
}

func BootstrapCapture() (CaptureApp, error) {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		return CaptureApp{}, err
	}
	logger := observability.NewLogger(cfg)
	return CaptureApp{Config: cfg, Logger: logger}, nil
}

func BootstrapWorker() (WorkerApp, error) {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		return WorkerApp{}, err
	}
	logger := observability.NewLogger(cfg)
	return WorkerApp{Config: cfg, Logger: logger}, nil
}
