package dbal

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type DBOption = func(c context.Context, db *gorm.DB) (context.Context, *gorm.DB)

func ApplyDBOptions(c context.Context, db *gorm.DB, options ...DBOption) (context.Context, *gorm.DB) {
	for _, opt := range options {
		c, db = opt(c, db)
	}
	return c, db
}

func WithSessions(includeClosedSessions bool) DBOption {
	return func(c context.Context, db *gorm.DB) (context.Context, *gorm.DB) {
		if includeClosedSessions {
			return c, db.Preload("Sessions")
		} else {
			return c, db.Preload("Sessions", "closedAt > ?", time.Now().UTC())
		}
	}
}
func WithSessionEvents() DBOption {
	return func(c context.Context, db *gorm.DB) (context.Context, *gorm.DB) {
		return c, db.Preload("Sessions.Events")
	}
}
func WithEvents() DBOption {
	return func(c context.Context, db *gorm.DB) (context.Context, *gorm.DB) {
		return c, db.Preload("Events")
	}
}
func WithApplication() DBOption {
	return func(c context.Context, db *gorm.DB) (context.Context, *gorm.DB) {
		return c, db.Preload("Application")
	}
}
func WithSession() DBOption {
	return func(c context.Context, db *gorm.DB) (context.Context, *gorm.DB) {
		return c, db.Preload("Session")
	}
}
