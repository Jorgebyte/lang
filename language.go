package lang

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Language holds all translations for a specific locale.
type Language struct {
	locale       string
	name         string
	translations map[string]string
}

// NewLanguageFromData creates a new Language object from pre-existing data.
func NewLanguageFromData(locale, name string, translations map[string]string) *Language {
	return &Language{
		locale:       locale,
		name:         name,
		translations: translations,
	}
}

// LoadLanguageFile loads a language file (JSON or YAML) from a given path.
// The locale is inferred from the filename (e.g., "en_US.json").
func LoadLanguageFile(path string) (*Language, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read language file %s: %w", path, err)
	}

	ext := filepath.Ext(path)
	locale := strings.TrimSuffix(filepath.Base(path), ext)

	if !IsMinecraftLocale(locale) {
		return nil, fmt.Errorf("locale '%s' derived from filename is not a valid Minecraft locale", locale)
	}

	var translations map[string]string

	// Select the appropriate parser based on the file extension \\
	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &translations); err != nil {
			return nil, fmt.Errorf("failed to decode JSON from %s: %w", path, err)
		}
	case ".yml", ".yaml":
		if err := yaml.Unmarshal(data, &translations); err != nil {
			return nil, fmt.Errorf("failed to decode YAML from %s: %w", path, err)
		}
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}

	langName, ok := translations["language.name"]
	if !ok {
		return nil, fmt.Errorf("language file %s is missing the required 'language.name' key", path)
	}
	delete(translations, "language.name")

	return NewLanguageFromData(locale, langName, translations), nil
}

// Locale returns the language code e.g., "en_US".
func (l *Language) Locale() string {
	return l.locale
}

// Name returns the display name of the language e.g., "English (US)".
func (l *Language) Name() string {
	return l.name
}

// translation looks up a translation key. It returns the translated string and
// a boolean indicating if the key was found.
func (l *Language) translation(key string) (string, bool) {
	val, ok := l.translations[key]
	return val, ok
}
