package srv

import "errors"

var (
	// ErrPortsRequired is returned when a healthcheck port has not been provided
	ErrPortsRequired = errors.New("no ports provided")

	errInvalidObjectNameLength = errors.New("object name must be less than 53 characters")
	errSubscriberCreate        = errors.New("unable to create subscriber")
	errSubscriptionCreate      = errors.New("unable to subscribe to topic")
	errInvalidHelmClient       = errors.New("unable to create helm client")
)
