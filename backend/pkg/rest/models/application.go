package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/JakeDoughty/customer-io-homework-backend/pkg/db/dbal"
	"github.com/JakeDoughty/customer-io-homework-backend/pkg/db/entities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Application struct {
	ID       uuid.UUID `json:"id,omitempty"`
	CreateAt time.Time `json:"createdAt"`
	Name     string    `json:"name" binding:"required"`
}

func (application *Application) ToEntity() *entities.Application {
	return &entities.Application{
		Model: entities.Model{
			ID: application.ID,
		},
		Name: application.Name,
	}
}

func ApplicationFromEntity(entity *entities.Application) *Application {
	return &Application{
		ID:       entity.ID,
		CreateAt: entity.CreatedAt,
		Name:     entity.Name,
	}
}

// GetAllApplications handle:
// GET /applications
func GetAllApplications(c *gin.Context) {
	pagination, err := readPaginationFromQuery(c)
	if err != nil {
		ErrorToResponse(err).Send(c, http.StatusBadRequest)
		return
	}

	if applications, err := dbal.GetAllApplications(c.Request.Context(), pagination); err != nil {
		ErrorToResponse(err).Send(c, http.StatusInternalServerError)
	} else {
		ConvertToListResponse(&applications, ApplicationFromEntity).Send(c, http.StatusOK)
	}
}

// GetApplicationById handle:
// GET /applications/:appId
func GetApplicationById(c *gin.Context) {
	if appId, err := readAppId(c); err != nil {
		ErrorToResponse(err).Send(c, http.StatusNotFound)
	} else if application, err := dbal.GetApplicationById(c.Request.Context(), *appId); err != nil {
		ErrorToResponse(err).Send(c, http.StatusInternalServerError)
	} else if application == nil {
		NewErrorResponsef("can't find an application with ID=%q", (*appId).String()).
			Send(c, http.StatusNotFound)
	} else {
		NewSuccessResponse(ApplicationFromEntity(application)).Send(c, http.StatusOK)
	}
}

// CreateApplication handles:
// POST /applications
func CreateApplication(c *gin.Context) {
	var application Application
	if err := c.ShouldBindJSON(&application); err != nil {
		ErrorToResponse(err).Send(c, http.StatusBadRequest)
		return
	}

	if app, err := dbal.GetApplicationByName(c.Request.Context(), application.Name); err != nil {
		ErrorToResponse(err).Send(c, http.StatusInternalServerError)
	} else if app != nil {
		NewErrorResponsef("another application with same name already exists with ID=%q", app.ID.String()).
			Send(c, http.StatusConflict)
	} else if app, err = dbal.CreateApplication(c.Request.Context(), application.ToEntity()); err != nil {
		ErrorToResponse(err).Send(c, http.StatusInternalServerError)
	} else {
		NewSuccessResponse(ApplicationFromEntity(app)).Send(c, http.StatusCreated)
	}
}

func InstallApplicationsApi_v1(apiRouter *gin.RouterGroup) {
	apiRouter.GET("/applications", GetAllApplications)
	apiRouter.POST("/applications", CreateApplication)
	apiRouter.GET(fmt.Sprintf("/applications/:%s", appId_Param), GetApplicationById)
}
