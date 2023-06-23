package dbal

import (
	"context"

	"github.com/JakeDoughty/customer-io-homework-backend/pkg/db/entities"
	"github.com/JakeDoughty/customer-io-homework-backend/pkg/db/utils"
	"github.com/google/uuid"
)

func GetAllApplications(
	ctx context.Context,
	pageination *Pagination,
	options ...DBOption,
) (List[entities.Application], error) {
	db := GetDatabase(ctx)
	_, db = ApplyDBOptions(ctx, db, options...)

	var list List[entities.Application]
	err := list.ApplyPagination(db, pageination).Find(&list.Items).Error
	return list, err
}
func GetApplicationById(ctx context.Context, applicationID uuid.UUID, options ...DBOption) (*entities.Application, error) {
	db := GetDatabase(ctx)
	_, db = ApplyDBOptions(ctx, db, options...)

	var application *entities.Application
	if err := db.First(&application, applicationID).Error; err != nil {
		return nil, err
	} else if application == nil || utils.IsZeroUUID(application.ID) {
		return nil, nil
	} else {
		return application, nil
	}
}
func GetApplicationByName(ctx context.Context, name string, options ...DBOption) (*entities.Application, error) {
	db := GetDatabase(ctx)
	_, db = ApplyDBOptions(ctx, db, options...)

	var application *entities.Application
	if err := db.Where("name=?", name).First(&application).Error; err != nil {
		return nil, err
	} else if application == nil || utils.IsZeroUUID(application.ID) {
		return nil, nil
	} else {
		return application, nil
	}
}
func CreateApplication(ctx context.Context, application *entities.Application, options ...DBOption) (*entities.Application, error) {
	db := GetDatabase(ctx)
	_, db = ApplyDBOptions(ctx, db, options...)

	if err := db.Create(application).Error; err != nil {
		return nil, err
	}
	return application, nil
}
