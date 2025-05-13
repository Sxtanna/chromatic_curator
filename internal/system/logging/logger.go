package logging

import (
	"emperror.dev/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	zapadapter "logur.dev/adapter/zap"
	"logur.dev/logur"
)

const (
	encodingJson    = "json"
	encodingConsole = "console"
)

func New(config Config) (logur.LoggerFacade, error) {
	var zc zap.Config

	if !config.Dev {
		zc = zap.NewProductionConfig()
		zc.EncoderConfig.EncodeTime = zapcore.EpochMillisTimeEncoder
	} else {
		zc = zap.NewDevelopmentConfig()
		zc.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	if config.Encoding == encodingJson {
		zc.Encoding = encodingJson
	} else if config.Encoding == encodingConsole {
		zc.Encoding = encodingConsole
	}

	zc.EncoderConfig.TimeKey = "timestamp"
	zc.EncoderConfig.EncodeName = zapcore.FullNameEncoder

	zc.OutputPaths = config.Output

	if lvl, err := zap.ParseAtomicLevel(config.Level); err == nil {
		zc.Level = lvl
	}

	if l, err := zc.Build(zap.AddCallerSkip(1)); err != nil {
		return nil, errors.Wrap(err, "failed to build logger")
	} else {
		return zapadapter.New(l), nil
	}
}

func SetStandardLogger(logger logur.Logger) {
	log.SetOutput(logur.NewLevelWriter(logger, logur.Info))
}
