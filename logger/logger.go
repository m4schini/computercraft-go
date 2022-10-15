package logger

import (
	"go.uber.org/zap"
	"log"
	"os"
	"strings"
)

var logger *zap.Logger

func init() {
	var err error

	lvl := os.Getenv("CC_GO_LOGGING_LEVEL")
	switch strings.ToLower(lvl) {
	case "dev":
		logger, err = zap.NewDevelopment()
		if err != nil {
			log.Fatalln(err)
		}
		break
	case "prod":
		logger, err = zap.NewProduction()
		if err != nil {
			log.Fatalln(err)
		}
		break
	default:
		logger = zap.NewNop()
		break
	}
}

func Sub(names ...string) *zap.Logger {
	subLogger := logger.Sugar()
	for _, name := range names {
		subLogger.Named(name)
	}

	return subLogger.Desugar()
}

func UseLogger(newLogger *zap.Logger) {
	logger = newLogger
}
