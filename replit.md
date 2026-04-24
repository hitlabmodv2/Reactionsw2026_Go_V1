# hisoka — Go WhatsApp Bot

## Overview
A console (TUI) WhatsApp bot written in Go using the
[`go.mau.fi/whatsmeow`](https://pkg.go.dev/go.mau.fi/whatsmeow) library.
On startup it either prints a QR code to the terminal (default) or generates a
pairing code if `PAIRING_NUMBER` is set in `.env`. Session state is persisted
locally in `session.db` (SQLite via `mattn/go-sqlite3`, requires CGO).

This is **not** a web application — there is no HTTP server, frontend, or open
port. The user interacts with it through the workflow console (to scan the QR
code) and via WhatsApp itself.

## Project Layout
- `main.go` — entry point, loads `.env` and calls `conn.StartClient`.
- `src/hisoka.go` — initializes the SQLite store, connects the whatsmeow
  client, and prints QR / pairing code.
- `src/handlers/` — message and event handlers.
- `src/commands/` — bot commands (auto, main, owner).
- `src/helpers/` — logging, message parsing, write helpers.
- `src/libs/` — client wrapper, command registry, message types.

## Configuration (`.env`)
- `PUBLIC` — `true` to allow everyone, `false` for owner only.
- `OWNER` — comma-separated owner phone numbers (for `eval`/`exec`).
- `PAIRING_NUMBER` — phone number for pairing-code login (leave empty for QR).
- `REACT_STATUS` — `true` to react to every status update.
- `PREFIX` — regex of accepted command prefixes.

## Replit Setup
- **Language runtime:** Go 1.25 (auto-installed via `toolchain` directive in `go.mod`).
- **Build deps:** `gcc` (already provided by the Replit Nix environment) — required because `go-sqlite3` uses CGO.
- **Workflow:** `Start application` — `go run main.go`, output type `console`.
  It is a long-running TUI process; no port is exposed.

## Recent Changes (Replit import, 2026-04-24)
- Updated `go.mau.fi/whatsmeow` to the latest release because WhatsApp's
  servers were rejecting the older client version with `Client outdated (405)`.
- Adjusted three call sites for the new whatsmeow API that now requires a
  `context.Context` as the first argument:
  - `src/libs/client.go` — `WA.GetGroupInfo`
  - `src/handlers/message.go` — `SendPresence` and `WA.MarkRead`
  - `src/commands/Auto/readsw.go` — `WA.MarkRead`
- Ran `go mod tidy` to refresh `go.sum` and pull in transitive updates.

## Running
1. Open the `Start application` workflow output.
2. Scan the QR code shown in the console with WhatsApp → Linked Devices.
3. Once paired, `session.db` stores credentials, so subsequent restarts
   reconnect automatically.
