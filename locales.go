package lang

// MinecraftLocales contains all officially supported language codes (locales) for
// Minecraft Bedrock.
// See: https://github.com/Mojang/bedrock-samples/blob/main/resource_pack/texts/language_names.json
var MinecraftLocales = map[string]struct{}{
	"en_US": {}, "en_GB": {}, "de_DE": {}, "es_ES": {}, "es_MX": {}, "fr_FR": {},
	"fr_CA": {}, "it_IT": {}, "ja_JP": {}, "ko_KR": {}, "pt_BR": {}, "pt_PT": {},
	"ru_RU": {}, "zh_CN": {}, "zh_TW": {}, "nl_NL": {}, "bg_BG": {}, "cs_CZ": {},
	"da_DK": {}, "el_GR": {}, "fi_FI": {}, "hu_HU": {}, "id_ID": {}, "nb_NO": {},
	"pl_PL": {}, "sk_SK": {}, "sv_SE": {}, "tr_TR": {}, "uk_UA": {},
}

// IsMinecraftLocale checks if a locale is officially supported.
func IsMinecraftLocale(locale string) bool {
	_, ok := MinecraftLocales[locale]
	return ok
}
