package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

func main() {

	var gameJson []byte
	gameJson, err = ParseGame(gameId)
	if err != nil {
		fmt.Println(err)
		return
	}

	game, err := UnmarshalGame(gameJson)
	if err != nil {
		fmt.Printf("An error occurred while attempting to unmarshal the JSON... (%s)\n", err)
		return
	} else if !game.Success {
		fmt.Println("The app ID provided does not exist or does not have a Store page...")
		return
	}

	outputFile, err := os.Create(fmt.Sprintf("output/%s.txt", gameId))
	if err != nil {
		fmt.Println("Failed to create the output file... Process stopped!")
		return
	}

	fmt.Println("* [1/25] Adding stub")
	outputFile.WriteString("{{stub}}\n")

	fmt.Println("* [2/25] Adding app cover")
	outputFile.WriteString(fmt.Sprintf("{{Infobox game\n|cover        = %s cover.jpg", SanitiseName(game.Data.Name, true)))

	fmt.Println("* [3/25] Adding app developers")
	outputFile.WriteString("\n|developers   = ")
	for _, developer := range game.Data.Developers {
		outputFile.WriteString(fmt.Sprintf("\n{{Infobox game/row/developer|%s}}", SanitiseName(developer, false)))
	}

	fmt.Println("* [4/25] Adding app publishers")
	outputFile.WriteString("\n|publishers   = ")
	for _, publisher := range game.Data.Publishers {
		if len(game.Data.Publishers) == 1 {
			skip := false
			for _, developer := range game.Data.Developers {
				if developer == publisher {
					skip = true
					break
				}
			}
			if skip {
				continue
			}
		}
		outputFile.WriteString(fmt.Sprintf("\n{{Infobox game/row/publisher|%s}}", SanitiseName(publisher, false)))
	}

	fmt.Println("* [5/25] Adding app release date")
	outputFile.WriteString("\n|engines      =\n<!-- {{Infobox game/row/engine|}} -->\n|release dates= ")

	date := ""
	if game.HasSteamGenre(EarlyAccess) {
		date += "EA"
	} else if game.Data.ReleaseDate.ComingSoon {
		if success, _ := IsDate(game.Data.ReleaseDate.Date); success {
			date += ParseDate(game.Data.ReleaseDate.Date)
		} else {
			date += "TBA"
		}
	} else {
		date += ParseDate(game.Data.ReleaseDate.Date)
	}

	if game.Data.Platforms.Windows {
		outputFile.WriteString(fmt.Sprintf("\n{{Infobox game/row/date|Windows| %s }}", date))
	}

	if game.Data.Platforms.MAC {
		outputFile.WriteString(fmt.Sprintf("\n{{Infobox game/row/date|OS X| %s }}", date))
	}

	if game.Data.Platforms.Linux {
		outputFile.WriteString(fmt.Sprintf("\n{{Infobox game/row/date|Linux| %s }}", date))
	}

	fmt.Println("* [6/25] Adding reception score")
	outputFile.WriteString("\n|reception    = \n{{Infobox game/row/reception|Metacritic|")
	if game.Data.Metacritic != nil {
		outputFile.WriteString(fmt.Sprintf("%s|%d}}", strings.TrimPrefix(game.Data.Metacritic.URL, "https://metacritic.com/game/pc/"), game.Data.Metacritic.Score))
	} else if val, ok := game.Data.Ratings["Metascore"]; ok {
		outputFile.WriteString(fmt.Sprintf("%s|%d}}", strings.TrimPrefix(val.URL, "https://metacritic.com/game/pc/"), val.Score))
	} else {
		outputFile.WriteString("link|rating}}")
	}

	outputFile.WriteString("\n{{Infobox game/row/reception|OpenCritic|")
	if val, ok := game.Data.Ratings["OpenCritic"]; ok {
		outputFile.WriteString(fmt.Sprintf("%s|%d}}", strings.TrimPrefix(val.URL, "https://opencritic.com/game/"), val.Score))
	} else {
		outputFile.WriteString("link|rating}}")
	}
	outputFile.WriteString("\n{{Infobox game/row/reception|IGDB|link|rating}}")

	outputFile.WriteString("\n|taxonomy     =\n{{Infobox game/row/taxonomy/monetization      | ")
	if game.Data.IsFree {
		fmt.Println("* [7/25] Game is F2P")
		outputFile.WriteString(("Free-to-play }}"))
	} else {
		fmt.Println("* [7/25] Game is not F2P")
		outputFile.WriteString(("One-time game purchase }}"))
	}

	fmt.Println("* [8/25] Taxonomy...")
	outputFile.WriteString("\n{{Infobox game/row/taxonomy/microtransactions | ")
	if !game.HasCategory(InAppPurchases) {
		outputFile.WriteString("None ")
	}
	outputFile.WriteString("}}\n{{Infobox game/row/taxonomy/modes             | ")

	modes := ""

	if game.HasCategory(Singleplayer) {
		modes += "Singleplayer, "
	}

	if game.HasCategory(Multiplayer) {
		modes += "Multiplayer, "
	}

	modes = strings.TrimSuffix(modes, ", ")
	outputFile.WriteString(modes)

	outputFile.WriteString(" }}\n{{Infobox game/row/taxonomy/pacing            | ")
	outputFile.WriteString(game.Data.Pacing)

	outputFile.WriteString(" }}\n{{Infobox game/row/taxonomy/perspectives      | ")
	outputFile.WriteString(game.Data.Perspectives)

	outputFile.WriteString(" }}\n{{Infobox game/row/taxonomy/controls          | ")
	outputFile.WriteString(game.Data.Controls)

	outputFile.WriteString(" }}\n{{Infobox game/row/taxonomy/genres            | ")
	outputFile.WriteString(game.Data.Genres)

	outputFile.WriteString(" }}\n{{Infobox game/row/taxonomy/sports            | ")
	outputFile.WriteString(game.Data.Sports)

	outputFile.WriteString(" }}\n{{Infobox game/row/taxonomy/vehicles          | ")
	outputFile.WriteString(game.Data.Vehicles)

	outputFile.WriteString(" }}\n{{Infobox game/row/taxonomy/art styles        | ")
	outputFile.WriteString(game.Data.ArtStyles)

	outputFile.WriteString(" }}\n{{Infobox game/row/taxonomy/themes            | ")
	outputFile.WriteString(game.Data.Themes)

	outputFile.WriteString(" }}\n{{Infobox game/row/taxonomy/series            | ")
	if len(game.Data.Franchise) != 0 {
		outputFile.WriteString(game.Data.Franchise)
		outputFile.WriteString(" }}\n")
	} else {
		outputFile.WriteString("}}\n")
	}

	outputFile.WriteString(fmt.Sprintf("|steam appid  = %s\n|steam appid side = ", gameId))
	if game.Data.Dlc != nil {
		var dlcs string = ""
		for _, v := range game.Data.Dlc {
			dlcs += fmt.Sprintf("%v, ", v)
		}
		dlcs = strings.TrimSuffix(dlcs, ", ")
		outputFile.WriteString(dlcs)
	}
	outputFile.WriteString("\n|gogcom id    = \n|gogcom id side = \n|official site= ")

	if game.Data.Website != nil {
		outputFile.WriteString(*game.Data.Website)
	} else {
		outputFile.WriteString(game.Data.SupportInfo.URL)
	}
	outputFile.WriteString("\n|hltb         = \n|igdb         = <!-- Only needs to be set if there is no IGDB reception row -->\n|lutris       = \n|mobygames    = \n|strategywiki = \n|wikipedia    = \n|winehq       = \n|license      = commercial\n}}")

	fmt.Println("* [9/25] Processing introduction...")
	outputFile.WriteString("\n\n{{Introduction\n|introduction      = ")
	// outputFile.WriteString(removeTags(game.Data.AboutTheGame))

	outputFile.WriteString("\n\n|release history   = ")

	outputFile.WriteString("\n\n|current state     = ")
	outputFile.WriteString("\n}}")

	outputFile.WriteString("\n\n'''General information'''")
	outputFile.WriteString("\n{{mm}} [https://steamcommunity.com/app/" + gameId + "/discussions/ Steam Community Discussions]")

	fmt.Println("* [10/25] Processing Availability!")

	outputFile.WriteString("\n\n==Availability==\n{{Availability|\n")

	platforms := ""
	if game.Data.Platforms.Windows {
		platforms += "Windows, "
	}
	if game.Data.Platforms.MAC {
		platforms += "OS X, "
	}
	if game.Data.Platforms.Linux {
		platforms += "Linux, "
	}

	platforms = strings.TrimSuffix(platforms, ", ")
	var editions []string
	trimPrice := regexp.MustCompile(`(\$.+ USD)`)

	for _, v := range game.Data.PackageGroups {
		diplayType, _ := GetInt(v.DisplayType)
		if diplayType == 1 {
			continue
		}

		for _, sub := range v.Subs {
			edition := RemoveTags(sub.OptionText, "")
			edition = strings.ReplaceAll(edition, SanitiseName(game.Data.Name, true), "")
			edition = strings.TrimSpace(edition)
			edition = strings.TrimPrefix(edition, ": ")
			edition = strings.TrimPrefix(edition, "- ")
			edition = trimPrice.ReplaceAllLiteralString(edition, "")
			edition = strings.TrimSpace(edition)
			edition = strings.TrimSuffix(edition, " -")

			if len(edition) != 0 {
				editions = append(editions, "'''"+edition+"'''")
			}
		}
	}

	editionList := ""
	for i := 0; i < len(editions); i++ {
		editionList += editions[i]
		if i == len(editions)-2 {
			editionList += " and "
		} else {
			editionList += ", "
		}
	}

	if len(editionList) != 0 {
		editionList = strings.TrimSuffix(editionList, ", ")
		editionList += " also available"
	}

	outputFile.WriteString(fmt.Sprintf("{{Availability/row| Steam | %s | Steam | %s | | %s ", gameId, editionList, platforms))

	if len(game.Data.Packages) == 0 {
		outputFile.WriteString("| unavailable ")
	}

	outputFile.WriteString("}}")

	for store, data := range game.Data.Stores {
		outputFile.WriteString(fmt.Sprintf("\n{{Availability/row| %s | %s | DRM | %s | | %s }}", store, data.URL, editionList, data.Platforms))
	}

	outputFile.WriteString("\n}}")

	// Third party account check
	if len(game.Data.ExternalAccountNotice) != 0 {
		outputFile.WriteString(fmt.Sprintf("\n{{ii}} Requires 3rd-Party Account: %s", game.Data.ExternalAccountNotice))
	}

	// DRM check
	drms := ""
	if strings.Contains(game.Data.DRMNotice, "Denuvo") {
		drms += "{{DRM|Denuvo}}, "
	}
	// if strings.Contains(game.Data.PCRequirements["minimum"].(string), "VMProtect") {
	// 	drms += "{{DRM|VMProtect}}, "
	// }

	drms = strings.TrimSuffix(drms, ", ")
	if len(drms) == 0 {
		drms += game.Data.DRMNotice
		if len(drms) != 0 {
			outputFile.WriteString(fmt.Sprintf("\n{{ii}} All versions require %s.", drms))
		}
	}

	if len(editionList) > 1 {
		outputFile.WriteString("\n\n===Version differences===\n{{ii}} ")
		outputFile.WriteString(editionList)
	}

	outputFile.WriteString("\n\n<!-- PAGE GENERATED BY STEAM2PCGW -->")

	fmt.Println("* [11/25] Processing Monetization!")
	outputFile.WriteString("\n\n==Monetization==\n")

	outputFile.WriteString("{{Monetization")
	outputFile.WriteString("\n|ad-supported                = ")
	outputFile.WriteString("\n|dlc                         = ")
	outputFile.WriteString("\n|expansion pack              = ")
	outputFile.WriteString("\n|freeware                    = ")
	outputFile.WriteString("\n|free-to-play                = ")
	if game.Data.IsFree {
		outputFile.WriteString("The game has such monetization.")
	}
	outputFile.WriteString("\n|one-time game purchase      = ")
	if !game.Data.IsFree {
		outputFile.WriteString("The game requires an upfront purchase to access.")
	}
	outputFile.WriteString("\n|sponsored                   = ")
	outputFile.WriteString("\n|subscription                = ")
	outputFile.WriteString("\n|subscription gaming service = ")
	outputFile.WriteString("\n}}")

	fmt.Println("* [12/25] Processing Microtransactions!")

	outputFile.WriteString("\n\n===Microtransactions===\n{{Microtransactions")

	outputFile.WriteString("\n|boost               = ")
	outputFile.WriteString("\n|cosmetic            = ")
	outputFile.WriteString("\n|currency            = ")
	outputFile.WriteString("\n|finite spend        = ")
	outputFile.WriteString("\n|infinite spend      = ")
	outputFile.WriteString("\n|free-to-grind       = ")
	outputFile.WriteString("\n|loot box            = ")
	outputFile.WriteString("\n|none                = ")
	if !game.HasCategory(InAppPurchases) {
		outputFile.WriteString("None")
	}
	outputFile.WriteString("\n|player trading      = ")
	outputFile.WriteString("\n|time-limited        = ")
	outputFile.WriteString("\n|unlock              = ")
	outputFile.WriteString("\n}}")

	fmt.Println("* [13/25] Processing DLCs!")
	outputFile.WriteString("\n\n{{DLC|\n<!-- DLC rows goes below: -->\n}}")

	fmt.Println("* [14/25] Processing Config File Location!")

	outputFile.WriteString("\n\n==Game data==\n===Configuration file(s) location===")
	outputFile.WriteString("\n{{Game data|")
	if game.Data.Platforms.Windows {
		outputFile.WriteString("\n{{Game data/config|Windows|}}")
	}
	if game.Data.Platforms.MAC {
		outputFile.WriteString("\n{{Game data/config|OS X|}}")
	}
	if game.Data.Platforms.Linux {
		outputFile.WriteString("\n{{Game data/config|Linux|}}")
	}
	outputFile.WriteString("\n}}")

	fmt.Println("* [15/25] Processing Save Game Location!")

	outputFile.WriteString("\n\n===Save game data location===")
	outputFile.WriteString("\n{{Game data|")
	if game.Data.Platforms.Windows {
		outputFile.WriteString("\n{{Game data/saves|Windows|}}")
	}
	if game.Data.Platforms.MAC {
		outputFile.WriteString("\n{{Game data/saves|OS X|}}")
	}
	if game.Data.Platforms.Linux {
		outputFile.WriteString("\n{{Game data/saves|Linux|}}")
	}
	outputFile.WriteString("\n}}")

	fmt.Println("* [16/25] Processing Save Game Sync!")

	outputFile.WriteString("\n\n===[[Glossary:Save game cloud syncing|Save game cloud syncing]]===\n{{Save game cloud syncing\n")
	outputFile.WriteString(`|discord                   =
|discord notes             =
|epic games launcher       =
|epic games launcher notes =
|gog galaxy                =
|gog galaxy notes          =
|origin                    =
|origin notes              =
|steam cloud               = `)

	// Game has steam cloud, then we can just add it
	// Otherwise, we can check if the game is out yet or not
	// to determine whether we should add `unknown` or `false`
	if game.HasCategory(SteamCloud) {
		outputFile.WriteString("true")
	} else {
		if game.Data.ReleaseDate.ComingSoon {
			outputFile.WriteString("unknown")
		} else {
			outputFile.WriteString("false")
		}
	}

	outputFile.WriteString(`
|steam cloud notes         =
|ubisoft connect           =
|ubisoft connect notes     =
|xbox cloud                =
|xbox cloud notes          =
}}`)

	// TODO: Scan the description to search for widescreen, ray tracing etc support
	fmt.Println("* [17/25] Processing Video!")
	outputFile.WriteString("\n\n==Video==\n{{Video\n")
	outputFile.WriteString(`|wsgf link                  =
|widescreen wsgf award      =
|multimonitor wsgf award    =
|ultrawidescreen wsgf award =
|4k ultra hd wsgf award     =
|widescreen resolution      = unknown
|widescreen resolution notes=
|multimonitor               = unknown
|multimonitor notes         =
|ultrawidescreen            = unknown
|ultrawidescreen notes      =
|4k ultra hd                = unknown
|4k ultra hd notes          =
|fov                        = unknown
|fov notes                  =
|windowed                   = unknown
|windowed notes             =
|borderless windowed        = unknown
|borderless windowed notes  =
|anisotropic                = unknown
|anisotropic notes          =
|antialiasing               = unknown
|antialiasing notes         =
|upscaling                  = unknown
|upscaling tech             =
|upscaling notes            =
|vsync                      = unknown
|vsync notes                =
|60 fps                     = unknown
|60 fps notes               =
|120 fps                    = unknown
|120 fps notes              =
|hdr                        = unknown
|hdr notes                  =
|ray tracing                = unknown
|ray tracing notes          =
|color blind                = unknown
|color blind notes          =
}}`)

	fmt.Println("* [18/25] Processing Input!")

	outputFile.WriteString("\n\n==Input==\n{{Input")

	controller := false
	if game.Data.ControllerSupport != nil {
		controller = true
	}

	outputFile.WriteString(`
|key remap                 = unknown
|key remap notes           =
|acceleration option       = unknown
|acceleration option notes =
|mouse sensitivity         = unknown
|mouse sensitivity notes   =
|mouse menu                = unknown
|mouse menu notes          =
|invert mouse y-axis       = unknown
|invert mouse y-axis notes =
|touchscreen               = unknown
|touchscreen notes         = `)

	outputFile.WriteString(fmt.Sprintf("\n|controller support        = %v\n|controller support notes  = \n|full controller           = ", controller))
	if controller && *game.Data.ControllerSupport == "full" {
		outputFile.WriteString("true")
	} else {
		outputFile.WriteString("false")
	}
	outputFile.WriteString("\n|full controller notes     = ")

	outputFile.WriteString(`
|controller remap          = unknown
|controller remap notes    =
|controller sensitivity    = unknown
|controller sensitivity notes=
|invert controller y-axis  = unknown
|invert controller y-axis notes=
|xinput controllers        = unknown
|xinput controllers notes  =
|xbox prompts              = unknown
|xbox prompts notes        =
|impulse triggers          = unknown
|impulse triggers notes    =
|dualshock 4               = unknown
|dualshock 4 notes         =
|dualshock prompts         = unknown
|dualshock prompts notes   =
|light bar support         = unknown
|light bar support notes   =
|dualshock 4 modes         = unknown
|dualshock 4 modes notes   =
|tracked motion controllers= unknown
|tracked motion controllers notes =
|tracked motion prompts    = unknown
|tracked motion prompts notes =
|other controllers         = unknown
|other controllers notes   =
|other button prompts      = unknown
|other button prompts notes=
|controller hotplug        = unknown
|controller hotplug notes  =
|haptic feedback           = unknown
|haptic feedback notes     =
|simultaneous input        = unknown
|simultaneous input notes  =
|steam input api           = unknown
|steam input api notes     =
|steam hook input          = unknown
|steam hook input notes    =
|steam input presets       = unknown
|steam input presets notes =
|steam controller prompts  = unknown
|steam controller prompts notes =
|steam cursor detection    = unknown
|steam cursor detection notes =
}}`)

	fmt.Println("* [19/25] Processing Audio!")

	game.ProcessLanguages()

	outputFile.WriteString("\n\n")
	outputFile.WriteString(`==Audio==
{{Audio
|separate volume           = unknown
|separate volume notes     =
|surround sound            = unknown
|surround sound notes      = `)
	outputFile.WriteString(fmt.Sprintf("\n|subtitles                 = %v\n", game.Data.Subtitles))

	outputFile.WriteString(`|subtitles notes           =
|closed captions           = unknown
|closed captions notes     =
|mute on focus lost        = unknown
|mute on focus lost notes  =
|eax support               =
|eax support notes         =
|royalty free audio        = unknown
|royalty free audio notes  =
|red book cd audio         =
|red book cd audio notes   =
|general midi audio        =
|general midi audio notes  =
}}`)

	fmt.Println("* [20/25] Processing Languages!")

	outputFile.WriteString("\n\n{{L10n|content=")

	orderedLangauges := make([]string, 0, len(game.Data.Languages))
	for key := range game.Data.Languages {
		sanitisedKey := key
		orderedLangauges = append(orderedLangauges, sanitisedKey)
	}

	sort.Strings(orderedLangauges)

	// find English and swap it to be the first language instead...

	if orderedLangauges[0] != "English" {
		// Only swap English to the first language if it isn't already...

		foundIndex := 0
		for i := 1; i < len(orderedLangauges); i++ {
			if orderedLangauges[i] == "English" {
				foundIndex = i
				break
			}
		}

		if foundIndex != 0 {
			for i := foundIndex; i > 0; i-- {
				orderedLangauges[i] = orderedLangauges[i-1]
			}
			orderedLangauges[0] = "English"
		}
	}

	for _, key := range orderedLangauges {
		outputFile.WriteString(game.FormatLanguage(key))
	}

	outputFile.WriteString("\n}}\n")

	fmt.Println("* [21/25] Processing Network!")

	if game.HasCategory(Multiplayer) {
		outputFile.WriteString("\n\n==Network==")
		outputFile.WriteString("\n{{Network/Multiplayer")
		outputFile.WriteString("\n|local play           = ")
		if game.HasCategory(LocalMultiPlayer) || game.HasCategory(LocalCoOp) {
			outputFile.WriteString("true")
		} else {
			outputFile.WriteString("false")
		}
		outputFile.WriteString(`
|local play players   =
|local play modes     =
|local play notes     = `)

		outputFile.WriteString("\n|lan play             = ")
		if game.HasCategory(CoOp) {
			outputFile.WriteString("true")
		} else {
			outputFile.WriteString("false")
		}
		outputFile.WriteString(`
|lan play players     =
|lan play modes       =
|lan play notes       = `)

		outputFile.WriteString("\n|online play          = ")
		if game.HasCategory(OnlineMultiPlayer) || game.HasCategory(OnlineCoOp) {
			outputFile.WriteString("true")
		} else {
			outputFile.WriteString("false")
		}
		outputFile.WriteString(`
|online play players  =
|online play modes    =
|online play notes    =
|asynchronous         =
|asynchronous notes   =
}}`)
		outputFile.WriteString("\n{{Network/Connections")
		outputFile.WriteString(`
|matchmaking        =
|matchmaking notes  =
|p2p                =
|p2p notes          =
|dedicated          =
|dedicated notes    =
|self-hosting       =
|self-hosting notes =
|direct ip          =
|direct ip notes    =
}}{{Network/Ports
|tcp  =
|udp  =
|upnp =
}}`)
	}

	fmt.Println("* [22/25] Processing API!")

	outputFile.WriteString("\n\n==Other information==\n===API===\n{{API\n")
	outputFile.WriteString(fmt.Sprintf("|direct3d versions      = %s\n", game.FindDirectX()))
	outputFile.WriteString(fmt.Sprintf(`|direct3d notes         =
|directdraw versions    =
|directdraw notes       =
|wing                   =
|wing notes             =
|opengl versions        =
|opengl notes           =
|glide versions         =
|glide notes            =
|software mode          =
|software mode notes    =
|mantle support         =
|mantle support notes   =
|metal support          =
|metal support notes    =
|vulkan versions        =
|vulkan notes           =
|dos modes              =
|dos modes notes        =
|windows 32-bit exe     = %s
|windows 64-bit exe     = %s
|windows arm app        = false
|windows exe notes      =
|mac os x powerpc app   = false
|macos intel 32-bit app = %s
|macos intel 64-bit app = %s
|macos arm app          = unknown
|macos app notes        =
|linux powerpc app      = false
|linux 32-bit executable= %s
|linux 64-bit executable= %s
|linux arm app          = false
|linux 68k app          = false
|linux executable notes =
|mac os powerpc app     = false
|mac os 68k app         = false
|mac os executable notes=
}}`,
		GetExeBit(true, "windows", game.Data.Platforms, game.Data.PCRequirements), GetExeBit(false, "windows", game.Data.Platforms, game.Data.PCRequirements),
		GetExeBit(true, "mac", game.Data.Platforms, game.Data.MACRequirements), GetExeBit(false, "mac", game.Data.Platforms, game.Data.MACRequirements),
		GetExeBit(true, "linux", game.Data.Platforms, game.Data.LinuxRequirements), GetExeBit(false, "linux", game.Data.Platforms, game.Data.LinuxRequirements)))

	fmt.Println("* [23/25] Processing Middleware!")

	outputFile.WriteString("\n\n===Middleware===\n{{Middleware")
	outputFile.WriteString(`
|physics          =
|physics notes    =
|audio            =
|audio notes      =
|interface        =
|interface notes  =
|input            =
|input notes      =
|cutscenes        =
|cutscenes notes  =
|multiplayer      =
|multiplayer notes=
|anticheat        =
|anticheat notes  =
}}`)

	fmt.Println("* [24/25] Processing System Requirements!")
	outputFile.WriteString("\n\n==System requirements==")

	outputFile.WriteString(game.OutputSpecs())

	fmt.Println("* [25/25] Processing References!")
	outputFile.WriteString("\n{{References}}")

	outputFile.Close()

	println(fmt.Sprintf("Successfully parsed information for game: '%s'", SanitiseName(game.Data.Name, true)))
}
