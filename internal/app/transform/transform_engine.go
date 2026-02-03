package transform

import "synthema/internal/observability"

type Engine struct {
	logger *observability.Logger
}

func NewEngine(logger *observability.Logger) *Engine {
	return &Engine{logger: logger}
}
