package app

import (
	"context"
	"to-do/api"
	"to-do/repository"

	log "github.com/sirupsen/logrus"
)

type ToDoService struct {
	db repository.Storage
}

func NewToDoService(db repository.Storage) (*ToDoService, error) {
	return &ToDoService{db: db}, nil
}

func (t *ToDoService) CreateToDo(ctx context.Context, todo api.ToDo) error {
	err := t.db.CreateToDo(ctx, todo)
	if err != nil {
		log.Error("cant create new todo: ", err)
		return err
	}
	return nil
}

func (t *ToDoService) UpdateToDo(ctx context.Context, todo api.ToDo) error {
	err := t.db.UpdateToDo(ctx, todo)
	if err != nil {
		log.Error("cant update todo: ", err)
		return err
	}
	return nil
}

func (t *ToDoService) GetTodo(ctx context.Context, todoID int64) (*api.ToDo, error) {
	todo, err := t.db.GetToDo(ctx, todoID)
	// TODO: handle db error NotFound
	if err != nil {
		log.Error("cant return todo: ", err)
		return nil, err
	}
	return todo, nil
}

func (t *ToDoService) DeleteTodo(ctx context.Context, todoID int64) error {
	err := t.db.DeleteToDo(ctx, todoID)
	if err != nil {
		log.Error("cant delete todo: ", err)
		return err
	}
	return nil
}
