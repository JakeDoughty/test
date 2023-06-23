package models

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/JakeDoughty/customer-io-homework-backend/pkg/db/dbal"
	"github.com/JakeDoughty/customer-io-homework-backend/pkg/db/entities"
	"github.com/JakeDoughty/customer-io-homework-backend/pkg/db/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	appId_Param   = "appId"
	sessId_Param  = "sessId"
	eventId_Param = "eventId"
)

func readId(c *gin.Context, key string) (*uuid.UUID, error) {
	s := c.Param(key)
	if s == "" {
		return nil, nil
	} else if id, err := uuid.Parse(s); err != nil {
		return nil, err
	} else {
		return &id, nil
	}
}
func readAppId(c *gin.Context) (*uuid.UUID, error)     { return readId(c, appId_Param) }
func readSessionId(c *gin.Context) (*uuid.UUID, error) { return readId(c, sessId_Param) }
func readEventId(c *gin.Context) (*uuid.UUID, error)   { return readId(c, eventId_Param) }

func getSession(c *gin.Context, options ...dbal.DBOption) (session *entities.Session, handled bool) {
	if sessId, err := readSessionId(c); err != nil {
		ErrorToResponse(err).Send(c, http.StatusNotFound)
		return nil, true
	} else if sessId == nil {
		NewErrorResponse("missing session ID").Send(c, http.StatusNotFound)
		return nil, true
	} else if session, err := dbal.GetSessionByID(c.Request.Context(), *sessId, options...); err != nil {
		ErrorToResponse(err).Send(c, http.StatusNotFound)
		return nil, true
	} else if session == nil {
		NewErrorResponsef("can't find a session with ID==%q", (*sessId).String())
		return nil, true
	} else if appId, err := readAppId(c); err != nil {
		ErrorToResponse(err).Send(c, http.StatusNotFound)
		return nil, true
	} else if appId != nil && !utils.EqualUUID(*sessId, *appId) {
		NewErrorResponsef("can't find a session with ID==%q and ApplicationID==%q",
			session.ID.String(), (*appId).String()).
			Send(c, http.StatusNotFound)
		return nil, true
	} else {
		return session, false
	}
}
func getEvent(c *gin.Context, options ...dbal.DBOption) (event *entities.Event, handled bool) {
	if eventId, err := readEventId(c); err != nil {
		ErrorToResponse(err).Send(c, http.StatusNotFound)
		return nil, true
	} else if eventId == nil {
		NewErrorResponse("missing event ID").Send(c, http.StatusNotFound)
		return nil, true
	} else if event, err := dbal.GetEventByID(c.Request.Context(), *eventId, options...); err != nil {
		ErrorToResponse(err).Send(c, http.StatusInternalServerError)
		return nil, true
	} else if event == nil {
		NewErrorResponsef("can't find a session with ID==%q", (*eventId).String())
		return nil, true
	} else if appId, err := readAppId(c); err != nil {
		ErrorToResponse(err).Send(c, http.StatusNotFound)
		return nil, true
	} else if appId != nil && !utils.EqualUUID(*appId, event.ApplicationID) {
		NewErrorResponsef("can't find an event with ID==%q and ApplicationID==%q",
			event.ID.String(), (*appId).String()).
			Send(c, http.StatusNotFound)
		return nil, true
	} else if sessId, err := readSessionId(c); err != nil {
		ErrorToResponse(err).Send(c, http.StatusNotFound)
		return nil, true
	} else if sessId != nil && !utils.EqualUUID(*sessId, event.SessionID) {
		NewErrorResponsef("can't find an event with ID==%q and SessionID==%q",
			event.ID.String(), (*sessId).String()).
			Send(c, http.StatusNotFound)
		return nil, true
	} else {
		return event, false
	}
}

func readIpFromRemoteAddress(remoteAddr string) string {
	if remoteAddr == "" {
		return ""
	} else if n := strings.IndexByte(remoteAddr, ':'); n != -1 {
		return remoteAddr[:n]
	} else {
		return remoteAddr
	}
}

func readPaginationFromQuery(c *gin.Context) (*dbal.Pagination, error) {
	var pagination dbal.Pagination
	pageSizeStr := c.Query("pageSize")
	pageIndexStr := c.Query("pageIndex")
	if pageSizeStr != "" {
		if n, err := strconv.ParseUint(pageSizeStr, 10, 32); err != nil {
			return nil, fmt.Errorf("%s is not a valid pageSize", pageSizeStr)
		} else {
			pagination.PageSize = int(n)
		}
	}
	if pageIndexStr != "" {
		if n, err := strconv.ParseUint(pageIndexStr, 10, 32); err != nil {
			return nil, fmt.Errorf("%s is not a valid pageIndex", pageIndexStr)
		} else {
			pagination.PageIndex = int(n)
		}
	}

	pagination.Sort = c.Query("sort")
	return &pagination, nil
}
