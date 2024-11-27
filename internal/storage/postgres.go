package storage

import (
	"context"
	"fmt"
	"github.com/hyperfyodor/yq_sample/db/postgres"
	"github.com/hyperfyodor/yq_sample/internal/helpers"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	db      *pgxpool.Pool
	queries *postgres.Queries
}

func (p *PostgresStorage) Close() {
	p.db.Close()
}

func NewPostgresStorage(ctx context.Context, connection string, ping bool) (*PostgresStorage, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		db, err := pgxpool.New(ctx, connection)

		if err != nil {
			return &PostgresStorage{}, err
		}

		if ping {
			if err := db.Ping(ctx); err != nil {
				return &PostgresStorage{}, err
			}
		}

		queries := postgres.New(db)

		return &PostgresStorage{db, queries}, nil
	}
}

func (p *PostgresStorage) SaveTask(ctx context.Context, taskType int, taskValue int) (int, error) {
	const op = "internal.storage.postgres.SaveTask"

	select {
	case <-ctx.Done():
		return 0, helpers.WrapErr(op, ctx.Err())
	default:
		id, err := p.queries.CreateTask(ctx, postgres.CreateTaskParams{Type: int32(taskType), Value: int32(taskValue)})

		if err != nil {
			return 0, helpers.WrapErr(op, err)
		}

		return int(id), nil
	}
}

func (p *PostgresStorage) State(ctx context.Context, id int) (string, error) {
	const op = "internal.storage.postgres.State"

	select {
	case <-ctx.Done():
		return "", helpers.WrapErr(op, ctx.Err())
	default:
		state, err := p.queries.GetTaskState(ctx, int32(id))

		if err != nil {
			return "", helpers.WrapErr(op, err)
		}

		return state, nil
	}
}

func (p *PostgresStorage) Done(ctx context.Context, id int) error {
	const op = "internal.storage.postgres.Done"

	select {
	case <-ctx.Done():
		return helpers.WrapErr(op, ctx.Err())
	default:
		id2, err := p.queries.SetStateToDone(ctx, int32(id))

		if err != nil {
			return helpers.WrapErr(op, err)
		}

		if int32(id) != id2 {
			return helpers.WrapErr(op, fmt.Errorf("id=%v not found", id))
		}

		return nil
	}
}

func (p *PostgresStorage) Processing(ctx context.Context, id int) error {
	const op = "internal.storage.postgres.Processing"

	select {
	case <-ctx.Done():
		return helpers.WrapErr(op, ctx.Err())
	default:
		id2, err := p.queries.SetStateToProcessing(ctx, int32(id))

		if err != nil {
			return helpers.WrapErr(op, err)
		}

		if int32(id) != id2 {
			return helpers.WrapErr(op, fmt.Errorf("id=%v not found", id))
		}

		return nil

	}
}
