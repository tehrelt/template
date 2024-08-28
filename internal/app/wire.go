//go:build wire

package app

import (
	"fmt"
	"log/slog"

	"github.com/google/wire"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"gitverse.ru/icyre/template/internal/config"
	"gitverse.ru/icyre/template/internal/storage/pg"
	server "gitverse.ru/icyre/template/internal/transport/grpc"
)

func New() (*App, func(), error) {
	panic(wire.Build(
		newApp,
		server.New,
		pg.NewUserStorage,

		initPG,
		config.New,

		wire.Bind(new(server.Repository), new(*pg.UserStorage)),
	))
}

func initPG(cfg *config.Config) (*sqlx.DB, func(), error) {
	host := cfg.Pg.Host
	port := cfg.Pg.Port
	user := cfg.Pg.User
	pass := cfg.Pg.Pass
	name := cfg.Pg.Name

	cs := fmt.Sprintf(`postgres://%s:%s@%s:%d/%s?sslmode=disable`, user, pass, host, port, name)

	slog.Info("connecting to database", slog.String("conn", cs))

	db, err := sqlx.Connect("pgx", cs)
	if err != nil {
		return nil, nil, err
	}

	slog.Info("send ping to database")

	if err := db.Ping(); err != nil {
		slog.Error("failed to connect to database", slog.String("err", err.Error()), slog.String("conn", cs))
		return nil, func() { db.Close() }, err
	}

	slog.Info("connected to database", slog.String("conn", cs))

	return db, func() { db.Close() }, nil
}
