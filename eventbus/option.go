package eventbus

import (
	"github.com/sirupsen/logrus"
)

type Option func(*EventBus) error

func WithLogger(logger *logrus.Entry) Option {
	return func(e *EventBus) error {
		e.logger = logger
		return nil
	}
}

func WithName(name string) Option {
	return func(e *EventBus) error {
		e.name = name
		return nil
	}
}

func WithBucket(name string) Option {
	return func(e *EventBus) error {
		e.bucket = name
		return nil
	}
}
