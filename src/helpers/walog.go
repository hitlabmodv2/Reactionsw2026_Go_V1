package helpers

import (
        "fmt"
        "strings"

        waLog "go.mau.fi/whatsmeow/util/log"
)

type prettyLog struct {
        module string
}

func PrettyWALogger(module string) waLog.Logger {
        return &prettyLog{module: module}
}

func (l *prettyLog) Sub(module string) waLog.Logger {
        if l.module != "" {
                module = l.module + "/" + module
        }
        return &prettyLog{module: module}
}

func (l *prettyLog) Debugf(msg string, args ...interface{}) {}

func (l *prettyLog) Infof(msg string, args ...interface{}) {
        text := fmt.Sprintf(msg, args...)
        if isNoisyInfo(text) {
                return
        }
        l.print(cBlue, "i", msg, args...)
}

func (l *prettyLog) Warnf(msg string, args ...interface{}) {
        text := fmt.Sprintf(msg, args...)
        if isNoisy(text) {
                return
        }
        l.print(cYellow, "!", msg, args...)
}

func (l *prettyLog) Errorf(msg string, args ...interface{}) {
        text := fmt.Sprintf(msg, args...)
        if isNoisy(text) {
                return
        }
        l.print(cRed, "✗", msg, args...)
}

func (l *prettyLog) print(color, sigil, msg string, args ...interface{}) {
        mod := l.module
        if mod == "" {
                mod = "wa"
        }
        fmt.Printf("  %s%s%s %s%s%s %s\n",
                color, sigil, cReset,
                cDim, mod, cReset,
                fmt.Sprintf(msg, args...),
        )
}

func isNoisyInfo(text string) bool {
        patterns := []string{
                "Stored ",
                "Uploading ",
                "Updating contact store",
                "push names from history sync",
                "history sync",
                "prekeys to server",
                "message secret keys",
                "Got 515 code",
                "Sending presence",
                "AppStateSyncComplete",
                "Successfully sent",
                "Migrated ",
                "identity keys",
                "sender keys",
        }
        for _, p := range patterns {
                if strings.Contains(text, p) {
                        return true
                }
        }
        return false
}

func isNoisy(text string) bool {
        patterns := []string{
                "failed to read frame header: EOF",
                "failed to get reader: failed to read frame header",
                "websocket: close 1006",
                "unexpected EOF",
                "Received stream end frame",
                "Disconnecting websocket",
                "frame header",
                "Error in websocket",
                "keepalive ping",
                "different participant list hash",
        }
        for _, p := range patterns {
                if strings.Contains(text, p) {
                        return true
                }
        }
        return false
}
