package server

import (
	"log/slog"
)

type Repository interface {
}

type Handler struct {
	// CHANGE IT
	// userspb.UnimplementedUserServiceServer
	repository Repository
	logger     *slog.Logger
}

func New(repo Repository) *Handler {
	return &Handler{
		repository: repo,
		logger:     slog.Default().With(slog.String("struct", "Handler")),
	}
}
