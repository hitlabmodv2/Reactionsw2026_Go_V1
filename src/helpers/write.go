package helpers

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

// settingsFilePath is the path to the Go file that holds user-editable
// configuration (replaces the old .env file).
const settingsFilePath = "src/settings/settings.go"

// UpdateSettingsBool rewrites a top-level boolean setting in
// `src/settings/settings.go`. It looks for a line like:
//
//	Name = true
//	Name = false
//
// and replaces the literal with the new value, leaving everything else
// (indentation, comments) untouched. Returns an error if the file can't
// be read/written or if the named setting can't be found.
func UpdateSettingsBool(name string, value bool) error {
	data, err := os.ReadFile(settingsFilePath)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`(?m)^(\s*` + regexp.QuoteMeta(name) + `\s*=\s*)(true|false)\b`)
	if !re.Match(data) {
		return fmt.Errorf("setting %q not found in %s", name, settingsFilePath)
	}

	updated := re.ReplaceAllString(string(data), "${1}"+strconv.FormatBool(value))
	return os.WriteFile(settingsFilePath, []byte(updated), 0644)
}
