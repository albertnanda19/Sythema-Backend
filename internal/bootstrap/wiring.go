package bootstrap

import (
	"context"
	"database/sql"
	"time"

	redisadapter "synthema/internal/adapters/redis"
	"synthema/internal/app/health"
	"synthema/internal/config"
	"synthema/internal/http"
	"synthema/internal/observability"
	"synthema/internal/repository"
	"synthema/internal/service"

	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
)

type APIApp struct {
	Config config.Config
	Logger *observability.Logger
	App    *fiber.App
	DB     *sql.DB
	Redis  *redis.Client
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
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return APIApp{}, err
	}

	redisClient, err := redisadapter.NewClient(cfg.Redis)
	if err != nil {
		_ = db.Close()
		return APIApp{}, err
	}
	if err := redisadapter.Ping(context.Background(), redisClient, cfg.Redis.DialTimeout); err != nil {
		_ = redisClient.Close()
		_ = db.Close()
		return APIApp{}, err
	}

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	authService := service.NewAuthService(userRepo, sessionRepo, cfg.Auth.SessionTTL)

	authHandler := http.NewAuthHandler(authService, cfg.Auth.CookieName, cfg.Auth.CookieSecure)

	app := fiber.New(fiber.Config{ErrorHandler: http.FiberErrorHandler(logger)})

	v1 := app.Group("/api/v1")

	postgresCheck := func(ctx context.Context) error {
		return db.PingContext(ctx)
	}
	redisCheck := func(ctx context.Context) error {
		return redisadapter.Ping(ctx, redisClient, cfg.Redis.ReadTimeout)
	}
	healthHandler := health.NewHandler(postgresCheck, redisCheck, time.Second)
	v1.Get("/health", healthHandler.Health)
	app.Get("/healthz", healthHandler.Healthz)

	v1.Post("/auth/login", authHandler.Login)
	v1.Post("/auth/logout", http.AuthMiddleware(userRepo, sessionRepo, cfg.Auth.CookieName), authHandler.Logout)

	api := v1.Group("", http.AuthMiddleware(userRepo, sessionRepo, cfg.Auth.CookieName))
	api.Get("/protected", func(c *fiber.Ctx) error {
		return http.Success(c, fiber.StatusOK, http.MsgProtectedOK, fiber.Map{"user_id": c.Locals("userID")})
	})

	return APIApp{Config: cfg, Logger: logger, App: app, DB: db, Redis: redisClient}, nil
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
