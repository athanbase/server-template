//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"server-template/internal/biz"
	"server-template/internal/conf"
	"server-template/internal/data"
	"server-template/internal/server"
	"server-template/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Config, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet,
		wire.FieldsOf(new(*conf.Config), "Server", "Redis", "Db"),
		newApp,
	))
}
