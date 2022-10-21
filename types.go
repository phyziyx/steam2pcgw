package main

// JSON related

type Game map[string]GameValue

type GameValue struct {
	Success bool `json:"success"`
	Data    Data `json:"data"`
}

type Data struct {
	Type                string             `json:"type"`
	Name                string             `json:"name"`
	SteamAppid          int64              `json:"steam_appid"`
	RequiredAge         int64              `json:"required_age"`
	IsFree              bool               `json:"is_free"`
	ControllerSupport   string             `json:"controller_support"`
	Dlc                 []int64            `json:"dlc,omitempty"`
	DetailedDescription string             `json:"detailed_description"`
	AboutTheGame        string             `json:"about_the_game"`
	ShortDescription    string             `json:"short_description"`
	SupportedLanguages  string             `json:"supported_languages"`
	HeaderImage         string             `json:"header_image"`
	Website             *string            `json:"website"`
	PCRequirements      Requirements       `json:"pc_requirements"`
	MACRequirements     Requirements       `json:"mac_requirements"`
	LinuxRequirements   Requirements       `json:"linux_requirements"`
	Developers          []string           `json:"developers"`
	Publishers          []string           `json:"publishers"`
	Packages            []int64            `json:"packages"`
	PackageGroups       []PackageGroup     `json:"package_groups"`
	Platforms           Platforms          `json:"platforms"`
	Metacritic          *Metacritic        `json:"metacritic,omitempty"`
	Categories          []Category         `json:"categories"`
	Genres              []Genre            `json:"genres"`
	Screenshots         []Screenshot       `json:"screenshots"`
	Movies              []Movie            `json:"movies"`
	Recommendations     Recommendations    `json:"recommendations"`
	Achievements        Achievements       `json:"achievements"`
	ReleaseDate         ReleaseDate        `json:"release_date"`
	SupportInfo         SupportInfo        `json:"support_info"`
	Background          string             `json:"background"`
	BackgroundRaw       string             `json:"background_raw"`
	ContentDescriptors  ContentDescriptors `json:"content_descriptors"`
	LegalNotice         *string            `json:"legal_notice,omitempty"`
	PriceOverview       *PriceOverview     `json:"price_overview,omitempty"`
}

type Achievements struct {
	Total       int64         `json:"total"`
	Highlighted []Highlighted `json:"highlighted"`
}

type Highlighted struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type Category struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
}

type ContentDescriptors struct {
	IDS   []int64 `json:"ids"`
	Notes *string `json:"notes"`
}

type Genre struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type Requirements struct {
	Minimum     string  `json:"minimum"`
	Recommended *string `json:"recommended,omitempty"`
}

type Metacritic struct {
	Score int64  `json:"score"`
	URL   string `json:"url"`
}

type Movie struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Thumbnail string `json:"thumbnail"`
	Webm      Mp4    `json:"webm"`
	Mp4       Mp4    `json:"mp4"`
	Highlight bool   `json:"highlight"`
}

type Mp4 struct {
	The480 string `json:"480"`
	Max    string `json:"max"`
}

type PackageGroup struct {
	Name                    string `json:"name"`
	Title                   string `json:"title"`
	Description             string `json:"description"`
	SelectionText           string `json:"selection_text"`
	SaveText                string `json:"save_text"`
	DisplayType             int64  `json:"display_type"`
	IsRecurringSubscription string `json:"is_recurring_subscription"`
	Subs                    []Sub  `json:"subs"`
}

type Sub struct {
	Packageid                int64  `json:"packageid"`
	PercentSavingsText       string `json:"percent_savings_text"`
	PercentSavings           int64  `json:"percent_savings"`
	OptionText               string `json:"option_text"`
	OptionDescription        string `json:"option_description"`
	CanGetFreeLicense        string `json:"can_get_free_license"`
	IsFreeLicense            bool   `json:"is_free_license"`
	PriceInCentsWithDiscount int64  `json:"price_in_cents_with_discount"`
}

type Platforms struct {
	Windows bool `json:"windows"`
	MAC     bool `json:"mac"`
	Linux   bool `json:"linux"`
}

type PriceOverview struct {
	Currency         string `json:"currency"`
	Initial          int64  `json:"initial"`
	Final            int64  `json:"final"`
	DiscountPercent  int64  `json:"discount_percent"`
	InitialFormatted string `json:"initial_formatted"`
	FinalFormatted   string `json:"final_formatted"`
}

type Recommendations struct {
	Total int64 `json:"total"`
}

type ReleaseDate struct {
	ComingSoon bool   `json:"coming_soon"`
	Date       string `json:"date"`
}

type Screenshot struct {
	ID            int64  `json:"id"`
	PathThumbnail string `json:"path_thumbnail"`
	PathFull      string `json:"path_full"`
}

type SupportInfo struct {
	URL   string `json:"url"`
	Email string `json:"email"`
}

// Language

type Language map[string]LanguageValue

type LanguageValue struct {
	UI        bool
	Audio     bool
	Subtitles bool
}
