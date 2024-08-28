package pg

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
	"gitverse.ru/icyre/template/internal/dto"
	"gitverse.ru/icyre/template/internal/entity"
	"gitverse.ru/icyre/template/internal/lib/logger/sl"
	"gitverse.ru/icyre/template/internal/server"
	"gitverse.ru/icyre/template/internal/storage"
)

var _ server.Repository = (*Storage)(nil)

type Storage struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func (u *Storage) Create(ctx context.Context, cu *dto.CreateEntity) (*entity.Entity, error) {

	log := u.logger.With(slog.String("method", "Create"))

	query, args, err := squirrel.Insert(table).
		// Columns("username", "email").
		// Values(cu.Username, cu.Email).
		Suffix("RETURNING *").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		log.Error("cannot build query", sl.Err(err))
		return nil, fmt.Errorf("cannot build query: %w", err)
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	user := new(entity.Entity)
	if err := u.db.Get(user, query, args...); err != nil {

		if e, ok := err.(pgx.PgError); ok {
			log.Debug("pg error", sl.PgError(e))
			if e.Code == "23505" {
				return nil, storage.ErrAlreadyExists
			}
		}

		log.Error("cannot execute query", sl.Err(err))
		return nil, fmt.Errorf("cannot execute query: %w", err)
	}

	log.Debug("user saved", slog.Any("user", user))

	return user, nil
}

func (u *Storage) Delete(ctx context.Context, id string) error {

	log := u.logger.With(slog.String("method", "Delete"))

	query, args, err := squirrel.Delete(table).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		log.Error("cannot build query", sl.Err(err))
		return fmt.Errorf("cannot build query: %w", err)
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	if _, err := u.db.Exec(query, args...); err != nil {
		if e, ok := err.(pgx.PgError); ok {
			log.Debug("pg error", sl.PgError(e))
			if e.Code == "22P02" {
				return storage.ErrNotFound
			}
		}

		log.Error("cannot execute query", sl.Err(err))
		return fmt.Errorf("cannot execute query: %w", err)
	}

	log.Info("user deleted", slog.String("userId", id))

	return nil
}

func (u *Storage) List(ctx context.Context, filters *dto.ListEntity) ([]*entity.Entity, error) {
	log := u.logger.With(slog.String("method", "List"))

	query, args, err := squirrel.Select("*").
		From(table).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		log.Error("cannot build query", sl.Err(err))
		return nil, fmt.Errorf("cannot build query: %w", err)
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	users := make([]*entity.Entity, 0)
	if err := u.db.Select(&users, query, args...); err != nil {
		log.Error("cannot execute query", sl.Err(err))
		return nil, fmt.Errorf("cannot execute query: %w", err)
	}

	log.Debug("users", slog.Any("users", users))

	return users, nil
}

func (u *Storage) Find(ctx context.Context, id string) (*entity.Entity, error) {
	log := u.logger.With(slog.String("method", "Find"))

	query, args, err := squirrel.Select("*").
		From(table).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		log.Error("cannot build query", sl.Err(err))
		return nil, fmt.Errorf("cannot build query: %w", err)
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	user := new(entity.Entity)
	if err := u.db.Get(user, query, args...); err != nil {
		if e, ok := err.(pgx.PgError); ok {
			log.Debug("pg error", sl.PgError(e))
			if e.Code == "22P02" {
				return nil, storage.ErrNotFound
			}
		}

		log.Error("cannot execute query", sl.Err(err))
		return nil, fmt.Errorf("cannot execute query: %w", err)

	}

	log.Debug("user found", slog.Any("user", user))

	return user, nil
}

func (u *Storage) Update(ctx context.Context, dto *dto.UpdateEntity) (*entity.Entity, error) {

	log := u.logger.With(slog.String("method", "Update"))

	builder := squirrel.
		Update(table).
		Where(squirrel.Eq{"id": dto.Id}).
		Suffix("RETURNING *").
		PlaceholderFormat(squirrel.Dollar)

		// FIELDS
	// if dto.Email != nil {
	// 	builder = builder.Set("email", dto.Email)
	// }

	// if dto.Username != nil {
	// 	builder = builder.Set("username", dto.Username)
	// }

	query, args, err := builder.ToSql()
	if err != nil {
		log.Error("cannot build query", sl.Err(err))
		return nil, fmt.Errorf("cannot build query: %w", err)
	}

	log.Debug("query", slog.String("query", query), slog.Any("args", args))

	user := new(entity.Entity)
	if err := u.db.Get(user, query, args...); err != nil {

		if e, ok := err.(pgx.PgError); ok {
			log.Debug("pg error", sl.PgError(e))
			if e.Code == "22P02" {
				return nil, storage.ErrNotFound
			}
		}

		log.Error("cannot execute query", sl.Err(err))
		return nil, fmt.Errorf("cannot execute query: %w", err)
	}

	log.Info("user updated", slog.Any("user", dto))

	return user, nil
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		db:     db,
		logger: slog.Default().With(slog.String("struct", "Storage")),
	}
}
