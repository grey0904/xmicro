package global

import (
    "go.uber.org/zap"
)

type Application struct {
	Log *zap.Logger
}

var App = new(Application)
