package messages

import (
	"errors"
	"fmt"
	"strings"
)

// --- Supported Languages ---
const (
	LangID      = "id"
	LangEN      = "en"
	LangDefault = LangEN // used when Accept-Language header is absent or unrecognised
)

// ParseLang normalises an Accept-Language header value to a supported language
// code. Returns LangDefault ("en") when the header is empty or not recognised.
//
// Usage in middleware:
//
//	lang := messages.ParseLang(c.GetHeader("Accept-Language"))
//	messages.Translate(lang, err)
func ParseLang(header string) string {
	switch {
	case strings.HasPrefix(header, LangID):
		return LangID
	case strings.HasPrefix(header, LangEN):
		return LangEN
	default:
		return LangDefault
	}
}

// --- Internal Registries ---
var (
	errorTranslations   = make(map[error]map[string]string)
	successTranslations = make(map[SuccessType]map[string]string)
)

// NewError defines a new standard error with its English translation registered automatically.
// The primary `.Error()` output will use the Indonesian (LangID) translation.
func NewError(idMsg, enMsg string) error {
	err := errors.New(idMsg)
	errorTranslations[err] = map[string]string{
		LangID: idMsg,
		LangEN: enMsg,
	}
	return err
}

// --- Dynamic Title Helper ---

// FormatTitle returns the provided title if it's not empty, otherwise defaults to "data".
func FormatTitle(lang, title string) string {
	if strings.TrimSpace(title) == "" {
		if lang == LangEN {
			return "data"
		}
		return "data"
	}
	return title
}

// --- Success Messages ---

type SuccessType string

func newSuccess(op string, idMsg string, enMsg string) SuccessType {
	s := SuccessType(op)
	successTranslations[s] = map[string]string{
		LangID: idMsg,
		LangEN: enMsg,
	}
	return s
}

// Translate returns a localized string for either an error or a SuccessType.
// For SuccessType, you can provide an optional title argument as the third parameter.
func Translate(lang string, msg interface{}, args ...string) string {
	switch v := msg.(type) {
	case error:
		if v == nil {
			return ""
		}

		// Iterate through registered errors using errors.Is to properly match wrapped errors
		for targetErr, langMap := range errorTranslations {
			if errors.Is(v, targetErr) {
				if translation, ok := langMap[lang]; ok {
					return translation
				}
				return langMap[LangID] // Fallback to ID translation
			}
		}

		// Fallback to original error string if no translation is found at all
		return v.Error()

	case SuccessType:
		title := ""
		if len(args) > 0 {
			title = args[0]
		}

		langMap, ok := successTranslations[v]
		if !ok {
			// Fallback if an unknown operation is requested
			return fmt.Sprintf("%s success", FormatTitle(lang, title))
		}

		template, ok := langMap[lang]
		if !ok {
			template = langMap[LangID] // Fallback to ID
		}

		return fmt.Sprintf(template, FormatTitle(lang, title))
	}

	return ""
}
