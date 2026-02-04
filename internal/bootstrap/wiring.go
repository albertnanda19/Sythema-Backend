package bootstrap

import (
	"context"
	"database/sql"
	"time"

	redisadapter "synthema/internal/adapters/redis"
	"synthema/internal/app/health"
	"synthema/internal/config"
	authctx "synthema/internal/context"
	authhandlers "synthema/internal/handlers/auth"
	"synthema/internal/http"
	"synthema/internal/middleware"
	"synthema/internal/observability"
	"synthema/internal/repositories"
	"synthema/internal/routes"
	"synthema/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
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

	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	authMW := middleware.Auth(userRepo, sessionRepo, cfg.Auth.CookieName)

	authService := service.NewAuthService(userRepo, sessionRepo, cfg.Auth.SessionTTL)

	cookieSecure := cfg.Auth.CookieSecure
	cookieSameSite := cfg.Auth.CookieSameSite
	cookieDomain := cfg.Auth.CookieDomain
	// Dev-friendly overrides
	if cfg.Environment == "dev" {
		cookieSecure = false
		if cookieSameSite == "Strict" {
			cookieSameSite = "Lax"
		}
		cookieDomain = ""
	}

	authHandler := http.NewAuthHandler(authService, cfg.Auth.CookieName, cookieSecure, cookieSameSite, cookieDomain)

	app := fiber.New(fiber.Config{ErrorHandler: http.FiberErrorHandler(logger)})
	app.Use(middleware.RequestLogger(logger))
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool { return true },
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "*",
		MaxAge:           3600,
	}))
	app.Use(func(c *fiber.Ctx) error {
		err := c.Next()
		if c.Get("Origin") == "" {
			c.Set("Access-Control-Allow-Origin", "*")
		}
		return err
	})
	app.Use(recover.New())

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

	meHandler := authhandlers.NewMeHandler()
	logoutMW := middleware.Logout(cfg.Auth.CookieName)
	logoutHandler := authhandlers.NewLogoutHandler(authService, cfg.Auth.CookieName, cookieSecure, cookieSameSite, cookieDomain)
	routes.RegisterAuthRoutes(v1, authMW, meHandler, logoutMW, logoutHandler)

	api := v1.Group("", authMW)
	api.Get("/protected", func(c *fiber.Ctx) error {
		userID, ok := authctx.UserID(c.UserContext())
		if !ok {
			return http.Success(c, fiber.StatusOK, http.MsgProtectedOK, fiber.Map{"user_id": nil})
		}
		return http.Success(c, fiber.StatusOK, http.MsgProtectedOK, fiber.Map{"user_id": userID})
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
