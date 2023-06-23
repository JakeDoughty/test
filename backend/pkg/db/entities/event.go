package entities

import (
	"errors"

	"github.com/JakeDoughty/customer-io-homework-backend/pkg/db/types"
	"github.com/JakeDoughty/customer-io-homework-backend/pkg/db/utils"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Event struct {
	Model

	ApplicationID uuid.UUID      `gorm:"type:uuid;column:applicationID;not null"`
	SessionID     uuid.UUID      `gorm:"type:uuid;column:sessionID;not null"`
	EventType     string         `gorm:"column:type;not null"`
	EventData     datatypes.JSON `gorm:"column:data;not null"`
	Location      types.Location `gorm:"embedded;embeddedPrefix:loc_"`

	Session     *Session
	Application *Application
}

func (event *Event) TableName() string { return "events" }

func (event *Event) BeforeCreate(tx *gorm.DB) error {
	event.Model.BeforeCreate(tx)
	return event.BeforeUpdate(tx)
}
func (event *Event) BeforeUpdate(tx *gorm.DB) error {
	if utils.IsZeroUUID(event.ApplicationID) {
		if event.Application != nil {
			event.ApplicationID = event.Application.ID
		}
		if utils.IsZeroUUID(event.ApplicationID) {
			return errors.New("event missing application ID")
		}
	}
	if utils.IsZeroUUID(event.SessionID) {
		if event.Session != nil {
			event.SessionID = event.Session.ID
		}
		if utils.IsZeroUUID(event.SessionID) {
			return errors.New("event missing session ID")
		}
	}
	return nil
}
func (event *Event) GetLocation() *types.Location {
	if event.Location.IsZero() {
		return nil
	} else {
		loc := event.Location
		return &loc
	}
}
