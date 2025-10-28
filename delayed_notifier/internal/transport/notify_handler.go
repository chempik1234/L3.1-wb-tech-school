package transport

import (
	"context"
	"errors"
	"fmt"
	"github.com/chempik1234/L3.1-wb-tech-school/delayed_notifier/internal/dto"
	internalerrors "github.com/chempik1234/L3.1-wb-tech-school/delayed_notifier/internal/errors"
	"github.com/chempik1234/L3.1-wb-tech-school/delayed_notifier/internal/models"
	"github.com/chempik1234/L3.1-wb-tech-school/delayed_notifier/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// NotifyHandler is the HTTP routes handler, used in AssembleRouter
//
// Validates request and passes it to service layer
type NotifyHandler struct {
	crudService *service.NotificationCRUDService
}

// NewNotifyHandler creates a new NotifyHandler with given service
func NewNotifyHandler(crudService *service.NotificationCRUDService) *NotifyHandler {
	return &NotifyHandler{crudService: crudService}
}

// CreateNotification POST /notify
func (h *NotifyHandler) CreateNotification(c *gin.Context) {
	var body dto.CreateNotificationBody

	err := c.BindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid body (parsing): %s", err.Error())},
		)
		return
	}

	var createModel *models.Notification
	createModel, err = body.ToEntity()
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid body (validating): %s", err.Error())},
		)
		return
	}

	_, err = h.crudService.CreateNotification(context.Background(), createModel)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusConflict,
			gin.H{"error": fmt.Sprintf("couldn't perform operation: %s", err.Error())},
		)
		return
	}

	// createModel has been mutated: ID is now assigned!
	c.JSON(http.StatusCreated, dto.FullNotificationBodyFromEntity(createModel))
}

// GetNotification GET /notify/id
func (h *NotifyHandler) GetNotification(c *gin.Context) {
	req, err := dto.BindGetNotificationRequest(c)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"error": fmt.Sprintf("invalid ID parameter: %s", err.Error())},
		)
		return
	}

	id, err := req.ToUUID()
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"error": fmt.Sprintf("invalid UUID format: %s", err.Error())},
		)
		return
	}

	notification, err := h.crudService.GetNotification(context.Background(), id)
	if err != nil {
		if errors.Is(err, internalerrors.ErrNotificationNotFound) {
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				gin.H{"error": "notification not found"},
			)
			return
		}

		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("couldn't get notification: %s", err.Error())},
		)
		return
	}

	c.JSON(http.StatusOK, dto.FullNotificationBodyFromEntity(notification))
}

// DeleteNotification DELETE /notify/id
func (h *NotifyHandler) DeleteNotification(c *gin.Context) {
	req, err := dto.BindGetNotificationRequest(c)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"error": fmt.Sprintf("invalid ID parameter: %s", err.Error())},
		)
		return
	}

	id, err := req.ToUUID()
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"error": fmt.Sprintf("invalid UUID format: %s", err.Error())},
		)
		return
	}

	err = h.crudService.DeleteNotification(context.Background(), id)
	if err != nil {
		if errors.Is(err, internalerrors.ErrNotificationNotFound) {
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				gin.H{"error": "notification not found"},
			)
			return
		}

		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("couldn't delete notification: %s", err.Error())},
		)
		return
	}

	c.Status(http.StatusNoContent)
}
