package entities

import (
	"time"

	"github.com/JakeDoughty/customer-io-homework-backend/pkg/db/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Model struct {
	ID        uuid.UUID `gorm:"type:uuid;column:id;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (model *Model) BeforeCreate(tx *gorm.DB) error {
	if utils.IsZeroUUID(model.ID) {
		model.ID = utils.GenerateNewSequentialUUID()
	}
	return nil
}
