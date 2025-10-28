package internaltypes

import (
	"fmt"
	"github.com/chempik1234/L3.1-wb-tech-school/delayed_notifier/pkg/types"
	"net/mail"
	"strconv"
)

// ErrInvalidNotificationChannelValue describes an error when invalid string was put into NotificationChannel
var ErrInvalidNotificationChannelValue = fmt.Errorf("invalid notification channel value: possible ones are: '%s', '%s', '%s'", EMAIL, TELEGRAM, CONSOLE)

const (
	// EMAIL is the constant value for email channel string value
	EMAIL = "email"
	// TELEGRAM is the constant value for telegram channel string value
	TELEGRAM = "telegram"
	// CONSOLE is the constant value for console channel string value
	CONSOLE = "console"
)

var (
	// ChannelEmail is an example channel with value EMAIL
	ChannelEmail = NotificationChannel{val: EMAIL}
	// ChannelTelegram is an example channel with value TELEGRAM
	ChannelTelegram = NotificationChannel{val: TELEGRAM}
	// ChannelConsole is an example channel with value CONSOLE
	ChannelConsole = NotificationChannel{val: CONSOLE}
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

// SendTo is a value type that stores an address of user according to notification channel
//
//	email	-> some@email.com
//	telegram	-> 13123123129 user id
//	console	-> any
type SendTo struct {
	val types.AnyText
}

// String returns value stored in SendTo
func (s SendTo) String() string {
	return s.val.String()
}

// NewSendTo creates a new SendTo with a valid address for given channel
func NewSendTo(val types.AnyText, channel NotificationChannel) (SendTo, error) {
	switch channel {
	case ChannelEmail:
		_, err := mail.ParseAddress(channel.val.String())
		if err != nil {
			return SendTo{}, fmt.Errorf("invalid email address: %w", err)
		}
	case ChannelTelegram:
		_, err := strconv.ParseInt(val.String(), 10, 64)
		if err != nil {
			return SendTo{}, fmt.Errorf("invalid telegram address: %s", val.String())
		}
	default:
		break
	}
	return SendTo{val: val}, nil
}
