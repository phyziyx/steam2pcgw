package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
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

func ParseGame(gameId string) (body []byte, err error) {
	fileName := fmt.Sprintf("%s.json", gameId)
	fi, err := os.Stat(fileName)

	if err != nil || time.Since(fi.ModTime()).Hours() > (7*24) {
		fmt.Println("Did not find game cache or cache is older than 7 days...")

		var response *http.Response
		response, err = http.Get(fmt.Sprintf("%s%s%s", API_LINK, gameId, LOCALE))

		if err != nil {
			fmt.Printf("Failed to connect to the Steam API... (error: %s)\n", err)
			return
		}

		if response.StatusCode != http.StatusOK {
			fmt.Printf("Failed to connect to the Steam API... (HTTP code: %d)\n", response.StatusCode)
			return
		}
		defer response.Body.Close()

		body, err = io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("An error occurred while attempting to parse the response body...")
			return
		}

		var cache *os.File
		cache, err = os.Create(fmt.Sprintf("%s.json", gameId))
		if err != nil {
			fmt.Println("Failed to cache the game... Process continuing...")
		} else {
			cache.WriteString(string(body))
			fmt.Println("Cached the game...")
		}
		cache.Close()
	} else {
		body, err = os.ReadFile(fileName)
	}

	return body, err
}

func TakeInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	// For Windows and Linux
	text = strings.TrimSuffix(text, "\n")
	// For Windows
	text = strings.TrimSuffix(text, "\r")

	fmt.Printf("Sanitised text: '%s' (len: %d)\n", text, len(text))

	if len(text) == 0 {
		return "", errors.New("Invalid input")
	}

	return text, nil
}

func isBlacklistedGenre(genre string) bool {
	blacklisted := []string{"Early Access", "Indie", "Casual"}
	for _, listedGenre := range blacklisted {
		if genre == listedGenre {
			return true
		}
	}

	return false
}

func OutputGenres(genres []Genre) string {
	var output = ""

	for _, v := range genres {
		if isBlacklistedGenre(v.Description) {
			continue
		}

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

func FindDirectX(pcRequirements Requirement) string {
	if len(pcRequirements["minimum"].(string)) == 0 {
		return ""
	}

	sanitised := removeTags(pcRequirements["minimum"].(string))
	dxRegex := regexp.MustCompile(`DirectX:(.+)\n`)
	version := dxRegex.FindStringSubmatch(sanitised)
	if len(version) == 2 {
		return strings.Trim(version[1], "Version ")
	}

	return ""
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
	output = strings.ReplaceAll(output, "available space", "")
	output = strings.ReplaceAll(output, "RAM", "")
	output = strings.ReplaceAll(output, "Version ", "")
	output = strings.ReplaceAll(output, "Windows ", "")

	networkRe := regexp.MustCompile(`Network:(.+)\n`)
	output = networkRe.ReplaceAllLiteralString(output, "")

	// Determine
	if isMin {
		level = "min"
		output = strings.Replace(output, "Minimum:\n", "", 1)
	} else {
		level = "rec"
		output = strings.Replace(output, "Recommended:\n", "", 1)
	}

	// Replace
	output = strings.Replace(output, "OS:", fmt.Sprintf("|%sOS    = ", level), 1)

	// Processor stuff
	if strings.Contains(output, "Processor:") {
		cpuRegEx := regexp.MustCompile(`(Processor:)(.+?)(?: or |/|,|\|)+(.+)\n`)
		cpus := cpuRegEx.FindStringSubmatch(output)

		if len(cpus) == 4 {
			output = cpuRegEx.ReplaceAllLiteralString(output, fmt.Sprintf("|%sCPU   = %s\n|%sCPU2  = %s\n", level, cpus[2], level, strings.TrimPrefix(cpus[3], " ")))
		} else {
			cpuRegEx = regexp.MustCompile(`Processor:(.+)\n`)
			cpus = cpuRegEx.FindStringSubmatch(output)
			output = cpuRegEx.ReplaceAllLiteralString(output, fmt.Sprintf("|%sCPU   = %s\n|%sCPU2  = %s\n", level, cpus[1], level, cpus[1]))
		}
	}

	output = strings.TrimSuffix(strings.Replace(output, "Storage:", fmt.Sprintf("|%sHD    = ", level), 1), " ")

	// Graphics stuff
	if strings.Contains(output, "Graphics:") {
		gpuRegEx := regexp.MustCompile(`Graphics:(.+)\n`)
		gpus := gpuRegEx.FindStringSubmatch(output)
		if strings.Contains(gpus[0], "OpenGL") {
			output = gpuRegEx.ReplaceAllLiteralString(output, fmt.Sprintf("|%sOGL   = %s\n", level, strings.ReplaceAll(strings.ReplaceAll(gpus[1], " or greater", ""), "OpenGL ", "")))
		} else {
			// Did not find OpenGL stuff, this means we can do a different regex then...
			gpuRegEx2 := regexp.MustCompile(`(Graphics:)(.+)(?: or |/|,|\|)+(.+)\n`)
			gpus := gpuRegEx2.FindStringSubmatch(output)
			if len(gpus) == 4 {
				output = gpuRegEx2.ReplaceAllLiteralString(output, fmt.Sprintf("|%sGPU   = %s\n|%sGPU2  = %s\n", level, gpus[2], level, strings.TrimPrefix(gpus[3], " ")))
			} else {
				gpus := gpuRegEx.FindStringSubmatch(output)
				output = gpuRegEx.ReplaceAllLiteralString(output, fmt.Sprintf("|%sGPU   = %s\n|%sGPU2  = %s\n", level, gpus[1], level, gpus[1]))
			}
		}
	}

	output = strings.TrimSuffix(strings.Replace(output, "Memory:", fmt.Sprintf("|%sRAM   = ", level), 1), " ")
	output = strings.Replace(output, "OS:", fmt.Sprintf("|%sVRAM    = ", level), 1)
	output = strings.Replace(output, "DirectX:", fmt.Sprintf("|%sDX    = ", level), 1)
	output = strings.Replace(output, "Sound Card:", fmt.Sprintf("|%saudio = ", level), 1)

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

	if platforms.Windows {
		output += "\n{{System requirements\n"
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
		output += "\n}}\n"
	}

	if platforms.MAC {
		output += "\n{{System requirements\n"
		output += ("|OSfamily = OS X")
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

func ParseDate(date string) (output string) {
	output = date
	dateRe, _ := regexp.Compile(`(\d+) (\w+), (\d+)`)
	tokens := dateRe.FindStringSubmatch(date)
	if len(tokens) != 0 {
		output = dateRe.ReplaceAllString(date, `$2 $1 $3`)
	}
	return output
}

// TODO:

func HasInAppPurchases(Categories []Category) bool {
	for _, v := range Categories {
		if v.ID == 35 {
			return true
		}
	}
	return false
}

func HasFullControllerSupport(Categories []Category) bool {
	for _, v := range Categories {
		if v.ID == 28 {
			return true
		}
	}
	return false
}

func HasMultiplayerSupport(Categories []Category) bool {
	for _, v := range Categories {
		if v.ID == 1 {
			return true
		}
	}

	return false
}

func IsEarlyAccess(genres []Genre) bool {
	for _, v := range genres {
		if v.Description == "Early Access" {
			return true
		}
	}
	return false
}

func HasSteamCloud(Categories []Category) bool {
	for _, v := range Categories {
		if v.ID == 23 {
			return true
		}
	}

	return false
}
