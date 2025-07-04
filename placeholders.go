package lang

// P is an alias for a map of placeholders. It improves readability when calling
// the Translate function.
// Example: lang.P{"{player}": "Steve", "{world}": "overworld"}
type P map[string]string
