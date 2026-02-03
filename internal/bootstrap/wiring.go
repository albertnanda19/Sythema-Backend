package bootstrap

import (
	"synthema/internal/config"
	"synthema/internal/observability"
)

type APIApp struct {
	Config config.Config
	Logger *observability.Logger
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
	return APIApp{Config: cfg, Logger: logger}, nil
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
