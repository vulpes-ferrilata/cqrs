package gorm

import (
	"context"

	"gorm.io/gorm"

	"github.com/vulpes-ferrilata/cqrs/pkg/db"
)

func newCommitter(db *gorm.DB) db.Committer {
	return &committer{
		db: db,
	}
}

type committer struct {
	db *gorm.DB
}

func (c committer) CommitTransaction(ctx context.Context) error {
	return c.db.Commit().Error
}

func (c committer) RollbackTransaction(ctx context.Context) error {
	return c.db.Rollback().Error
}
