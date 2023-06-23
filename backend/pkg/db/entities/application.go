package entities

import "gorm.io/gorm"

type Application struct {
	Model

	Name     string     `gorm:"uniqueIndex;not null"`
	Sessions []*Session `gorm:"foreignKey:ApplicationID"`
	Events   []*Event   `gorm:"foreignKey:ApplicationID"`
}

func (app *Application) TableName() string { return "Applications" }
func (app *Application) BeforeCreate(tx *gorm.DB) error {
	app.Model.BeforeCreate(tx)
	return app.BeforeUpdate(tx)
}
func (app *Application) BeforeUpdate(tx *gorm.DB) error {
	// update application ID of all the session and events
	for _, session := range app.Sessions {
		session.Application = app
		session.ApplicationID = app.ID
	}
	for _, event := range app.Events {
		event.Application = app
		event.ApplicationID = app.ID
	}
	return nil
}
