package rest

import (
	"github.com/JakeDoughty/customer-io-homework-backend/pkg/rest/models"
	"github.com/gin-gonic/gin"
)

func RestServer_v1(apiRouter_v1 *gin.RouterGroup) {
	models.InstallApplicationsApi_v1(apiRouter_v1)
	models.InstallSessionsApi_v1(apiRouter_v1)
	models.InstallEventsApi_v1(apiRouter_v1)
}
