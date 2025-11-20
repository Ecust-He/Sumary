package i18n

type Ii18nOpr interface {
	// Translate Get dist-language value by input key from dict. if not found,
	// return original key.
	// and then return dist-language sentence by vars_template filled with
	// addition vars. ${var} in vars_template will be replaced by vars.get(var).
	Translate(language, key string, vars ...string) string
	// PutValue dynamic to put i18n value to cache
	PutValue(language, key, value string) error
	// RemoveValue dynamic to remove i18n value to cache
	RemoveValue(language, key string) error
}
