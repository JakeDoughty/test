package dbal

import (
	"context"

	"github.com/JakeDoughty/customer-io-homework-backend/pkg/db/entities"
	"github.com/JakeDoughty/customer-io-homework-backend/pkg/db/utils"
	"github.com/google/uuid"
)

func GetAllEvents(
	ctx context.Context,
	applicationID *uuid.UUID,
	sessionID *uuid.UUID,
	pagination *Pagination,
	options ...DBOption,
) (List[entities.Event], error) {
	db := GetDatabase(ctx)
	_, db = ApplyDBOptions(ctx, db, options...)

	var list List[entities.Event]
	if applicationID != nil {
		db = db.Where("applicationID=?", *applicationID)
	}
	if sessionID != nil {
		db = db.Where("sessionID=?", *sessionID)
	}
	err := list.ApplyPagination(db, pagination).
		Find(&list.Items).
		Error
	return list, err
}
func GetEventByID(ctx context.Context, eventID uuid.UUID, options ...DBOption) (*entities.Event, error) {
	db := GetDatabase(ctx)
	_, db = ApplyDBOptions(ctx, db, options...)

	var event *entities.Event
	if err := db.First(&event, eventID).Error; err != nil {
		return nil, err
	} else if event == nil || utils.IsZeroUUID(event.ID) {
		return nil, nil
	} else {
		return event, nil
	}
}
func CreateEvent(ctx context.Context, event *entities.Event, options ...DBOption) (*entities.Event, error) {
	db := GetDatabase(ctx)
	_, db = ApplyDBOptions(ctx, db, options...)

	if err := db.Create(event).Error; err != nil {
		return nil, err
	}

	return event, nil
}
