package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func UnmarshalGame(data []byte) (Game, error) {
	var r Game
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Game) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func ProcessSpecs(input string, isMin bool) string {
	// Create vars
	var level string
	output := input

	if len(output) == 0 {
		return output
	}

	// Sanitise input and remove HTML tags
	noTag, _ := regexp.Compile(`(<[^>]*>)+`)
	output = noTag.ReplaceAllLiteralString(output, "\n")

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
	output = strings.Replace(output, "OS:\n", fmt.Sprintf("|%sOS    = ", level), 1)

	output = strings.Replace(output, "Processor:\n", fmt.Sprintf("|%sCPU    = |%sCPU2    = ", level, level), 1)

	output = strings.Replace(output, "Storage:\n", fmt.Sprintf("|%sHD    = ", level), 1)

	output = strings.Replace(output, "Graphics:\n", fmt.Sprintf("|%sGPU    = |%sGPU2    = ", level, level), 1)
	output = strings.Replace(output, "Memory:\n", fmt.Sprintf("|%sRAM   = ", level), 1)
	output = strings.Replace(output, "OS:\n", fmt.Sprintf("|%sVRAM    = ", level), 1)
	output = strings.Replace(output, "DirectX:\n", fmt.Sprintf("|%sDX    = ", level), 1)

	output = strings.Replace(output, "Additional Notes:\n", "\n|notes    = ", 1)

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

func OutputSpecs(platforms Platforms, pcRequirements, macRequirements, linuxRequirements Requirements) string {
	var output string = ""
	var specs string = ""

	if platforms.Windows {
		output += "|OSfamily = Windows"
		specs = ProcessSpecs(pcRequirements.Minimum, true)
		output += (specs)

		// Handle recommended specs
		if pcRequirements.Recommended != nil {
			specs = ProcessSpecs(*pcRequirements.Recommended, false)
			output += (specs)
		} else {
			emptySpecs("rec")
		}

		// Output closure
		output += ("\n}}")
	}

	if platforms.MAC {
		output += ("|OSfamily = Mac")
		specs = ProcessSpecs(macRequirements.Minimum, true)
		output += (specs)

		// Handle recommended specs
		if macRequirements.Recommended != nil {
			specs = ProcessSpecs(*macRequirements.Recommended, false)
			output += (specs)
		} else {
			emptySpecs("rec")
		}

		// Output closure
		output += ("\n}}")
	}

	if platforms.Linux {
		output += ("|OSfamily = Linux")
		specs = ProcessSpecs(linuxRequirements.Minimum, true)
		output += (specs)

		// Handle recommended specs
		if linuxRequirements.Recommended != nil {
			specs = ProcessSpecs(*linuxRequirements.Recommended, false)
			output += (specs)
		} else {
			emptySpecs("rec")
		}

		// Output closure
		output += ("\n}}")
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
