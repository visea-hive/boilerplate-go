package config

import (
	"context"
	"log/slog"
	"time"

	"gorm.io/gorm/logger"
)

type SlogGormLogger struct {
	LogLevel logger.LogLevel
}

func NewSlogGormLogger(level logger.LogLevel) *SlogGormLogger {
	return &SlogGormLogger{LogLevel: level}
}

func (l *SlogGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &SlogGormLogger{LogLevel: level}
}

func (l *SlogGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		slog.LogAttrs(ctx, slog.LevelInfo, msg, slog.Any("data", data))
	}
}

func (l *SlogGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		slog.LogAttrs(ctx, slog.LevelWarn, msg, slog.Any("data", data))
	}
}

func (l *SlogGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		slog.LogAttrs(ctx, slog.LevelError, msg, slog.Any("data", data))
	}
}

func (l *SlogGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()

	attrs := []slog.Attr{
		slog.String("elapsed", elapsed.String()),
		slog.Int64("rows", rows),
		slog.String("sql", sql),
	}

	if err != nil && l.LogLevel >= logger.Error {
		attrs = append(attrs, slog.Any("error", err))
		slog.LogAttrs(ctx, slog.LevelError, "gorm trace", attrs...)
	} else if l.LogLevel >= logger.Info {
		slog.LogAttrs(ctx, slog.LevelInfo, "gorm trace", attrs...)
	}
}
