package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

const (
	// TODO を更新する SQL
	updateTODOQuery     = `UPDATE todos SET subject = ?, description = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	selectTODOByIDQuery = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id = ?`
)

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.Todo, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	// prepare statement
	stmt, err := s.db.PrepareContext(ctx, insert)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, subject, description)
	if err != nil {
		return nil, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	var todo model.Todo
	row := s.db.QueryRowContext(ctx, confirm, lastID)
	if err := row.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
		return nil, err
	}
	todo.ID = lastID

	return &todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.Todo, error) {
	const (
		readWithPrevID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id > ? ORDER BY id ASC LIMIT ?`
		readAll        = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id ASC LIMIT ?`
	)

	if size == 0 {
		size = 5
	}

	var rows *sql.Rows
	var err error

	if prevID > 0 {
		rows, err = s.db.QueryContext(ctx, readWithPrevID, prevID, size)
	} else {
		rows, err = s.db.QueryContext(ctx, readAll, size)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := make([]*model.Todo, 0)
	for rows.Next() {
		var todo model.Todo
		if err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

// DeleteTODO deletes TODOs on DB.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf("DELETE FROM todos WHERE id IN (%s)", strings.Join(placeholders, ","))
	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return &model.ErrNotFound{}
	}

	return nil
}

// UpdateTODO updates a TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.Todo, error) {
	res, err := s.db.ExecContext(ctx, updateTODOQuery, subject, description, id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, &model.ErrNotFound{}
	}

	var todo model.Todo
	row := s.db.QueryRowContext(ctx, selectTODOByIDQuery, id)
	if err := row.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
		return nil, err
	}

	return &todo, nil
}
