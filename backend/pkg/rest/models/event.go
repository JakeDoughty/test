package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/JakeDoughty/customer-io-homework-backend/pkg/db/dbal"
	"github.com/JakeDoughty/customer-io-homework-backend/pkg/db/entities"
	dbtypes "github.com/JakeDoughty/customer-io-homework-backend/pkg/db/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Event struct {
	ID            uuid.UUID         `json:"id,omitempty"`
	CreatedAt     time.Time         `json:"createAt"`
	ApplicationID uuid.UUID         `json:"appID,omitempty"`
	SessionID     uuid.UUID         `json:"sessionID,omitempty"`
	Type          string            `json:"type" binding:"required"`
	Data          json.RawMessage   `json:"data" binding:"required"`
	Location      *dbtypes.Location `json:"location,omitempty"`
	Link          string            `json:"link,omitempty"`
}

func (event *Event) ToEntity(appId, sessId uuid.UUID) *entities.Event {
	entity := &entities.Event{
		Model: entities.Model{
			ID:        event.ID,
			CreatedAt: event.CreatedAt,
		},
		ApplicationID: appId,
		SessionID:     sessId,
		EventType:     event.Type,
		EventData:     datatypes.JSON(event.Data),
	}
	if event.Location != nil {
		entity.Location = *event.Location
	}
	return entity
}

func EventFromEntity(event *entities.Event) *Event {
	return &Event{
		ID:            event.ID,
		CreatedAt:     event.CreatedAt,
		ApplicationID: event.ApplicationID,
		SessionID:     event.SessionID,
		Type:          event.EventType,
		Data:          json.RawMessage(event.EventData),
		Location:      event.GetLocation(),
		Link:          fmt.Sprintf("/api/v1/sessions/%s/events/%s", event.SessionID, event.ID),
	}
}

// GetAllEvents handle:
// GET /api/v1/applications/:appId/sessions/:sessId/events
// GET /api/v1/'applications/:appId/events
// GET /api/v1/sessions/:sessId/events
func GetAllEvents(c *gin.Context) {
	pagination, err := readPaginationFromQuery(c)
	if err != nil {
		ErrorToResponse(err).Send(c, http.StatusBadRequest)
		return
	}

	if appId, err := readAppId(c); err != nil {
		ErrorToResponse(err).Send(c, http.StatusBadRequest)
	} else if sessId, err := readSessionId(c); err != nil {
		ErrorToResponse(err).Send(c, http.StatusBadRequest)
	} else if events, err := dbal.GetAllEvents(
		c.Request.Context(),
		appId,
		sessId,
		pagination,
	); err != nil {
		ErrorToResponse(err).Send(c, http.StatusInternalServerError)
	} else {
		ConvertToListResponse(&events, EventFromEntity).Send(c, http.StatusOK)
	}
}

// GetEventByID handle:
// GET /application/:appId/sessions/:sessId/events/:eventId
// GET /application/:appId/events/:eventId
// GET /sessions/:sessId/events/:eventId
// GET /events/:eventId
func GetEventByID(c *gin.Context) {
	if event, handled := getEvent(c); !handled {
		NewSuccessResponse(EventFromEntity(event)).Send(c, http.StatusOK)
	}
}

// CreateEvent handle
// POST /
func CreateEvent(c *gin.Context) {
	var event Event
	if err := c.ShouldBindJSON(&event); err != nil {
		ErrorToResponse(err).Send(c, http.StatusBadRequest)
		return
	}

	if session, handled := getSession(c); !handled {
		if event.Type != "ping" {
			if entity, err := dbal.CreateEvent(c.Request.Context(), event.ToEntity(session.ApplicationID, session.ID)); err != nil {
				ErrorToResponse(err).Send(c, http.StatusInternalServerError)
			} else {
				// update closetime of the session
				session.CloseTime = time.Now().UTC().Add(15 * time.Minute)
				dbal.UpdateSession(c.Request.Context(), session)
				NewSuccessResponse(EventFromEntity(entity)).Send(c, http.StatusCreated)
			}
		} else {
			// just update close time of the session
			session.CloseTime = time.Now().UTC().Add(15 * time.Minute)
			if _, err := dbal.UpdateSession(c.Request.Context(), session); err != nil {
				ErrorToResponse(err).Send(c, http.StatusInternalServerError)
			} else {
				NewSuccessResponse(true).Send(c, http.StatusOK)
			}
		}
	}
}

func InstallEventsApi_v1(apiRouter *gin.RouterGroup) {
	apiRouter.GET(fmt.Sprintf("/applications/:%s/events", appId_Param), GetAllEvents)
	apiRouter.GET(fmt.Sprintf("/sessions/:%s/events", sessId_Param), GetAllEvents)
	apiRouter.POST(fmt.Sprintf("/sessions/:%s/events", sessId_Param), CreateEvent)
	apiRouter.GET(fmt.Sprintf("/sessions/:%s/events/:%s", sessId_Param, eventId_Param), GetEventByID)
	apiRouter.GET(fmt.Sprintf("/applications/:%s/events/:%s", appId_Param, eventId_Param), GetEventByID)
	apiRouter.GET(fmt.Sprintf("/events/:%s", eventId_Param), GetEventByID)
}
