package types

type PCGW struct {
	Name         string                  // Set by the name sanitizer code
	Subtitles    bool                    // Set by the language extraction code
	Languages    map[string]LanguageData // Extracted from the `SupportedLanguages` string
	Stores       map[string]Store        // Scrapped from IsThereAnyDeals
	Ratings      map[string]Rating       // Scraped from IsThereAnyDeals
	Genres       string                  // Scraped from App Tags (Taxonomy on PCGW)
	Franchise    string                  // Scraped from Steam Store (Series on PCGW)
	Pacing       string                  // Scraped from App Tags (Taxonomy on PCGW)
	Perspectives string                  // Scraped from App Tags (Taxonomy on PCGW)
	Controls     string                  // Scraped from App Tags (Taxonomy on PCGW)
	Sports       string                  // Scraped from App Tags (Taxonomy on PCGW)
	Vehicles     string                  // Scraped from App Tags (Taxonomy on PCGW)
	ArtStyles    string                  // Scraped from App Tags (Taxonomy on PCGW)
	Themes       string                  // Scraped from App Tags (Taxonomy on PCGW)
}
