package validator

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/config"
	"github.com/brienze1/crypto-robot-validator/internal/validator/delivery/adapters"
)

// Main class works as a proxy for the handler.Handler class. It's responsible for configuring env vars with
// config.LoadEnv and injecting dependencies with config.DependencyInjector before passing the request forward.
func Main() adapters.HandlerAdapter {
	config.LoadEnv()
	return config.DependencyInjector().WireDependencies().Handler
}
