package logging

import (
	"emperror.dev/errors"
	"github.com/Sxtanna/chromatic_curator/internal/common"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
	"log"
	"log/slog"
)

const (
	encodingJson    = "json"
	encodingConsole = "console"
)

const (
	configurationMissing = errors.Sentinel("logging configuration missing")
)

func New(config common.Configuration) (*slog.Logger, error) {
	loggingConfiguration := common.FindConfiguration[Config](config)
	if loggingConfiguration == nil {
		return nil, configurationMissing
	}

	var zc zap.Config

	if !loggingConfiguration.Dev {
		zc = zap.NewProductionConfig()
		zc.EncoderConfig.EncodeTime = zapcore.EpochMillisTimeEncoder
	} else {
		zc = zap.NewDevelopmentConfig()
		zc.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	if loggingConfiguration.Encoding == encodingJson {
		zc.Encoding = encodingJson
	} else if loggingConfiguration.Encoding == encodingConsole {
		zc.Encoding = encodingConsole
	}

	zc.EncoderConfig.TimeKey = "timestamp"
	zc.EncoderConfig.EncodeName = zapcore.FullNameEncoder

	zc.OutputPaths = loggingConfiguration.Output

	if lvl, err := zap.ParseAtomicLevel(loggingConfiguration.Level); err == nil {
		zc.Level = lvl
	}

	if l, err := zc.Build(zap.AddCallerSkip(1)); err != nil {
		return nil, errors.Wrap(err, "failed to build logger")
	} else {
		return slog.New(zapslog.NewHandler(l.Core())), nil
	}
}

func SetStandardLogger(logger *slog.Logger) {
	log.SetOutput(NewSlogWriter(logger))
}
