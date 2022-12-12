package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {

	if len(os.Args) > 1 && strings.ToLower(os.Args[1]) == "-v" {
		fmt.Println(APP_NAME, VERSION, "(", GH_LINK, ")")
		return
	}

	fmt.Println("Running", APP_NAME, VERSION, "(", GH_LINK, ")")

	var gameId string
	var err error = nil

	// Ask for input from the user
	for len(gameId) == 0 || err != nil {
		print("Insert the Steam app ID: ")
		gameId, err = TakeInput()
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("Fetching game app details...")

	var gameJson []byte
	gameJson, err = ParseGame(gameId)
	if err != nil {
		fmt.Println(err)
		return
	}

	var game Game
	game, err = UnmarshalGame(gameJson)

	if err != nil {
		fmt.Printf("An error occurred while attempting to unmarshal the JSON... (%s)\n", err)
		return
	}

	if !game[gameId].Success {
		fmt.Println("The app ID provided does not exist on the Steam database...")
		return
	}

	outputFile, err := os.Create(fmt.Sprintf("%s.txt", gameId))
	if err != nil {
		fmt.Println("Failed to create the output file... Process stopped!")
		return
	}
	defer outputFile.Close()

	fmt.Println("* [1/24] Adding stub")
	outputFile.WriteString("{{stub}}\n")

	fmt.Println("* [2/24] Adding app cover")
	outputFile.WriteString(fmt.Sprintf("{{Infobox game\n|cover        = %s cover.jpg", SanitiseName(game[gameId].Data.Name, true)))

	fmt.Println("* [3/24] Adding app developers")
	outputFile.WriteString("\n|developers   = ")
	for _, developer := range game[gameId].Data.Developers {
		outputFile.WriteString(fmt.Sprintf("\n{{Infobox game/row/developer|%s}}", SanitiseName(developer, false)))
	}

	fmt.Println("* [4/24] Adding app publishers")
	outputFile.WriteString("\n|publishers   = ")
	for _, publisher := range game[gameId].Data.Publishers {
		if len(game[gameId].Data.Publishers) == 1 {
			skip := false
			for _, developer := range game[gameId].Data.Developers {
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

	fmt.Println("* [5/24] Adding app release date")
	outputFile.WriteString("\n|engines      =\n{{Infobox game/row/engine|}}\n|release dates= \n")

	var date string = ""
	if IsEarlyAccess(game[gameId].Data.Genres) {
		date += "EA"
	} else if game[gameId].Data.ReleaseDate.ComingSoon {
		if success, _ := IsDate(game[gameId].Data.ReleaseDate.Date); success {
			date += ParseDate(game[gameId].Data.ReleaseDate.Date)
		} else {
			date += "TBA"
		}
	} else {
		date += ParseDate(game[gameId].Data.ReleaseDate.Date)
	}

	if game[gameId].Data.Platforms.Windows {
		outputFile.WriteString(fmt.Sprintf("{{Infobox game/row/date|Windows| %s }}\n", date))
	}

	if game[gameId].Data.Platforms.MAC {
		outputFile.WriteString(fmt.Sprintf("{{Infobox game/row/date|OS X| %s }}\n", date))
	}

	if game[gameId].Data.Platforms.Linux {
		outputFile.WriteString(fmt.Sprintf("{{Infobox game/row/date|Linux| %s }}\n", date))
	}

	fmt.Println("* [6/24] Adding reception score")
	if game[gameId].Data.Metacritic != nil {
		fmt.Println("* [6/24] Added Metacritic")
		outputFile.WriteString(fmt.Sprintf("\n|reception    = \n{{Infobox game/row/reception|Metacritic|%s|%d}}", game[gameId].Data.Metacritic.URL, game[gameId].Data.Metacritic.Score))
	} else {
		fmt.Println("* [6/24] Skipped Metacritic")
		outputFile.WriteString(("\n|reception    = \n{{Infobox game/row/reception|Metacritic|link|rating}}"))
	}

	outputFile.WriteString(("\n{{Infobox game/row/reception|OpenCritic|link|rating}}\n{{Infobox game/row/reception|IGDB|link|rating}}"))

	outputFile.WriteString("\n|taxonomy     =\n{{Infobox game/row/taxonomy/monetization      | ")
	if game[gameId].Data.IsFree {
		fmt.Println("* [7/24] Game is F2P")
		outputFile.WriteString(("Free-to-play }}"))
	} else {
		fmt.Println("* [7/24] Game is not F2P")
		outputFile.WriteString(("One-time game purchase }}"))
	}

	// TODO:
	fmt.Println("* [8/24] Taxonomy...")
	outputFile.WriteString("\n{{Infobox game/row/taxonomy/microtransactions | }}\n{{Infobox game/row/taxonomy/modes             | ")

	modes := ""

	if HasSinglePlayerSupport(game[gameId].Data.Categories) {
		modes += "Singleplayer, "
	}

	if HasMultiplayerSupport(game[gameId].Data.Categories) {
		modes += "Multiplayer, "
	}

	modes = strings.TrimSuffix(modes, ", ")
	outputFile.WriteString(modes)

	outputFile.WriteString(" }}\n{{Infobox game/row/taxonomy/pacing            | }}\n{{Infobox game/row/taxonomy/perspectives      | }}\n{{Infobox game/row/taxonomy/controls          | }}\n{{Infobox game/row/taxonomy/genres            | ")
	outputFile.WriteString(OutputGenres(game[gameId].Data.Genres))
	outputFile.WriteString(" }}\n{{Infobox game/row/taxonomy/sports            | }}\n{{Infobox game/row/taxonomy/vehicles          | }}\n{{Infobox game/row/taxonomy/art styles        | }}\n{{Infobox game/row/taxonomy/themes            | }}\n{{Infobox game/row/taxonomy/series            |  }}\n")

	outputFile.WriteString(fmt.Sprintf("|steam appid  = %s\n|steam appid side = \n", gameId))
	if game[gameId].Data.Dlc != nil {
		var dlcs string = ""
		for _, v := range game[gameId].Data.Dlc {
			dlcs += fmt.Sprintf("%v, ", v)
		}
		dlcs = strings.TrimSuffix(dlcs, ", ")
		dlcs += "\n"
		outputFile.WriteString(dlcs)
	}
	outputFile.WriteString("|gogcom id    = \n|gogcom id side = \n|official site= ")

	if game[gameId].Data.Website != nil {
		outputFile.WriteString(*game[gameId].Data.Website)
	} else {
		outputFile.WriteString(game[gameId].Data.SupportInfo.URL)
	}
	outputFile.WriteString("\n|hltb         = \n|igdb         = <!-- Only needs to be set if there is no IGDB reception row -->\n|lutris       = \n|mobygames    = \n|strategywiki = \n|wikipedia    = \n|winehq       = \n|license      = commercial\n}}")

	fmt.Println("* [9/24] Processing introduction...")
	outputFile.WriteString("\n\n{{Introduction\n|introduction      = ")
	// outputFile.WriteString(removeTags(game[gameId].Data.AboutTheGame))

	outputFile.WriteString("\n\n|release history   = ")

	outputFile.WriteString("\n\n|current state     = ")
	outputFile.WriteString("\n}}")

	outputFile.WriteString("\n\n'''General information'''")
	outputFile.WriteString("\n{{mm}} [https://steamcommunity.com/app/" + gameId + "/discussions/ Steam Community Discussions]")

	fmt.Println("* [10/24] Processing Availability!")

	outputFile.WriteString("\n\n==Availability==\n{{Availability|\n")

	platforms := ""

	if game[gameId].Data.Platforms.Windows {
		platforms += "Windows, "
	}
	if game[gameId].Data.Platforms.MAC {
		platforms += "OS X, "
	}
	if game[gameId].Data.Platforms.Linux {
		platforms += "Linux, "
	}

	platforms = strings.TrimSuffix(platforms, ", ")

	outputFile.WriteString(fmt.Sprintf("{{Availability/row| Steam | %s | Steam |  |  | %s }}", gameId, platforms))
	outputFile.WriteString("\n}}")

	outputFile.WriteString("\n\n<!-- PAGE GENERATED BY STEAM2PCGW -->\n")

	fmt.Println("* [11/24] Processing Monetization!")
	outputFile.WriteString("\n\n==Monetization==\n")

	outputFile.WriteString("{{Monetization")
	outputFile.WriteString("\n|ad-supported           = ")
	outputFile.WriteString("\n|dlc                    = ")
	outputFile.WriteString("\n|expansion pack         = ")
	outputFile.WriteString("\n|freeware               = ")
	outputFile.WriteString("\n|free-to-play           = ")
	if game[gameId].Data.IsFree {
		outputFile.WriteString("The game has such monetization.")
	}
	outputFile.WriteString("\n|one-time game purchase = ")
	if !game[gameId].Data.IsFree {
		outputFile.WriteString("The game requires an upfront purchase to access.")
	}
	outputFile.WriteString("\n|sponsored              = ")
	outputFile.WriteString("\n|subscription           = ")
	outputFile.WriteString("\n}}")

	fmt.Println("* [12/24] Processing Microtransactions!")

	outputFile.WriteString("\n\n===Microtransactions===\n{{Microtransactions")

	outputFile.WriteString("\n|boost               = ")
	outputFile.WriteString("\n|cosmetic            = ")
	outputFile.WriteString("\n|currency            = ")
	outputFile.WriteString("\n|finite spend        = ")
	outputFile.WriteString("\n|infinite spend      = ")
	outputFile.WriteString("\n|free-to-grind       = ")
	outputFile.WriteString("\n|loot box            = ")
	outputFile.WriteString("\n|none                = None")
	outputFile.WriteString("\n|player trading      = ")
	outputFile.WriteString("\n|time-limited        = ")
	outputFile.WriteString("\n|unlock              = ")
	outputFile.WriteString("\n}}")

	fmt.Println("* [13/24] Processing DLCs!")
	outputFile.WriteString("\n\n{{DLC|\n<!-- DLC rows goes below: -->\n}}")

	fmt.Println("* [14/24] Processing Config File Location!")

	outputFile.WriteString("\n\n==Game data==\n===Configuration file(s) location===")
	outputFile.WriteString("\n{{Game data|")
	if game[gameId].Data.Platforms.Windows {
		outputFile.WriteString("\n{{Game data/config|Windows|}}")
	}
	if game[gameId].Data.Platforms.MAC {
		outputFile.WriteString("\n{{Game data/config|OS X|}}")
	}
	if game[gameId].Data.Platforms.Linux {
		outputFile.WriteString("\n{{Game data/config|Linux|}}")
	}
	outputFile.WriteString("\n}}")

	fmt.Println("* [15/24] Processing Save Game Location!")

	outputFile.WriteString("\n\n===Save game data location===")
	outputFile.WriteString("\n{{Game data|")
	if game[gameId].Data.Platforms.Windows {
		outputFile.WriteString("\n{{Game data/saves|Windows|}}")
	}
	if game[gameId].Data.Platforms.MAC {
		outputFile.WriteString("\n{{Game data/saves|OS X|}}")
	}
	if game[gameId].Data.Platforms.Linux {
		outputFile.WriteString("\n{{Game data/saves|Linux|}}")
	}
	outputFile.WriteString("\n}}")

	fmt.Println("* [16/24] Processing Save Game Sync!")

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
	if HasSteamCloud(game[gameId].Data.Categories) {
		outputFile.WriteString("true")
	} else {
		if game[gameId].Data.ReleaseDate.ComingSoon {
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
	fmt.Println("* [17/24] Processing Video!")
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

	fmt.Println("* [18/24] Processing Input!")

	outputFile.WriteString("\n\n==Input==\n{{Input\n")

	controller := false
	if game[gameId].Data.ControllerSupport != nil {
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
	if controller && *game[gameId].Data.ControllerSupport == "full" {
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

	fmt.Println("* [19/24] Processing Audio!")

	outputFile.WriteString("\n\n")
	outputFile.WriteString(`==Audio==
{{Audio
|separate volume           = unknown
|separate volume notes     = 
|surround sound            = unknown
|surround sound notes      = 
|subtitles                 = unknown
|subtitles notes           = 
|closed captions           = unknown
|closed captions notes     = 
|mute on focus lost        = unknown
|mute on focus lost notes  = 
|royalty free audio        = unknown
|royalty free audio notes  = 
|eax support               = 
|eax support notes         = 
|red book cd audio         =
|red book cd audio notes   = 
|general midi audio        = 
|general midi audio notes  = 
}}`)

	fmt.Println("* [20/24] Processing Languages!")
	languages := ProcessLanguages(game[gameId].Data.SupportedLanguages)

	outputFile.WriteString("\n\n{{L10n|content=")

	orderedLangauges := make([]string, 0, len(languages))
	for key := range languages {
		orderedLangauges = append(orderedLangauges, key)
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
			for i := 0; i < foundIndex-1; i++ {
				orderedLangauges[i+1] = orderedLangauges[i]
			}
			orderedLangauges[0] = "English"
		}
	}

	for _, key := range orderedLangauges {
		outputFile.WriteString(FormatLanguage(key, languages))
	}

	outputFile.WriteString("\n}}\n")

	fmt.Println("* [21/24] Processing API!")

	outputFile.WriteString("\n\n==Other information==\n===API===\n{{API\n")
	outputFile.WriteString(fmt.Sprintf("|direct3d versions      = %s\n", FindDirectX(game[gameId].Data.PCRequirements)))
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
|linux executable notes = 
}}`,
		GetExeBit(true, "windows", game[gameId].Data.Platforms, game[gameId].Data.PCRequirements), GetExeBit(false, "windows", game[gameId].Data.Platforms, game[gameId].Data.PCRequirements),
		GetExeBit(true, "mac", game[gameId].Data.Platforms, game[gameId].Data.MACRequirements), GetExeBit(false, "mac", game[gameId].Data.Platforms, game[gameId].Data.MACRequirements),
		GetExeBit(true, "linux", game[gameId].Data.Platforms, game[gameId].Data.LinuxRequirements), GetExeBit(false, "linux", game[gameId].Data.Platforms, game[gameId].Data.LinuxRequirements)))

	fmt.Println("* [22/24] Processing Middleware!")

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

	if HasMultiplayerSupport(game[gameId].Data.Categories) {
		outputFile.WriteString("\n\n===Network===\n{{Network/Multiplayer")
		outputFile.WriteString(`
|local play           = 
|local play players   = 
|local play modes     = 
|local play notes     = 
|lan play             = 
|lan play players     = 
|lan play modes       = 
|lan play notes       = 
|online play          = 
|online play players  = 
|online play modes    = 
|online play notes    = 
|asynchronous         = 
|asynchronous notes   = 
}}`)
	}

	fmt.Println("* [23/24] Processing System Requirements!")
	outputFile.WriteString("\n\n==System requirements==\n")

	outputFile.WriteString(OutputSpecs(game[gameId].Data.Platforms, game[gameId].Data.PCRequirements, game[gameId].Data.MACRequirements, game[gameId].Data.LinuxRequirements))

	fmt.Println("* [24/24] Processing References!")
	outputFile.WriteString("\n{{References}}")

	println(fmt.Sprintf("Successfully parsed information for game: '%s'", SanitiseName(game[gameId].Data.Name, true)))
}
