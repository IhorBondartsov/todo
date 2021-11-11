package repository

import (
	"context"
	"database/sql"
	"fmt"
	"to-do/api"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	// TODO_LIST table query
	deleteTODOQuery = `DELETE FROM todo_app.todo_list WHERE id = $1`

	addToDoQuery = `
		INSERT INTO todo_app.todo_list 
			(id, user_id, created_at, updated_at, message)
    	VALUES 
			(DEFAULT, $1, DEFAULT, DEFAULT, $2)`

	getToDoQuery = `SELECT id, user_id, created_at, updated_at, message FROM todo_app.todo_list WHERE id = $1`

	updateToDoQuery = `UPDATE todo_app.todo_list SET message=$1  WHERE id = $2`

	// TODO: implement
	getTODOsQuery = `SELECT * FROM todo_app.todo_list WHERE user_id = $1 LIMIT $2 OFFSET $3`

	// USERS Query
	getUserQuery = `SELECT * FROM todo_app.users WHERE id = $1`
)

type StorageConfig struct {
	Driver string `json:"driver"`
	DSN    string `json:"dsn"`
}

func (c StorageConfig) Validate() error {
	var errs []error
	if len(c.Driver) == 0 {
		errs = append(errs, errors.New("db driver cannot be empty"))
	}

	if len(c.DSN) == 0 {
		errs = append(errs, errors.New("db DSN cannot be empty"))
	}

	if len(errs) != 0 {
		return errors.Errorf("validate errors - %v", errs)
	}
	return nil
}

type pgDatabase struct {
	cgf StorageConfig
	db  *sql.DB
}

func (pg *pgDatabase) initializeDatabase(ctx context.Context) error {
	db, err := sql.Open(pg.cgf.Driver, pg.cgf.DSN)
	if err != nil {
		return errors.Wrapf(err, "open(%s) database connection", pg.cgf.Driver)
	}
	pg.db = db
	if err := pg.db.PingContext(ctx); err != nil {
		return errors.Wrap(err, "ping database")
	}
	log.Info("Successfully connected to database.")
	return nil
}

func (pg *pgDatabase) CreateToDo(ctx context.Context, todo api.ToDo) error {
	result, err := pg.db.ExecContext(ctx, addToDoQuery, todo.UserID, todo.Message)
	if err != nil {
		return errors.Wrap(err, "insert todo to database")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "insert todo to database, cant return rows affected")
	}

	if rows != 1 {
		return fmt.Errorf("expected to affect 1 row, affected %d", rows)
	}
	log.Debugf("Successfully inserted materialization instance to database.")
	return nil
}

func (pg *pgDatabase) UpdateToDo(ctx context.Context, todo api.ToDo) error {
	result, err := pg.db.ExecContext(ctx, updateToDoQuery, todo.Message, todo.ID)
	if err != nil {
		return errors.Wrap(err, "update todo in database")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "update todo in database, cant return rows affected")
	}

	if rows != 1 {
		return fmt.Errorf("expected to affect 1 row, affected %d", rows)
	}
	log.Debugf("Successfully updated materialization instance to database.")
	return nil
}

func (pg *pgDatabase) DeleteToDo(ctx context.Context, todoID int64) error {
	_, err := pg.db.ExecContext(ctx, deleteTODOQuery, todoID)
	if err != nil {
		return errors.Wrap(err, "delete spec data in database")
	}
	log.Debugf("Successfully deleted todo in database.")

	return nil
}

func (pg *pgDatabase) GetToDo(ctx context.Context, todoID int64) (*api.ToDo, error) {
	var todo api.ToDo
	err := pg.db.QueryRowContext(ctx, getToDoQuery, todoID).Scan(
		&todo.ID,
		&todo.UserID,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.Message)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		log.Printf("no todo with id %d\n", todoID)
		return nil, nil
	case err != nil:
		return nil, errors.Wrap(err, "query error")
	}
	return &todo, nil
}

func (pg *pgDatabase) GetUser(ctx context.Context, id int64) (*api.User, error) {
	var user api.User
	err := pg.db.QueryRowContext(ctx, getUserQuery, id).Scan(
		&user.ID,
		&user.Name)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		log.Printf("no todo with id %d\n", user.ID)
		return nil, nil
	case err != nil:
		return nil, errors.Wrap(err, "query error")
	}
	return &user, nil
}
