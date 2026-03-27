package notifier

import (
	"context"
	"testing"
	"time"
)

// --- NoOpNotifier Tests ---

func TestNoOpNotifier_Send(t *testing.T) {
	n := NewNoOpNotifier()

	err := n.Send(context.Background(), "test message")
	if err != nil {
		t.Errorf("NoOpNotifier.Send() returned error: %v", err)
	}
}

func TestNoOpNotifier_SendWithTitle(t *testing.T) {
	n := NewNoOpNotifier()

	err := n.SendWithTitle(context.Background(), "Test Title", "test message")
	if err != nil {
		t.Errorf("NoOpNotifier.SendWithTitle() returned error: %v", err)
	}
}

// --- MockNotifier for testing consumers ---

// MockNotifier records all notifications sent, useful for testing
// services that depend on the Notifier interface.
type MockNotifier struct {
	Messages []MockMessage
}

type MockMessage struct {
	Level   Level
	Title   string
	Message string
}

func NewMockNotifier() *MockNotifier {
	return &MockNotifier{}
}

func (m *MockNotifier) Send(ctx context.Context, message string) error {
	m.Messages = append(m.Messages, MockMessage{Level: LevelInfo, Message: message})
	return nil
}

func (m *MockNotifier) SendWithTitle(ctx context.Context, title string, message string) error {
	m.Messages = append(m.Messages, MockMessage{Title: title, Message: message})
	return nil
}

func (m *MockNotifier) SendWithLevel(ctx context.Context, level Level, message string) error {
	m.Messages = append(m.Messages, MockMessage{Level: level, Message: message})
	return nil
}

// --- Test MockNotifier behavior ---

func TestMockNotifier_RecordsMessages(t *testing.T) {
	mock := NewMockNotifier()

	_ = mock.Send(context.Background(), "hello")
	_ = mock.SendWithTitle(context.Background(), "Alert", "something happened")
	_ = mock.SendWithLevel(context.Background(), LevelCritical, "database down")

	if len(mock.Messages) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(mock.Messages))
	}

	if mock.Messages[0].Message != "hello" || mock.Messages[0].Level != LevelInfo {
		t.Errorf("expected info message 'hello', got '%s' with level '%s'", mock.Messages[0].Message, mock.Messages[0].Level)
	}

	if mock.Messages[1].Title != "Alert" {
		t.Errorf("expected title 'Alert', got '%s'", mock.Messages[1].Title)
	}

	if mock.Messages[1].Message != "something happened" {
		t.Errorf("expected message 'something happened', got '%s'", mock.Messages[1].Message)
	}

	if mock.Messages[2].Message != "database down" || mock.Messages[2].Level != LevelCritical {
		t.Errorf("expected critical message 'database down', got '%v'", mock.Messages[2])
	}
}

// --- Async Tests ---

func TestAsync_SendDoesNotBlock(t *testing.T) {
	mock := NewMockNotifier()
	async := NewAsync(mock)

	// Send should return immediately (non-blocking)
	start := time.Now()
	async.Send("async message")
	elapsed := time.Since(start)

	if elapsed > 100*time.Millisecond {
		t.Errorf("Async.Send() took too long: %v (should be instant)", elapsed)
	}

	// Wait a bit for the goroutine to finish
	time.Sleep(100 * time.Millisecond)

	if len(mock.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(mock.Messages))
	}
	if mock.Messages[0].Message != "async message" {
		t.Errorf("expected 'async message', got '%s'", mock.Messages[0].Message)
	}
}

func TestAsync_SendWithTitleDoesNotBlock(t *testing.T) {
	mock := NewMockNotifier()
	async := NewAsync(mock)

	start := time.Now()
	async.SendWithTitle("Deploy", "v1.2.0 deployed to production")
	elapsed := time.Since(start)

	if elapsed > 100*time.Millisecond {
		t.Errorf("Async.SendWithTitle() took too long: %v", elapsed)
	}

	time.Sleep(100 * time.Millisecond)

	if len(mock.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(mock.Messages))
	}
	if mock.Messages[0].Title != "Deploy" {
		t.Errorf("expected title 'Deploy', got '%s'", mock.Messages[0].Title)
	}
}

func TestAsync_SendWithLevelDoesNotBlock(t *testing.T) {
	mock := NewMockNotifier()
	async := NewAsync(mock)

	start := time.Now()
	async.SendWithLevel(LevelError, "something broke")
	elapsed := time.Since(start)

	if elapsed > 100*time.Millisecond {
		t.Errorf("Async.SendWithLevel() took too long: %v", elapsed)
	}

	time.Sleep(100 * time.Millisecond)

	if len(mock.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(mock.Messages))
	}
	if mock.Messages[0].Level != LevelError {
		t.Errorf("expected level Error, got '%s'", mock.Messages[0].Level)
	}
}

// --- Context Cancellation Test ---

func TestWebhookNotifier_RespectsContextCancellation(t *testing.T) {
	// Use a fake URL — the request should be cancelled before it reaches the server
	n := NewWebhookNotifier("https://httpbin.org/delay/10", "test")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := n.Send(ctx, "should timeout")
	if err == nil {
		t.Error("expected error due to context timeout, got nil")
	}
}
