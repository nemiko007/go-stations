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

const (
	// TODO一覧を取得するSQL
	read = `SELECT id, subject, description, created_at, updated_at FROM todos LIMIT ? OFFSET ?`
)

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
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
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

	var (
		subj      string
		desc      string
		createdAt sql.NullTime
		updatedAt sql.NullTime
	)

	row := s.db.QueryRowContext(ctx, confirm, lastID)
	if err := row.Scan(&subj, &desc, &createdAt, &updatedAt); err != nil {
		return nil, err
	}

	todo := &model.TODO{
		ID:          int(lastID),
		Subject:     subj,
		Description: desc,
	}
	if createdAt.Valid {
		todo.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		todo.UpdatedAt = updatedAt.Time
	}

	return todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, limit, offset int64) ([]*model.TODO, error) {
	rows, err := s.db.QueryContext(ctx, read, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*model.TODO
	for rows.Next() {
		todo := &model.TODO{}
		if err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
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
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return &model.ErrNotFound{}
	}
	return nil
}

func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
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

	var todo model.TODO
	row := s.db.QueryRowContext(ctx, selectTODOByIDQuery, id)
	if err := row.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
		return nil, err
	}

	return &todo, nil
}
