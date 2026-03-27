package logger

import (
	"context"
	"log/slog"

	"github.com/visea-hive/auth-core/pkg/notifier"
)

// Logger is a unified logger that writes to slog and optionally sends
// a notification via the Notifier. Notification is only delivered when
// the embedded Async notifier wraps a real (non-noop) notifier — i.e.
// when notifications are enabled in config.
//
// Use the package-level functions (Info, Warn, Error, Critical) directly
// after calling SetDefault, or create an instance with New.
type Logger struct {
	notify *notifier.Async
}

// New creates a Logger backed by the given Async notifier.
// Pass the same *notifier.Async that is built in main.go.
func New(n *notifier.Async) *Logger {
	return &Logger{notify: n}
}

// --- Instance methods ---

// Info logs at INFO level and sends an info notification.
func (l *Logger) Info(msg string, args ...any) {
	slog.Info(msg, args...)
	if l.notify != nil {
		l.notify.Send(msg)
	}
}

// Warn logs at WARN level and sends a warning notification.
func (l *Logger) Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
	if l.notify != nil {
		l.notify.SendWithLevel(notifier.LevelWarning, msg)
	}
}

// Error logs at ERROR level and sends an error notification.
func (l *Logger) Error(msg string, args ...any) {
	slog.Error(msg, args...)
	if l.notify != nil {
		l.notify.SendWithLevel(notifier.LevelError, msg)
	}
}

// Critical logs at ERROR level (slog has no "critical") and sends a
// critical notification.
func (l *Logger) Critical(msg string, args ...any) {
	slog.Error(msg, args...)
	if l.notify != nil {
		l.notify.SendWithLevel(notifier.LevelCritical, msg)
	}
}

// Notify sends a titled notification without writing a slog entry.
// Useful when you want fine-grained control over the notification message.
func (l *Logger) Notify(title, msg string) {
	if l.notify != nil {
		l.notify.SendWithTitle(title, msg)
	}
}

// --- Package-level default ---

var defaultLogger *Logger

// SetDefault sets the package-level default Logger.
// Call this once during application startup (after initNotifier).
func SetDefault(l *Logger) {
	defaultLogger = l
}

// Default returns the package-level Logger. Panics if SetDefault has
// not been called — this is intentional to surface misconfiguration early.
func Default() *Logger {
	if defaultLogger == nil {
		panic("logger: SetDefault has not been called")
	}
	return defaultLogger
}

// Package-level helpers that mirror the Logger methods.

func Info(msg string, args ...any) {
	Default().Info(msg, args...)
}

func Warn(msg string, args ...any) {
	Default().Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Default().Error(msg, args...)
}

func Critical(msg string, args ...any) {
	Default().Critical(msg, args...)
}

// Notify sends a titled notification without a slog entry.
func Notify(title, msg string) {
	Default().Notify(title, msg)
}

// NotifyCtx is a context-aware variant; reserved for future use when the
// underlying Notifier gains context support on the Async layer.
func NotifyCtx(_ context.Context, title, msg string) {
	Default().Notify(title, msg)
}
