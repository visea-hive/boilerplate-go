package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// Level represents the severity of a notification.
type Level string

const (
	LevelInfo     Level = "info"     // ℹ️ Informational messages
	LevelWarning  Level = "warning"  // ⚠️ Warning messages
	LevelError    Level = "error"    // 🐛 Error/bug messages
	LevelCritical Level = "critical" // 🚨 Critical/emergency messages
)

// Emoji returns the emoji prefix for the level.
func (l Level) Emoji() string {
	switch l {
	case LevelWarning:
		return "⚠️"
	case LevelError:
		return "🐛"
	case LevelCritical:
		return "🚨"
	default:
		return "ℹ️"
	}
}

// Notifier is the interface for sending notifications.
// Any notification provider (Slack, Mattermost, Discord, etc.) must implement this.
type Notifier interface {
	// Send delivers a notification message (defaults to Info level).
	Send(ctx context.Context, message string) error

	// SendWithTitle delivers a notification with a title and body.
	SendWithTitle(ctx context.Context, title string, message string) error

	// SendWithLevel delivers a notification with a severity level.
	SendWithLevel(ctx context.Context, level Level, message string) error
}

// --- WebhookNotifier ---

// WebhookNotifier sends notifications via an incoming webhook URL.
// Compatible with Slack, Mattermost, and any service that accepts
// a JSON payload with a "text" field.
type WebhookNotifier struct {
	webhookURL string
	client     *http.Client
	provider   string
}

type webhookPayload struct {
	Text string `json:"text"`
}

// NewWebhookNotifier creates a new WebhookNotifier with a 5-second timeout.
func NewWebhookNotifier(webhookURL string, provider string) *WebhookNotifier {
	return &WebhookNotifier{
		webhookURL: webhookURL,
		provider:   provider,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (w *WebhookNotifier) Send(ctx context.Context, message string) error {
	return w.SendWithLevel(ctx, LevelInfo, message)
}

func (w *WebhookNotifier) SendWithTitle(ctx context.Context, title string, message string) error {
	formatted := fmt.Sprintf("*%s*\n%s", title, message)
	return w.send(ctx, formatted)
}

func (w *WebhookNotifier) SendWithLevel(ctx context.Context, level Level, message string) error {
	formatted := fmt.Sprintf("%s %s", level.Emoji(), message)
	return w.send(ctx, formatted)
}

func (w *WebhookNotifier) send(ctx context.Context, text string) error {
	body, err := json.Marshal(webhookPayload{Text: text})
	if err != nil {
		return fmt.Errorf("failed to marshal notification payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.webhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create notification request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("%s webhook request failed: %w", w.provider, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("%s webhook returned status %d", w.provider, resp.StatusCode)
	}

	slog.Debug("Notification sent", "provider", w.provider)
	return nil
}

// --- NoOpNotifier ---

// NoOpNotifier logs messages instead of sending them.
// Used in local/dev environments or when notifications are disabled.
type NoOpNotifier struct{}

func NewNoOpNotifier() *NoOpNotifier {
	return &NoOpNotifier{}
}

func (n *NoOpNotifier) Send(ctx context.Context, message string) error {
	slog.Debug("Notification skipped (noop)", "message", message)
	return nil
}

func (n *NoOpNotifier) SendWithTitle(ctx context.Context, title string, message string) error {
	slog.Debug("Notification skipped (noop)", "title", title, "message", message)
	return nil
}

func (n *NoOpNotifier) SendWithLevel(ctx context.Context, level Level, message string) error {
	slog.Debug("Notification skipped (noop)", "level", string(level), "message", message)
	return nil
}

// --- Async Wrapper ---

// Async wraps a Notifier to send notifications in a background goroutine.
// Errors are logged but never propagated — your app is never affected.
type Async struct {
	notifier Notifier
}

// NewAsync wraps the given Notifier for async delivery.
func NewAsync(n Notifier) *Async {
	return &Async{notifier: n}
}

func (a *Async) Send(message string) {
	a.SendWithLevel(LevelInfo, message)
}

func (a *Async) SendWithTitle(title string, message string) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := a.notifier.SendWithTitle(ctx, title, message); err != nil {
			slog.Error("Async notification failed", "title", title, "error", err)
		}
	}()
}

func (a *Async) SendWithLevel(level Level, message string) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := a.notifier.SendWithLevel(ctx, level, message); err != nil {
			slog.Error("Async notification failed", "level", string(level), "error", err)
		}
	}()
}

// --- Helper ---

// ParseOrigins splits a comma-separated string of origins.
func ParseOrigins(origins string) []string {
	return strings.Split(origins, ",")
}
