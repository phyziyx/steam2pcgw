# Steam 2 PCGW

The goal of this tool is to simplify the process of creating new articles for the PC Gaming Wiki by just simply entering the Steam App ID.

## Version

v0.0.69

## How To

1. Either visit the Releases page and download the latest build from there.
2. Clone or download the ZIP (Press 'Code' → 'Download ZIP'), enter the directory and type `go run`.  Do note, you'll require Go (https://go.dev/doc/install) to do this!

## Contributions

- You are welcome to contribute and improve the code as you see fit.
- If you wish to discuss your plans for the repo, then please make an issue first.

## Plans

- [ ] Convert into a CLI app
- [ ] Clean-up the code
- [ ] Ability to pass in more than one app ID at a time/ parsing more than one article in one go
- [ ] Add support for other APIs to make the data output more complete
- [ ] Save cache in a sub-folder, fetch new data if cache is older than seven days
- [ ] IGDB and OpenCritic support
- [ ] Download game covers

### Article Status

- [x] Marks the article as stub
- [x] Infobox: Game Cover (needs manual review)
- [x] Infobox: Developers
- [x] Infobox: Publishers
- [x] Infobox: Release Date
- [x] Infobox: Reception: Metacritic (if available)
- [ ] Infobox: Reception: OpenCritic (if available)
- [ ] Infobox: Reception: IGDB (if available)
- [x] Infobox: Taxomony: F2P / One-time Game Purchase
- [ ] Infobox: Taxonomy: Microtransactions
- [ ] Infobox: Taxonomy: Modes (defaults to Singleplayer for now)
- [ ] Infobox: Taxonomy: Pacing
- [ ] Infobox: Taxonomy: Perspectives
- [ ] Infobox: Taxonomy: Controls
- [x] Infobox: Taxonomy: Genres
- [ ] Infobox: Taxonomy: Sports
- [ ] Infobox: Taxonomy: Vehicles
- [ ] Infobox: Taxonomy: Art Styles
- [ ] Infobox: Taxonomy: Themes
- [ ] Infobox: Taxonomy: Series
- [x] Infobox: Steam App ID
- [ ] Infobox: GOG App ID
- [x] Infobox: Official Website (or Support Website, whichever is available)
- [ ] Infobox: HLTB
- [ ] Infobox: IGDB (Only needs to be set if there is no IGDB reception row, Empty by default for now)
- [ ] Infobox: Lutris
- [ ] Infobox: MobyGames
- [ ] Infobox: StrategyWiki
- [ ] Infobox: Wikipedia
- [ ] Infobox: WineHQ
- [ ] Infobox: License (defaults to Commercial for now)
- [x] Introduction: Introduction
- [x] Introduction: Release History (generic)
- [ ] Introduction: Current State (impossible?)
- [x] Availability: Steam
- [ ] Availability: Other Stores (not the scope of the project yet)
- [ ] Monetization: Ad-Supported
- [ ] Monetization: DLC
- [ ] Monetization: Expansion Pack
- [ ] Monetization: freeware
- [x] Monetization: free-to-play (F2P / One-time Game Purchase)
- [ ] Monetization: sponsored
- [ ] Monetization: subscription
- [ ] Microtransactions: Microtransactions
- [x] Microtransactions: DLCs
- [x] Game Data: Config File Location (needs manual review)
- [x] Save Game Data: File location (needs manual review)
- [x] Save Game Sync (Steam cloud detected! needs manual review)
- [ ] Video
- [ ] Input: Key remapping
- [ ] Input: Touchscreen
- [x] Input: Controller Support, Full Controller
- [ ] Input: Controller (PS/Xbox/Others)
- [ ] Audio
- [x] Languages (There maybe some discrepancies as Steam API provides very vague info)
- [ ] API (App executable bits detected, needs manual review)
- [ ] Middleware
- [x] System Requirements: Windows (Needs manual review)
- [x] System Requirements: Mac (Needs manual review)
- [x] System Requirements: Linux (Needs manual review)
- [x] References

## Special Thanks

- Dandelion Sprout - first contribution, vital feedback and testing
- Baron Smoki - vital feedback and testing
- Dave247 - vital feedback and testing
