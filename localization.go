package main

import (
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"os"
)

// Localization holds the localizer and supported languages
type Localization struct {
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
}

// NewLocalization initializes the localization with supported languages
func NewLocalization(defaultLang string, supportedLangs []string) *Localization {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", i18n.UnmarshalJSON)

	for _, lang := range supportedLangs {
		bundle.MustLoadMessageFile(fmt.Sprintf("active.%s.json", lang))
	}

	localizer := i18n.NewLocalizer(bundle, defaultLang)
	return &Localization{bundle: bundle, localizer: localizer}
}

// LocalizeMessage localizes a message based on the user's preferred language
func (l *Localization) LocalizeMessage(messageID string, lang string) string {
	localizer := i18n.NewLocalizer(l.bundle, lang)
	translated, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: messageID})
	if err != nil {
		return messageID
	}
	return translated
}

func main() {
	defaultLang := "en"
	supportedLangs := []string{"en", "es", "fr", "de", "zh"}

	localization := NewLocalization(defaultLang, supportedLangs)

	// Example usage
	userLang := os.Getenv("USER_LANG")
	if userLang == "" {
		userLang = defaultLang
	}

	messageID := "hello_world"
	translatedMessage := localization.LocalizeMessage(messageID, userLang)
	fmt.Println(translatedMessage)
}
