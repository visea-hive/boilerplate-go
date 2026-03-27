package validate

import (
	"testing"

	"github.com/visea-hive/auth-core/pkg/messages"
)

func TestScore_AllCriteriaMet(t *testing.T) {
	result := Score("MySecret123!")
	if result.Score != 6 {
		t.Errorf("expected score 6, got %d", result.Score)
	}
	if result.Strength != StrengthVeryStrong {
		t.Errorf("expected VeryStrong, got %s", result.Strength)
	}
	if len(result.Feedback) != 0 {
		t.Errorf("expected no feedback, got: %v", result.Feedback)
	}
}

func TestScore_ShortWeakPassword(t *testing.T) {
	result := Score("abc")
	if result.Score > 2 {
		t.Errorf("expected weak score (<=2), got %d", result.Score)
	}
	if result.Strength != StrengthWeak {
		t.Errorf("expected Weak, got %s", result.Strength)
	}
}

func TestScore_LengthBonus(t *testing.T) {
	// Has all char types, length >= 12 → max score 6
	r1 := Score("Short1!")      // length 7  → no length points
	r2 := Score("LongEnough1!") // length 12 → both length points
	if r2.Score <= r1.Score {
		t.Errorf("longer password should score higher: short=%d long=%d", r1.Score, r2.Score)
	}
}

func TestScore_NoSpecialChar(t *testing.T) {
	result := Score("Password123")
	// missing special char → max 5 points
	if result.Score > 5 {
		t.Errorf("expected max 5 without special char, got %d", result.Score)
	}
	found := false
	for _, f := range result.Feedback {
		if f == "add special characters (e.g. !@#$%)" {
			found = true
		}
	}
	if !found {
		t.Error("expected feedback about missing special characters")
	}
}

func TestStrong_PassesForFairPassword(t *testing.T) {
	// 8 chars, lower+upper+digit = score 4 (fair)
	err := Strong("Abc1defg")
	if err != nil {
		t.Errorf("Strong() returned error for fair password: %v", err)
	}
}

func TestStrong_FailsForWeakPassword(t *testing.T) {
	err := Strong("abc")
	if err == nil {
		t.Error("Strong() expected error for weak password, got nil")
	}
	if err != messages.ErrPasswordTooWeak {
		t.Errorf("expected ErrPasswordTooWeak, got: %v", err)
	}
}

func TestStrength_String(t *testing.T) {
	cases := []struct {
		s    Strength
		want string
	}{
		{StrengthWeak, "weak"},
		{StrengthFair, "fair"},
		{StrengthStrong, "strong"},
		{StrengthVeryStrong, "very strong"},
	}
	for _, c := range cases {
		if got := c.s.String(); got != c.want {
			t.Errorf("Strength(%d).String() = %q, want %q", c.s, got, c.want)
		}
	}
}
