// Package settings holds all user-facing configuration for the bot.
//
// This file replaces the old `.env` file. Edit the values below to
// reconfigure the bot — no environment variables are read anymore.
//
// NOTE: The `mode` owner-command rewrites the Public value of this file
// at runtime via helpers.UpdateSettingsBool, so keep the format
// `Public = true` / `Public = false` on its own line.
package settings

var (
	// Public: true → everyone can use commands.
	//         false → only OWNER numbers can use commands.
	Public = false

	// Owner phone numbers authorized for owner-only commands
	// (eval, exec, mode). No "+" or spaces.
	Owner = []string{"6285815663170"}

	// PairingNumber: phone number used for pairing-code login.
	// Leave as "" to use QR-code login instead.
	PairingNumber = "6289667923162"

	// ReactStatus: true → automatically react to every WhatsApp status
	// with a random emoji. false → only auto-read, no reaction.
	ReactStatus = true

	// Prefix: regex of accepted command prefixes
	// (e.g. "[?!.#]" allows ?, !, ., #).
	Prefix = "[?!.#]"
)
