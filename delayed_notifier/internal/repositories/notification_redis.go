package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chempik1234/L3.1-wb-tech-school/delayed_notifier/internal/models"
	"github.com/chempik1234/L3.1-wb-tech-school/delayed_notifier/pkg/types"
	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/retry"
	"time"
)

// NotificationRedis implements ports.NotificationCRUDCacheRepository
//
// json.Marshal for serialization
//
// KEY is set like "notification:<uuid>"
type NotificationRedis struct {
	redisClient *redis.Client
	expiration  time.Duration
	strategy    retry.Strategy
}

// NewNotificationRedis creates a new NotificationRedis
func NewNotificationRedis(redisClient *redis.Client, retryStrategy retry.Strategy, expiration time.Duration) *NotificationRedis {
	return &NotificationRedis{redisClient: redisClient, strategy: retryStrategy, expiration: expiration}
}

// SaveNotification is both the Create and Update methods of this Cache CRUD
func (r *NotificationRedis) SaveNotification(ctx context.Context, notification *models.Notification) error {
	key := r.key(*notification.ID)
	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}
	return r.redisClient.SetWithExpiration(ctx, key, data, r.expiration)
}

// GetNotification is the Read method of this Cache CRUD
//
// err on redis.NoMatches
func (r *NotificationRedis) GetNotification(ctx context.Context, id types.UUID) (*models.Notification, error) {
	key := r.key(id)

	data, err := r.redisClient.Get(ctx, key)
	if err != nil {
		if errors.Is(err, redis.NoMatches) {
			return nil, errors.New("notification not found")
		}
		return nil, err
	}
	var notification models.Notification

	if err = json.Unmarshal([]byte(data), &notification); err != nil {
		return nil, err
	}
	return &notification, nil
}

// DeleteNotification is the Delete method of this Cache CRUD
func (r *NotificationRedis) DeleteNotification(ctx context.Context, id types.UUID) error {
	key := r.key(id)
	err := r.redisClient.Del(ctx, key)
	if err != nil {
		return fmt.Errorf("error deleting from redis notification (id '%s'): %w", id, err)
	}
	return nil
}

func (r *NotificationRedis) key(id types.UUID) string {
	return fmt.Sprintf("notification:%s", id)
}
