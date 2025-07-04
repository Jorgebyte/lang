package lang

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
)

// Manager handles a collection of languages for a plugin or server.
// It is safe for concurrent use.
type Manager struct {
	mu          sync.RWMutex
	languages   map[string]*Language
	defaultLang *Language
	logger      *slog.Logger
}

// NewManager creates a new language manager.
// The logger is optional but recommended for debugging translation issues. If no
// logger is provided, a default one writing to os.Stdout is used.
func NewManager(logger *slog.Logger) *Manager {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	return &Manager{
		languages: make(map[string]*Language),
		logger:    logger,
	}
}

// Register adds a language to the manager. If a language with the same
// locale already exists, it is overwritten.
func (m *Manager) Register(lang *Language) error {
	if lang == nil || lang.Locale() == "" {
		return fmt.Errorf("language and its locale cannot be nil or empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.languages[lang.Locale()] = lang

	// If no default language is set, the first registered one becomes the default.
	if m.defaultLang == nil {
		m.defaultLang = lang
		m.logger.Info("Default language set automatically", "locale", lang.Locale())
	}
	return nil
}

// SetDefault sets the default language to use as a fallback. The language for
// the given locale must be registered first.
func (m *Manager) SetDefault(locale string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	lang, ok := m.languages[locale]
	if !ok {
		return fmt.Errorf("cannot set default to an unregistered language: %s", locale)
	}
	m.defaultLang = lang
	return nil
}

// Language returns a registered language by its locale.
func (m *Manager) Language(locale string) (*Language, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	lang, ok := m.languages[locale]
	return lang, ok
}

// Languages returns a slice of all registered languages.
func (m *Manager) Languages() []*Language {
	m.mu.RLock()
	defer m.mu.RUnlock()
	langs := make([]*Language, 0, len(m.languages))
	for _, lang := range m.languages {
		langs = append(langs, lang)
	}
	return langs
}

func (m *Manager) Translate(p Player, key string, placeholders P) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.defaultLang == nil {
		m.logger.Error("Translation failed: no default language configured", "key", key)
		return key
	}

	var text string
	var found bool

	// Attempt to find a translation in the player's language.
	if playerLang, ok := m.languages[p.Locale()]; ok {
		text, found = playerLang.translation(key)
	}

	if !found {
		text, found = m.defaultLang.translation(key)
	}

	if !found {
		m.logger.Warn("Translation key not found", "key", key, "fallback_locale", m.defaultLang.Locale())
		return key
	}

	if len(placeholders) > 0 {
		// strings.NewReplacer is highly efficient for multiple replacements.
		args := make([]string, 0, len(placeholders)*2)
		for k, v := range placeholders {
			args = append(args, k, v)
		}
		replacer := strings.NewReplacer(args...)
		return replacer.Replace(text)
	}

	return text
}
