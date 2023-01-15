# Steam 2 PCGW

The goal of this tool is to simplify the process of creating new articles for the PC Gaming Wiki by just simply entering the Steam App ID.

## Why Should I Use This?

Well, it's certainly up to you to decide whether you want to use it or not, and no one will force you to use it, but there are good reasons why you should, even if for certain templates:

- Games with massive DLCs, such as The Sims/Euro Truck Simulator 2/American Truck Simulator/Cities Skylines to name a few.
- Parses most system requirements FLAWLESSLY!
- Fills in most of the information available from the API which you may copy-paste and modify as you go!
- Makes it harder to not skip something while filling in the information.

## Version

v0.0.74 Release Candidate 4

## How To

1. Either visit the Releases page and download the latest build from there.
2. Clone or download the ZIP (Press 'Code' → 'Download ZIP'), enter the directory and type `go run`.  Do note, you'll require Go (<https://go.dev/doc/install>) to do this!

## Contributions

- You are welcome to contribute and improve the code as you see fit.
- If you wish to discuss your plans for the repo, then please make an issue first.

## Plans

- [x] Convert into a CLI app
- [ ] Clean-up the code (underway)
- [ ] Ability to pass in more than one app ID at a time/ parsing more than one article in one go
- [ ] Add support for other APIs to make the data output more complete
- [x] Save cache in a sub-folder, fetch new data if cache is older than seven days (partially done)
- [ ] IGDB and OpenCritic support
- [ ] Download game covers

### Article Status

- [x] Marks the article as stub
- [x] Infobox: Game Cover (needs manual review)
- [x] Infobox: Developers
- [x] Infobox: Publishers
- [x] Infobox: Release Date
- [x] Infobox: Reception: Metacritic (if available)
- [ ] Infobox: Reception: OpenCritic (NOT AVAILABLE)
- [ ] Infobox: Reception: IGDB (NOT AVAILABLE)
- [x] Infobox: Taxomony: F2P / One-time Game Purchase
- [ ] Infobox: Taxonomy: Microtransactions
- [x] Infobox: Taxonomy: Modes (Singleplayer and Multiplayer)
- [ ] Infobox: Taxonomy: Pacing
- [ ] Infobox: Taxonomy: Perspectives
- [ ] Infobox: Taxonomy: Controls
- [x] Infobox: Taxonomy: Genres (and removes genres not valid for PCGW)
- [ ] Infobox: Taxonomy: Sports
- [ ] Infobox: Taxonomy: Vehicles
- [ ] Infobox: Taxonomy: Art Styles
- [x] Infobox: Taxonomy: Themes (detects Adult only)
- [x] Infobox: Taxonomy: Series (detected)
- [x] Infobox: Steam App ID
- [ ] Infobox: GOG App ID
- [x] Infobox: Official Website (or Support Website as a fallback)
- [ ] Infobox: HLTB
- [ ] Infobox: IGDB (Only needs to be set if there is no IGDB reception row, Empty by default for now)
- [ ] Infobox: Lutris
- [ ] Infobox: MobyGames
- [ ] Infobox: StrategyWiki
- [ ] Infobox: Wikipedia
- [ ] Infobox: WineHQ
- [ ] Infobox: License (defaults to Commercial for now)
- [x] Introduction: Introduction
- [x] Introduction: Release History (Generic)
- [x] Introduction: Current State (EMPTY)
- [x] Availability: Steam (Game editions are automatically added)
- [ ] Availability: Other Stores (Outside the scope of the project)
- [ ] Monetization: Ad-Supported
- [ ] Monetization: DLC
- [ ] Monetization: Expansion Pack
- [ ] Monetization: freeware
- [x] Monetization: free-to-play (F2P / One-time Game Purchase)
- [ ] Monetization: sponsored
- [ ] Monetization: subscription
- [ ] Microtransactions: Microtransactions
- [x] Microtransactions: DLCs
- [x] Game Data: Config File Location (Add file location)
- [x] Save Game Data: File location (Add file location)
- [x] Save Game Sync (Steam cloud detected! Add file location)
- [ ] Video
- [ ] Input: Key remapping
- [ ] Input: Touchscreen
- [x] Input: Controller Support, Full Controller
- [ ] Input: Controller (PS/Xbox/Others) (IMPOSSIBLE)
- [ ] Audio
- [x] Languages (There maybe some discrepancies as Steam API provides very vague info - No discrepancies reported so far)
- [ ] API (App executables are guessed from the system specifications - mostly accurate)
- [ ] Middleware
- [x] System Requirements: Windows (CPU and GPU sections may need review)
- [x] System Requirements: Mac (CPU and GPU sections may need review)
- [x] System Requirements: Linux (CPU and GPU sections may need review)
- [x] References

## Special Thanks

- Dandelion Sprout - first contribution, vital feedback and testing
- Baron Smoki - vital feedback and testing
- Dave247 - vital feedback and testing
- Mine18 - vital feedback and testing
- Mrtnptrs - vital feedback and testing
- mastan - vital feedback and testing
