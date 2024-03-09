package types

type Language map[string]LanguageData

type LanguageData struct {
	UI        bool
	Audio     bool
	Subtitles bool
}
