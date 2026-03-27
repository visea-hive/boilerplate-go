package validate

import (
	"unicode"

	"github.com/visea-hive/auth-core/pkg/messages"
)

// Strength represents how strong a password is.
type Strength int

const (
	StrengthWeak       Strength = iota // score 0–2
	StrengthFair                       // score 3–4
	StrengthStrong                     // score 5
	StrengthVeryStrong                 // score 6
)

// String returns a human-readable label for the strength level.
func (s Strength) String() string {
	switch s {
	case StrengthFair:
		return "fair"
	case StrengthStrong:
		return "strong"
	case StrengthVeryStrong:
		return "very strong"
	default:
		return "weak"
	}
}

// ScoreResult holds the outcome of a password strength evaluation.
type ScoreResult struct {
	Score    int      // 0–6
	Strength Strength // derived label
	Feedback []string // tips for improving the password
}

// Score evaluates the strength of a password and returns a ScoreResult.
//
// Scoring criteria (1 point each):
//   - Length >= 8
//   - Length >= 12  (bonus)
//   - Contains a lowercase letter
//   - Contains an uppercase letter
//   - Contains a digit
//   - Contains a special character
func Score(password string) ScoreResult {
	points := 0
	var feedback []string

	var hasLower, hasUpper, hasDigit, hasSpecial bool
	for _, r := range password {
		switch {
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	if len(password) >= 8 {
		points++
	} else {
		feedback = append(feedback, "use at least 8 characters")
	}

	if len(password) >= 12 {
		points++ // bonus
	} else {
		feedback = append(feedback, "12+ characters makes your password significantly stronger")
	}

	if hasLower {
		points++
	} else {
		feedback = append(feedback, "add lowercase letters")
	}

	if hasUpper {
		points++
	} else {
		feedback = append(feedback, "add uppercase letters")
	}

	if hasDigit {
		points++
	} else {
		feedback = append(feedback, "add numbers")
	}

	if hasSpecial {
		points++
	} else {
		feedback = append(feedback, "add special characters (e.g. !@#$%)")
	}

	return ScoreResult{
		Score:    points,
		Strength: toStrength(points),
		Feedback: feedback,
	}
}

// Strong returns nil if the password is at least Fair strength (score >= 3).
// Returns messages.ErrPasswordTooWeak otherwise.
func Strong(password string) error {
	result := Score(password)
	if result.Score < 3 {
		return messages.ErrPasswordTooWeak
	}
	return nil
}

func toStrength(score int) Strength {
	switch {
	case score >= 6:
		return StrengthVeryStrong
	case score >= 5:
		return StrengthStrong
	case score >= 3:
		return StrengthFair
	default:
		return StrengthWeak
	}
}
