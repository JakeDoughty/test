package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/JakeDoughty/customer-io-homework-backend/pkg/db/dbal"
	"github.com/JakeDoughty/customer-io-homework-backend/pkg/db/entities"
	"github.com/JakeDoughty/customer-io-homework-backend/pkg/db/types"
	"github.com/JakeDoughty/customer-io-homework-backend/pkg/db/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID   `json:"id"`
	CreatedAt time.Time   `json:"createdAt"`
	CloseTime time.Time   `json:"closeTime"`
	Closed    bool        `json:"closed,omitempty"`
	IP        string      `json:"ip,omitempty"`
	Browser   string      `json:"browser,omitempty"`
	Screen    *types.Size `json:"screen,omitempty"`
	Link      string      `json:"link,omitempty"`
}

func (session *Session) ToEntity(applicationID *uuid.UUID) *entities.Session {
	entity := &entities.Session{
		Model: entities.Model{
			ID:        session.ID,
			CreatedAt: session.CreatedAt,
		},
		CloseTime: session.CloseTime,
		IP:        types.StringToNullString(session.IP),
		Browser:   types.StringToNullString(session.Browser),
	}
	if applicationID != nil {
		entity.ApplicationID = *applicationID
	}
	if session.Screen != nil {
		entity.Screen = *session.Screen
	}
	return entity
}

func SessionFromEntity(entity *entities.Session) *Session {
	return &Session{
		ID:        entity.ID,
		CreatedAt: entity.CreatedAt,
		CloseTime: entity.CloseTime,
		Closed:    entity.IsClosed(),
		IP:        entity.IP.ToString(),
		Browser:   entity.Browser.ToString(),
		Screen:    entity.GetScreen(),
		Link:      fmt.Sprintf("/api/v1/sessions/%s", entity.ID.String()),
	}
}

// GetAllSessions handle:
// GET /api/v1/applications/:appId/sessions
func GetAllSessions(c *gin.Context) {
	pagination, err := readPaginationFromQuery(c)
	if err != nil {
		ErrorToResponse(err).Send(c, http.StatusBadRequest)
		return
	}

	options := []dbal.DBOption{}
	if c.Query("includeClosedSessions") == "true" {
		options = append(options, dbal.IncludeClosedSessions())
	}
	if appId, err := readAppId(c); err != nil {
		ErrorToResponse(err).Send(c, http.StatusNotFound)
	} else if sessions, err := dbal.GetAllSessions(
		c.Request.Context(), appId, pagination, options...,
	); err != nil {
		ErrorToResponse(err).Send(c, http.StatusInternalServerError)
	} else {
		ConvertToListResponse(&sessions, SessionFromEntity).Send(c, http.StatusOK)
	}
}

// GetSessionByID handles:
// GET /api/v1/sessions/:sessId
// GET /api/v1/applications/:appId/sessions/:sessId
func GetSessionByID(c *gin.Context) {
	if sessId, err := readSessionId(c); err != nil {
		ErrorToResponse(err).Send(c, http.StatusNotFound)
	} else if session, err := dbal.GetSessionByID(c.Request.Context(), *sessId); err != nil {
		ErrorToResponse(err).Send(c, http.StatusInternalServerError)
	} else if session == nil {
		NewErrorResponsef("can't find a session with ID==%q", (*sessId).String()).
			Send(c, http.StatusNotFound)
	} else if appId, err := readAppId(c); err != nil {
		ErrorToResponse(err).Send(c, http.StatusNotFound)
	} else if appId != nil && !utils.EqualUUID(*appId, session.ApplicationID) {
		NewErrorResponsef("can't find a session with ID==%q and ApplicationID==%q",
			(*sessId).String(), (*appId).String()).
			Send(c, http.StatusNotFound)
	} else {
		NewSuccessResponse(SessionFromEntity(session)).Send(c, http.StatusOK)
	}
}

// CreateSession handle:
// POST /api/v1/applications/:appId/sessions
func CreateSession(c *gin.Context) {
	var session *Session
	if err := c.ShouldBindJSON(&session); err != nil {
		ErrorToResponse(err).Send(c, http.StatusBadRequest)
		return
	} else {
		session.IP = readIpFromRemoteAddress(c.Request.RemoteAddr)
		session.Browser = c.Request.Header.Get("browser")
		session.CreatedAt = time.Now().UTC()
		session.CloseTime = session.CreatedAt.Add(15 * time.Minute)
	}

	appId, err := readAppId(c)
	if err != nil {
		ErrorToResponse(err).Send(c, http.StatusNotFound)
		return
	}

	if application, err := dbal.GetApplicationById(c.Request.Context(), *appId); err != nil {
		ErrorToResponse(err).Send(c, http.StatusInternalServerError)
	} else if application == nil {
		ErrorToResponse(err).Send(c, http.StatusNotFound)
	} else if session, err := dbal.CreateSession(c.Request.Context(), session.ToEntity(appId)); err != nil {
		ErrorToResponse(err).Send(c, http.StatusInternalServerError)
	} else {
		NewSuccessResponse(SessionFromEntity(session)).Send(c, http.StatusCreated)
	}
}

func CloseSession(c *gin.Context) {
	if sessId, err := readSessionId(c); err != nil {
		ErrorToResponse(err).Send(c, http.StatusNotFound)
	} else if sessId == nil {
		NewErrorResponse("missing session ID").Send(c, http.StatusNotFound)
	} else if session, err := dbal.GetSessionByID(c.Request.Context(), *sessId); err != nil {
		ErrorToResponse(err).Send(c, http.StatusInternalServerError)
	} else if session == nil {
		NewErrorResponsef("can't find a session with ID=%q", (*sessId).String()).
			Send(c, http.StatusNotFound)
	} else {
		session.CloseTime = time.Now()
		if session, err = dbal.UpdateSession(c.Request.Context(), session); err != nil {
			ErrorToResponse(err).Send(c, http.StatusInternalServerError)
		} else {
			NewSuccessResponse(SessionFromEntity(session)).Send(c, http.StatusOK)
		}
	}
}

func InstallSessionsApi_v1(apiRouter *gin.RouterGroup) {
	apiRouter.GET(fmt.Sprintf("/applications/:%s/sessions", appId_Param), GetAllSessions)
	apiRouter.POST(fmt.Sprintf("/applications/:%s/sessions", appId_Param), CreateSession)
	apiRouter.GET(fmt.Sprintf("/sessions/:%s", sessId_Param), GetSessionByID)
	apiRouter.GET(fmt.Sprintf("/applications/:%s/sessions/:%s", appId_Param, sessId_Param), GetSessionByID)
	apiRouter.DELETE(fmt.Sprintf("/applications/:%s/sessions/:%s", appId_Param, sessId_Param), CloseSession)
}
