package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// Initialize initializes the global logger
func Initialize(level string) error {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapLevel)
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.FunctionKey = "function"
	config.EncoderConfig.StacktraceKey = "stacktrace"
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	logger, err := config.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return err
	}

	Log = logger
	return nil
}

// Sync flushes any buffered log entries
func Sync() error {
	if Log != nil {
		return Log.Sync()
	}
	return nil
}

// WithFields returns a logger with additional fields
func WithFields(fields ...zap.Field) *zap.Logger {
	if Log == nil {
		// Fallback to basic logger if not initialized
		logger, err := zap.NewProduction()
		if err != nil {
			// If we can't create a production logger, create a no-op logger
			return zap.NewNop().With(fields...)
		}
		return logger.With(fields...)
	}
	return Log.With(fields...)
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Debug(msg, fields...)
	}
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Info(msg, fields...)
	}
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Warn(msg, fields...)
	}
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Error(msg, fields...)
	} else {
		// Fallback to stderr if logger not initialized
		logger, err := zap.NewProduction()
		if err == nil {
			logger.Error(msg, fields...)
		}
	}
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Fatal(msg, fields...)
	} else {
		logger, err := zap.NewProduction()
		if err == nil {
			logger.Fatal(msg, fields...)
		}
	}
	os.Exit(1)
}

// Panic logs a panic message and panics
func Panic(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Panic(msg, fields...)
	}
}
