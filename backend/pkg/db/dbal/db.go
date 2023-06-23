package dbal

import (
	"context"

	"github.com/JakeDoughty/customer-io-homework-backend/pkg/db/entities"
	"gorm.io/gorm"
)

var db_Key any = &struct{}{}

func UseDatabase(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, db_Key, db.WithContext(ctx))
}
func GetDatabase(ctx context.Context) *gorm.DB {
	if db, ok := ctx.Value(db_Key).(*gorm.DB); ok {
		return db
	} else {
		return nil
	}
}
func AutoMigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(&entities.Application{}, &entities.Session{}, &entities.Event{})
}
