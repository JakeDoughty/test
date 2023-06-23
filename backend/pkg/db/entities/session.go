package entities

import (
	"errors"
	"time"

	"github.com/JakeDoughty/customer-io-homework-backend/pkg/db/types"
	"github.com/JakeDoughty/customer-io-homework-backend/pkg/db/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	Model

	ApplicationID uuid.UUID        `gorm:"type:uuid;column:applicationID;not null"`
	CloseTime     time.Time        `gorm:"column:closedAt"`
	IP            types.NullString `gorm:"type:varchar(39);null;column:ip"`
	Browser       types.NullString `gorm:"column:browser;null"`
	OS            types.NullString `gorm:"column:os;null"`
	Screen        types.Size       `gorm:"embedded;embeddedPrefix=screen_"`

	Application *Application
	Events      []*Event `gorm:"foreignKey:SessionID"`
}

func (session *Session) TableName() string { return "Sessions" }
func (session *Session) GetScreen() *types.Size {
	if session.Screen.IsZero() {
		return nil
	} else {
		return &session.Screen
	}
}
func (session *Session) IsClosedAt(t time.Time) bool {
	return session.CloseTime.After(t)
}
func (session *Session) IsClosed() bool {
	return session.IsClosedAt(time.Now().UTC())
}

func (session *Session) BeforeCreate(tx *gorm.DB) error {
	session.Model.BeforeCreate(tx)
	return session.BeforeUpdate(tx)
}
func (session *Session) BeforeUpdate(tx *gorm.DB) error {
	if utils.IsZeroUUID(session.ApplicationID) {
		if session.Application != nil {
			session.ApplicationID = session.Application.ID
		}
		if utils.IsZeroUUID(session.ApplicationID) {
			return errors.New("session is missing application")
		}
	}
	for _, event := range session.Events {
		event.Session = session
		event.SessionID = session.ID
		event.Application = session.Application
		event.ApplicationID = session.ApplicationID
	}
	return nil
}
