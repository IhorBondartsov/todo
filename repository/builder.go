package repository

import (
	"context"
	"to-do/api"
)

type TODOStorage interface {
	CreateToDo(ctx context.Context, todo api.ToDo) error
	UpdateToDo(ctx context.Context, todo api.ToDo) error
	DeleteToDo(ctx context.Context, todoID int64) error
	GetToDo(ctx context.Context, todoID int64) (*api.ToDo, error)
}

type UserStorage interface {
	GetUser(ctx context.Context, id int64) (*api.User, error)
}

type Storage interface {
	UserStorage
	TODOStorage
}

func NewDBClient(ctx context.Context, cfg StorageConfig) (Storage, error) {
	db := pgDatabase{
		cgf: cfg,
	}

	if err := db.initializeDatabase(ctx); err != nil {
		return nil, err
	}
	return &db, nil
}
