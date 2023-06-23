package dbal

import (
	"context"
	"time"

	"github.com/JakeDoughty/customer-io-homework-backend/pkg/db/entities"
	"github.com/JakeDoughty/customer-io-homework-backend/pkg/db/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var includeClosedSessions_Key any = "includeClosedSessions"

func SessionsBefore(t time.Time) DBOption {
	return func(c context.Context, db *gorm.DB) (context.Context, *gorm.DB) {
		return c, db.Where("createdAt < (?)", t.UTC())
	}
}
func SessionsAfter(t time.Time) DBOption {
	return func(c context.Context, db *gorm.DB) (context.Context, *gorm.DB) {
		return c, db.Where("createdAt > (?)", t.UTC())
	}
}
func SessionsBetween(start, stop time.Time) DBOption {
	return func(c context.Context, db *gorm.DB) (context.Context, *gorm.DB) {
		return c, db.Where("createdAt BETWEEN (?) AND (?)", start.UTC(), stop.UTC())
	}
}
func IncludeClosedSessions() DBOption {
	return func(c context.Context, db *gorm.DB) (context.Context, *gorm.DB) {
		return context.WithValue(c, includeClosedSessions_Key, true), db
	}
}

func GetAllSessions(
	ctx context.Context,
	applicationID *uuid.UUID,
	pagination *Pagination,
	options ...DBOption,
) (List[entities.Session], error) {
	db := GetDatabase(ctx)
	ctx, db = ApplyDBOptions(ctx, db, options...)
	if includeClosedSessions, _ := ctx.Value(includeClosedSessions_Key).(bool); !includeClosedSessions {
		db = db.Where("closedAt > ?", time.Now().UTC())
	}

	var list List[entities.Session]
	if applicationID != nil {
		db = db.Where("applicationID=?", *applicationID)
	}
	err := list.ApplyPagination(db, pagination).Find(&list.Items).Error
	return list, err
}
func GetSessionByID(ctx context.Context, id uuid.UUID, options ...DBOption) (*entities.Session, error) {
	db := GetDatabase(ctx)
	_, db = ApplyDBOptions(ctx, db, options...)

	var session *entities.Session
	if err := db.First(&session, id).Error; err != nil {
		return nil, err
	} else if session == nil || utils.IsZeroUUID(session.ID) {
		return nil, nil
	} else {
		return session, nil
	}
}
func CreateSession(ctx context.Context, session *entities.Session, options ...DBOption) (*entities.Session, error) {
	db := GetDatabase(ctx)
	_, db = ApplyDBOptions(ctx, db, options...)

	if err := db.Create(session).Error; err != nil {
		return nil, err
	} else {
		return session, nil
	}
}
func UpdateSession(ctx context.Context, session *entities.Session, options ...DBOption) (*entities.Session, error) {
	db := GetDatabase(ctx)
	_, db = ApplyDBOptions(ctx, db, options...)

	if err := db.Model(session).Clauses(clause.Returning{}).
		Updates(session).Error; err != nil {
		return nil, err
	} else {
		return session, nil
	}
}
func CloseSession(ctx context.Context, sessionID uuid.UUID, options ...DBOption) (*entities.Session, error) {
	db := GetDatabase(ctx)
	_, db = ApplyDBOptions(ctx, db, options...)

	var session *entities.Session
	if err := db.Model(&session).Clauses(clause.Returning{}).
		Where("ID", sessionID).Update("ExpiredTime", time.Now().UTC()).Error; err != nil {
		return nil, err
	} else {
		return session, nil
	}
}
