// get_notification.go
package dto

import (
	"delayed_notifier/pkg/types"
	"github.com/gin-gonic/gin"
)

// GetNotificationRequest is a DTO for get/delete endpoint path parameters
type GetNotificationRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// ToUUID converts string ID to types.UUID
func (r *GetNotificationRequest) ToUUID() (types.UUID, error) {
	return types.NewUUID(r.ID)
}

// BindGetNotificationRequest binds and validates get/delete request
func BindGetNotificationRequest(c *gin.Context) (*GetNotificationRequest, error) {
	var req GetNotificationRequest
	if err := c.ShouldBindUri(&req); err != nil {
		return nil, err
	}
	return &req, nil
}
