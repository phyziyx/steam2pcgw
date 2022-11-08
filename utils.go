package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func (req *Requirement) UnmarshalJSON(data []byte) error {
	if string(data) == `""` || string(data) == `{}` || string(data) == `[]` {
		return nil
	}

	type requirement Requirement
	return json.Unmarshal(data, (*requirement)(req))
}

func UnmarshalGame(data []byte) (Game, error) {
	var r Game
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Game) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func OutputGenres(genres []Genre) string {
	var output = ""

	for _, v := range genres {
		output += v.Description + ", "
	}
	output = strings.TrimSuffix(output, ", ")

	return output
}

func GetExeBit(is32 bool, platform string, platforms Platforms, requirements Requirement) string {
	value := "unknown"

	if (platform == "windows" && !platforms.Windows) ||
		(platform == "mac" && !platforms.MAC) ||
		(platform == "linux" && !platforms.Linux) {
	} else {
		var sanitised = strings.ToLower(requirements["minimum"].(string))
		sanitised = removeTags(sanitised)

		// Could have just used RAM but hey /shrug/
		ramFinder, _ := regexp.Compile(`memory:(\d+) gb`)
		ramFound := ramFinder.FindStringSubmatch(sanitised)
		var ram = 0
		if len(ramFound) != 0 {
			ram, _ = strconv.Atoi(ramFound[1])
		}

		// This may need to be modified!
		if is32 && (strings.Contains(sanitised, "64-bit") || strings.Contains(sanitised, "64 bit") || ram > 4) {
			value = "false"
		} else {
			value = "true"
		}
	}

	fmt.Printf("* [21/24] %s (32-bit: %v): %s\n", platform, is32, value)

	return value
}

func removeTags(input string) string {
	noTag, _ := regexp.Compile(`(<[^>]*>)+`)
	output := noTag.ReplaceAllLiteralString(input, "\n")
	output = strings.ReplaceAll(output, "\n ", "")
	return output
}

func ProcessSpecs(input string, isMin bool) string {
	// Create vars
	var level string
	output := input

	if len(output) == 0 {
		return output
	}

	// Sanitise input and remove HTML tags
	output = removeTags(output)

	// Cleanup some text, more texts must be added here...
	output = strings.Replace(output, "Requires a 64-bit processor and operating system", "", 1)

	// Determine
	if isMin {
		level = "min"
		output = strings.Replace(output, "Minimum:", "", 1)
	} else {
		level = "rec"
		output = strings.Replace(output, "Recommended:", "", 1)
	}

	// Replace
	output = strings.Replace(output, "OS:", fmt.Sprintf("|%sOS    = ", level), 1)

	output = strings.Replace(output, "Processor:", fmt.Sprintf("|%sCPU    = |%sCPU2    = ", level, level), 1)

	output = strings.Replace(output, "Storage:", fmt.Sprintf("|%sHD    = ", level), 1)

	output = strings.Replace(output, "Graphics:", fmt.Sprintf("|%sGPU    = |%sGPU2    = ", level, level), 1)
	output = strings.Replace(output, "Memory:", fmt.Sprintf("|%sRAM   = ", level), 1)
	output = strings.Replace(output, "OS:", fmt.Sprintf("|%sVRAM    = ", level), 1)
	output = strings.Replace(output, "DirectX:", fmt.Sprintf("|%sDX    = ", level), 1)

	output = strings.Replace(output, "Additional Notes:", "\n|notes    = ", 1)

	// Output
	return output
}

func emptySpecs(level string) string {
	return fmt.Sprintf(`|%sOS    = 
|%sCPU   = 
|%sCPU2  = 
|%sRAM   = 
|%sHD    = 
|%sGPU   = 
|%sGPU2  = 
|%sVRAM  = 
`, level, level, level, level, level, level, level, level)
}

func OutputSpecs(platforms Platforms, pcRequirements, macRequirements, linuxRequirements Requirement) string {
	var output string = ""
	var specs string = ""

	output += "\n{{System requirements\n"
	if platforms.Windows {
		output += "|OSfamily = Windows"
		specs = ProcessSpecs(pcRequirements["minimum"].(string), true)
		output += (specs)

		// Handle recommended specs
		if pcRequirements["recommended"] != nil {
			specs = ProcessSpecs(pcRequirements["recommended"].(string), false)
			output += (specs)
		} else {
			output += emptySpecs("rec")
		}
	}
	output += "\n}}\n"

	if platforms.MAC {
		output += "\n{{System requirements\n"
		output += ("|OSfamily = Mac")
		specs = ProcessSpecs(macRequirements["minimum"].(string), true)
		output += (specs)

		// Handle recommended specs
		if macRequirements["recommended"] != nil {
			specs = ProcessSpecs(macRequirements["recommended"].(string), false)
			output += (specs)
		} else {
			output += emptySpecs("rec")
		}
		output += "\n}}\n"
	}

	if platforms.Linux {
		output += "\n{{System requirements\n"
		output += ("|OSfamily = Linux")
		specs = ProcessSpecs(linuxRequirements["minimum"].(string), true)
		output += (specs)

		// Handle recommended specs
		if linuxRequirements["recommended"] != nil {
			specs = ProcessSpecs(linuxRequirements["recommended"].(string), false)
			output += (specs)
		} else {
			output += emptySpecs("rec")
		}
		output += "\n}}\n"
	}

	return output
}

func ProcessLanguages(input string) Language {
	languages := make(Language)
	var language string

	input = strings.Replace(input, "<br><strong>*</strong>languages with full audio support", "", 1)
	input = strings.ReplaceAll(input, ", ", "\n")
	input = strings.ReplaceAll(input, "<strong>", "")
	input = strings.ReplaceAll(input, "</strong>", "")

	for i := 0; i < len(input); i++ {

		// fmt.Printf("[ProcessLanguages] '%c' char found (language: '%s')\n", input[i], language)

		if rune(input[i]) == '\n' {
			// New line, new language!

			if len(language) != 0 {
				languages[language] = LanguageValue{
					UI:        true,
					Audio:     false,
					Subtitles: true,
				}
				// fmt.Printf("[ProcessLanguages] %s added (\\n found)\n", language)
			}

			language = ""
			continue
		}

		// Found * this means that it has complete support
		if input[i] == '*' {
			languages[language] = LanguageValue{
				UI:        true,
				Audio:     true,
				Subtitles: true,
			}
			// fmt.Printf("[ProcessLanguages] %s added (* found)\n", language)

			language = ""
			continue
		}

		// Append that language string
		language += string(input[i])
	}

	if len(language) != 0 {
		languages[language] = LanguageValue{
			UI:        true,
			Audio:     false,
			Subtitles: true,
		}
	}

	return languages
}

// TODO:

func HasInAppPurchases(Categories []Category) bool {
	for k := range Categories {
		if k == 35 {
			return true
		}
	}
	return false
}

func HasFullControllerSupport(Categories []Category) bool {
	for k := range Categories {
		if k == 28 {
			return true
		}
	}
	return false
}

func HasMultiplayerSupport(Categories []Category) bool {
	for k := range Categories {
		if k == 1 {
			return true
		}
	}

	return false
}
