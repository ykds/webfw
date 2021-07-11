package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(filename string, maxSize, maxBackups, maxAge int, compress bool) *zap.Logger {
	core := zapcore.NewCore(getEncoder(), getWriteSyncer(filename, maxSize, maxBackups, maxAge, compress), zapcore.InfoLevel)
	return zap.New(core, zap.AddCaller(),  zap.AddStacktrace(zapcore.ErrorLevel))
}

func getEncoder() zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	config.TimeKey = "datetime"
	return zapcore.NewJSONEncoder(config)
}

func getWriteSyncer(filename string, maxSize, maxBackups, maxAge int, compress bool) zapcore.WriteSyncer {
	logger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
	}
	return zapcore.AddSync(logger)
}
