package internaltypes

import (
	"fmt"
	"github.com/chempik1234/L3.1-wb-tech-school/delayed_notifier/pkg/types"
)

// ErrInvalidNotificationChannelValue describes an error when invalid string was put into NotificationChannel
var ErrInvalidNotificationChannelValue = fmt.Errorf("invalid notification channel value: possible ones are: '%s', '%s', '%s'", EMAIL, TELEGRAM, CONSOLE)

const (
	EMAIL    = "email"
	TELEGRAM = "telegram"
	CONSOLE  = "console"
)

// NotificationChannel is enum'd type for notification channels
//
// possible values: “email“, “telegram“, “console“
type NotificationChannel struct {
	val types.AnyText
}

// NotificationChannelFromString channel creates a new NotificationChannel object if it's valid
func NotificationChannelFromString(val string) (NotificationChannel, error) {
	switch val {
	case EMAIL, TELEGRAM, CONSOLE:
		break
	default:
		return NotificationChannel{}, ErrInvalidNotificationChannelValue
	}
	return NotificationChannel{val: types.NewAnyText(val)}, nil
}

// String method of NotificationChannel returns its string value
func (c *NotificationChannel) String() string {
	if c == nil {
		return "unset"
	}
	return c.val.String()
}
