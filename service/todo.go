package service

import (
	"context"
	"database/sql"

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
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	return nil, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	return nil, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	return nil
}
