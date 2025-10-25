package service

import (
	"context"
	"fmt"
	"github.com/chempik1234/L3.1-wb-tech-school/delayed_notifier/internal/errors"
	"github.com/chempik1234/L3.1-wb-tech-school/delayed_notifier/internal/models"
	"github.com/chempik1234/L3.1-wb-tech-school/delayed_notifier/internal/ports"
	"github.com/chempik1234/L3.1-wb-tech-school/delayed_notifier/pkg/types"
	"github.com/wb-go/wbf/zlog"
	"golang.org/x/sync/errgroup"
)

// SignalFunc is a function that's
type SignalFunc func(ctx context.Context, notification *models.Notification) error

// NotificationCRUDService is the service for storing, caching, receiving, updating/deleting notifications
type NotificationCRUDService struct {
	storageRepo ports.NotificationCRUDStorageRepository
	cacheRepo   ports.NotificationCRUDCacheRepository

	// funcOnCreate is called after CreateNotification
	//
	// for example: SenderService.QuickSend
	funcOnCreate SignalFunc
}

// NewNotificationCRUDService creates a new NotificationCRUDService with given adapters (storage, cache)
func NewNotificationCRUDService(
	storageRepo ports.NotificationCRUDStorageRepository,
	cacheRepo ports.NotificationCRUDCacheRepository,
	funcOnCreate SignalFunc,
) *NotificationCRUDService {
	return &NotificationCRUDService{storageRepo: storageRepo, cacheRepo: cacheRepo, funcOnCreate: funcOnCreate}
}

// CreateNotification saves a new notification
//
// Event ID is generated in service layer, so:
//
// 1. This mutates the model
//
// 2. This returns the model back
func (s *NotificationCRUDService) CreateNotification(ctx context.Context, model *models.Notification) (*models.Notification, error) {
	id := types.GenerateUUID()
	model.ID = &id

	err := s.storageRepo.CreateNotification(ctx, model) // retry is called inside
	if err != nil {
		return nil, fmt.Errorf("notification storage failed to create: %v", err)
	}

	s.tryCacheNotificationInBackground(ctx, model)

	// call signal in bg
	if s.funcOnCreate != nil {
		go func(model *models.Notification) {
			funcErr := s.funcOnCreate(ctx, model)
			if funcErr != nil {
				zlog.Logger.Error().Err(funcErr).Msg(fmt.Sprintf("error in funcOnCreate %v", s.funcOnCreate))
			}
		}(model)
	}

	return model, nil
}

// GetNotification returns notification, firstly checking the cache, then the storage
func (s *NotificationCRUDService) GetNotification(ctx context.Context, id types.UUID) (*models.Notification, error) {
	result, err := s.getObjectFromCache(ctx, id) // retry is called inside
	if err != nil {
		result, err = s.getObjectFromStorage(ctx, id) // retry is called inside
		if err != nil {
			return nil, fmt.Errorf("error getting object from storage: %w", err)
		}

		if result != nil {
			s.tryCacheNotificationInBackground(ctx, result)
		}
	}

	return result, err
}

// DeleteNotification deletes notification if exists (invalidates cache, affects storage)
//
// returns error on NotFound or Internal Error
func (s *NotificationCRUDService) DeleteNotification(ctx context.Context, id types.UUID) error {
	object, err := s.GetNotification(ctx, id)
	if err != nil {
		return fmt.Errorf("error checking object existence: %w", err)
	}

	// just to ensure, lol!
	if object == nil || object.ID == nil {
		return errors.ErrNotificationNotFound
	}

	errGroup := &errgroup.Group{}

	errGroup.Go(func() error { return s.storageRepo.DeleteNotification(ctx, *object.ID) })
	errGroup.Go(func() error { return s.cacheRepo.DeleteNotification(ctx, *object.ID) })

	return errGroup.Wait()
}

// PRIVATE METHODS

func (s *NotificationCRUDService) getObjectFromStorage(ctx context.Context, id types.UUID) (*models.Notification, error) {
	return s.storageRepo.GetNotification(ctx, id)
}

func (s *NotificationCRUDService) getObjectFromCache(ctx context.Context, id types.UUID) (*models.Notification, error) {
	return s.cacheRepo.GetNotification(ctx, id)
}

// tryCacheNotificationInBackground launches cache SET in the background, logs on error
func (s *NotificationCRUDService) tryCacheNotificationInBackground(ctx context.Context, model *models.Notification) {
	go func() {
		cacheErr := s.cacheRepo.SaveNotification(ctx, model)
		if cacheErr != nil {
			zlog.Logger.Error().Err(fmt.Errorf("error saving in cache: %w", cacheErr))
		}
	}()
}
