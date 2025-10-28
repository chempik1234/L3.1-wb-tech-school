package errors

import "errors"

// ErrNotificationNotFound occurs when searched notification couldn't be found
//
// Used by both service and repo
var ErrNotificationNotFound = errors.New("notification not found")
