package logging

import (
	"github.com/Roma7-7-7/workshops/calendar/internal/config"
	"go.uber.org/zap"
)

var Logger, _ = zap.NewDevelopment()

func Init(env string) {
	var err error
	if env != config.DEV_ENV {
		Logger, err = zap.NewProduction()
	}

	if err != nil {
		panic(err)
	}
}
