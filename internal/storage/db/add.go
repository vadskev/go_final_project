package db

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/vadskev/go_final_project/internal/logger"
	"github.com/vadskev/go_final_project/internal/models/task"
	"go.uber.org/zap"
)

func (r *Repository) Create(task *task.Info) (int64, error) {
	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(colDate, colTitle, colComment, colRepeat).
		Values(task.Date, task.Title, task.Comment, task.Repeat).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		logger.Error("storage.db.Create", zap.Any("error:", err.Error()))
		return 0, err
	}

	stmt, err := r.DB().Prepare(query)
	if err != nil {
		logger.Error("storage.db.Create", zap.Any("error:", err.Error()))
		return 0, err
	}

	defer stmt.Close()

	res, err := stmt.Exec(args...)
	if err != nil {
		logger.Error("storage.db.Create", zap.Any("error:", err.Error()))
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		logger.Error("storage.db.Create", zap.Any("error:", err.Error()))
		return 0, err
	}
	return id, nil
}
