package app

import (
	"github.com/google/wire"
	_ "github.com/jackc/pgx/stdlib"
	"gitverse.ru/icyre/template/internal/config"
	server "gitverse.ru/icyre/template/internal/transport/grpc"
)

func New() (*App, func(), error) {
	panic(wire.Build(
		newApp,
		server.New,

		config.New,
	))
}
