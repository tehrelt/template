package app

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"

	"gitverse.ru/icyre/template/internal/config"
	"gitverse.ru/icyre/template/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	handler *server.Handler
	cfg     *config.Config
}

func newApp(cfg *config.Config, handler *server.Handler) *App {
	return &App{
		handler, cfg,
	}
}

func (a *App) Run() {

	slog.Info("running server")

	server := grpc.NewServer()
	if a.cfg.App.UseReflection {
		slog.Info("enabling reflection")
		reflection.Register(server)
	}

	slog.Info("registering user service")
	// CHANGE IT
	// userspb.RegisterUserServiceServer(server, a.handler)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	go func() {
		addr := fmt.Sprintf("%s:%d", a.cfg.App.Host, a.cfg.App.Port)
		slog.Debug("starting listener", slog.String("addr", addr))
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			panic(fmt.Errorf("cannot bind port %d: %w", a.cfg.App.Port, err))
		}

		if err := server.Serve(listener); err != nil {
			panic(fmt.Errorf("serve error", err))
		}

	}()

	sig := <-sigChan
	slog.Info(fmt.Sprintf("Signal %v received, stopping server...\n", sig))
	server.GracefulStop()
}
