package lang

// Player represents an entity, such as a player, that has a language setting.
// The host application (e.g., a Dragonfly server) should provide an object that
// implements this interface.
type Player interface {
	// Locale returns the language code of the player for example, "en_US".
	Locale() string
}
