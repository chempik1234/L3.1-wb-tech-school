package transport

import (
	"context"
	"delayed_notifier/internal/dto"
	errors2 "delayed_notifier/internal/errors"
	"delayed_notifier/internal/models"
	"delayed_notifier/internal/service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NotifyHandler struct {
	crudService *service.NotificationCRUDService
}

func NewNotifyHandler(crudService *service.NotificationCRUDService) *NotifyHandler {
	return &NotifyHandler{crudService: crudService}
}

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
		if errors.Is(err, errors2.ErrNotificationNotFound) {
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
		if errors.Is(err, errors2.ErrNotificationNotFound) {
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
