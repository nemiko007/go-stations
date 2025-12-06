package model

// ErrNotFound は TODO が存在しない場合に返されるエラー
type ErrNotFound struct{}

func (e *ErrNotFound) Error() string {
	return "todo not found"
}
