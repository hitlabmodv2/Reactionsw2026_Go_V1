package commands

import (
	"hisoka/src/helpers"
	"hisoka/src/libs"
	"hisoka/src/settings"
)

func init() {
	libs.NewCommands(&libs.ICommand{
		Name:     `mode`,
		As:       []string{"mode"},
		Tags:     "owner",
		IsPrefix: true,
		IsOwner:  true,
		Execute: func(conn *libs.IClient, m *libs.IMessage) bool {
			newValue := !settings.Public

			if newValue {
				m.Reply("The bot is now in public mode.")
			} else {
				m.Reply("The bot is now in private mode.")
			}

			// Update the in-memory value so the change takes effect
			// immediately for the running process…
			settings.Public = newValue

			// …and persist it to settings.go so it survives a restart.
			if err := helpers.UpdateSettingsBool("Public", newValue); err != nil {
				m.Reply("Failed to update settings file: " + err.Error())
				return false
			}

			return true
		},
	})
}
