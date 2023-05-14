package logger

import (
	"fmt"
	"log"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type ILogger interface {
	//Infow(string, ...interface{})
	//Errorw(string, ...interface{})
	//Warnw(string, ...interface{})
	//Debugw(string, ...interface{})
	//Info(...interface{})
	//Error(...interface{})
	//Warn(...interface{})
	//Debug(...interface{})
	//Infof(string, ...interface{})
	//Errorf(string, ...interface{})
	//Warnf(string, ...interface{})
	//Debugf(string, ...interface{})
}

// InitLogger initializes a new zap logger with default
// JSON formatter and call stack property
func Init(appMode string) *zap.SugaredLogger {
	atom := zap.NewAtomicLevel()
	// if error is nil new log level will be set.
	// In case err is NOT nil default atom level will be used (INFO)
	newLevel, err := getAtomLogLevel(appMode)
	if err == nil {
		atom.SetLevel(newLevel)
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "timestamp",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "message",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.UTC().Format(time.RFC3339Nano))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	config := zap.Config{
		Level:       atom,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	core, err := config.Build()
	if err != nil {
		log.Fatalf("building logger config:%v", err)
	}

	logger := zap.New(core.Core())

	defer logger.Sync()

	sugar := logger.Sugar()

	// setting default fields
	sugar = sugar.With(
		"microservice", "url_shortener",
	)

	sugar.Debug("logger initialized")
	return sugar
}

func getAtomLogLevel(appMode string) (zapcore.Level, error) {
	appMode = strings.ToUpper(appMode)
	var newLevel zapcore.Level

	switch appMode {
	case "ERROR":
		newLevel = zap.ErrorLevel
	case "WARN":
		newLevel = zap.WarnLevel
	case "INFO":
		newLevel = zap.InfoLevel
	case "DEBUG":
		newLevel = zap.DebugLevel
	default:
		return newLevel, fmt.Errorf("invalid log level %s", appMode)
	}
	return newLevel, nil
}
