package helpers

import (
        "fmt"
        "strings"
        "time"
)

const (
        cReset = "\x1b[0m"
        cBold  = "\x1b[1m"
        cDim   = "\x1b[2m"

        cCyan    = "\x1b[36m"
        cGreen   = "\x1b[32m"
        cYellow  = "\x1b[33m"
        cMagenta = "\x1b[35m"
        cRed     = "\x1b[31m"
        cBlue    = "\x1b[94m"
        cGray    = "\x1b[90m"
)

const uiWidth = 44

func Banner() {
        bar := strings.Repeat("─", uiWidth-2)
        fmt.Println()
        fmt.Println(cCyan + "╭" + bar + "╮" + cReset)
        fmt.Println(cCyan + "│" + cReset + center(cBold+"✦ HISOKA · WhatsApp Bot ✦"+cReset, uiWidth-2) + cCyan + "│" + cReset)
        fmt.Println(cCyan + "│" + cReset + center(cDim+"Go + whatsmeow"+cReset, uiWidth-2) + cCyan + "│" + cReset)
        fmt.Println(cCyan + "╰" + bar + "╯" + cReset)
        fmt.Println()
}

func PairingPanel(number, code string) {
        row("⚙", "Mode", cMagenta+"Pairing Code"+cReset)
        row("📱", "Nomor", number)
        row("🔑", "Kode", cBold+cYellow+code+cReset)
        row("⏱", "Berlaku", "± 60 detik")
        fmt.Println()
        tutorial("CARA PAKAI", []string{
                "Buka WhatsApp di HP kamu",
                "Tap " + cBold + "⋮" + cReset + " → " + cBold + "Perangkat tertaut" + cReset,
                "Pilih " + cBold + "Tautkan dgn nomor telepon" + cReset,
                "Ketik kode " + cBold + cYellow + code + cReset + " di HP",
        })
}

func QRPanel() {
        row("⚙", "Mode", cMagenta+"QR Code"+cReset)
        fmt.Println()
        tutorial("CARA SCAN", []string{
                "Buka WhatsApp di HP kamu",
                "Tap " + cBold + "⋮" + cReset + " → " + cBold + "Perangkat tertaut" + cReset,
                "Pilih " + cBold + "Tautkan perangkat" + cReset,
                "Arahkan kamera ke QR di bawah",
        })
}

func Connected() {
        fmt.Println()
        fmt.Println("  " + cGreen + "✓ " + cBold + "Connected" + cReset + cDim + " — siap menerima pesan" + cReset)
        fmt.Println()
}

func SocketDown(msg string) {
        fmt.Println("  " + cYellow + "⚠ " + cReset + msg)
}

func PairingRefresh(code string) {
        fmt.Println()
        fmt.Println("  " + cMagenta + "↻ " + cBold + "Kode pairing baru" + cReset + cDim + " (yang lama kadaluarsa)" + cReset)
        fmt.Println("  🔑  " + cBold + "Kode    " + cReset + cDim + ": " + cReset + cBold + cYellow + code + cReset)
        fmt.Println("  ⏱  " + cBold + "Berlaku " + cReset + cDim + ": " + cReset + "± 60 detik")
        fmt.Println("  " + cDim + "Ketik kode di HP: ⋮ → Perangkat tertaut → Tautkan dgn nomor telepon" + cReset)
        fmt.Println()
}

func Step(msg string) {
        fmt.Println("  " + cCyan + "›" + cReset + " " + msg)
}

func MessageLog(name, sender, command, body, msgType string) {
        fmt.Println("  " + cGray + strings.Repeat("─", uiWidth-2) + cReset)
        field("From", cBlue+name+cReset+cDim+" ("+sender+")"+cReset)
        if command != "" {
                field("Cmd", cYellow+command+cReset)
        }
        if len(body) >= 350 {
                body = cDim + "<" + msgType + ">" + cReset
        }
        field("Msg", body)
}

// StatusViewLog prints a compact, mobile-friendly notification when the
// bot auto-views (auto-reads) someone's WhatsApp status / story.
func StatusViewLog(name, sender string) {
        statusBox("👁", "LIHAT STORY", cMagenta, name, sender, "")
}

// StatusReactLog prints a compact notification when the bot auto-reacts
// to someone's WhatsApp status with an emoji.
func StatusReactLog(name, sender, emoji string) {
        statusBox("💖", "REACT STORY", cYellow, name, sender, emoji)
}

// statusBox renders a 2-line rounded box that fits 44-col mobile consoles.
// Layout (uiWidth=44):
//
//      ╭─ 👁 LIHAT STORY ──────── 14:23 ─╮
//      │ Wily Wonka · 6285198168303      │
//      ╰─────────────────────────────────╯
func statusBox(icon, title, color, name, sender, extra string) {
        if name == "" {
                name = "(tanpa nama)"
        }
        stamp := time.Now().Format("15:04")
        inner := uiWidth - 2 // chars between the corner glyphs

        // --- header ---
        // "─ <icon> <TITLE> " + dashes + " <stamp> ─"
        left := " " + icon + " " + title + " "
        right := " " + stamp + " "
        leftLen := visibleLen(left)
        rightLen := visibleLen(right)
        dashes := inner - leftLen - rightLen - 2 // 2 = leading & trailing dash
        if dashes < 1 {
                dashes = 1
        }
        header := cGray + "╭─" + cReset +
                color + left + cReset +
                cGray + strings.Repeat("─", dashes) + cReset +
                cDim + right + cReset +
                cGray + "─╮" + cReset
        fmt.Println("  " + header)

        // --- body ---
        body := cBlue + name + cReset + cDim + " · " + cReset + sender
        if extra != "" {
                body = body + cDim + "  " + cReset + extra
        }
        body = padVisible(" "+body, inner)
        fmt.Println("  " + cGray + "│" + cReset + body + cGray + "│" + cReset)

        // --- footer ---
        footer := cGray + "╰" + strings.Repeat("─", inner) + "╯" + cReset
        fmt.Println("  " + footer)
}

// padVisible pads s with spaces so its visible (printable) length equals width.
func padVisible(s string, width int) string {
        v := visibleLen(s)
        if v >= width {
                // truncate gracefully on rune boundaries
                runes := []rune(stripANSI(s))
                if len(runes) > width {
                        runes = runes[:width-1]
                        return string(runes) + "…"
                }
                return s
        }
        return s + strings.Repeat(" ", width-v)
}

// stripANSI removes ANSI color escapes so we can measure / truncate text.
func stripANSI(s string) string {
        var b strings.Builder
        in := false
        for _, r := range s {
                if r == '\x1b' {
                        in = true
                        continue
                }
                if in {
                        if r == 'm' {
                                in = false
                        }
                        continue
                }
                b.WriteRune(r)
        }
        return b.String()
}

func field(label, value string) {
        fmt.Printf("  "+cCyan+"%-4s"+cReset+cDim+" : "+cReset+"%s\n", label, value)
}

func row(icon, label, value string) {
        fmt.Printf("  %s  "+cBold+"%-8s"+cReset+cDim+": "+cReset+"%s\n", icon, label, value)
}

func tutorial(title string, steps []string) {
        inner := uiWidth - len(title) - 4
        if inner < 2 {
                inner = 2
        }
        left := inner / 2
        right := inner - left
        fmt.Println("  " + cDim + strings.Repeat("─", left) + cReset + " " + cBold + title + cReset + " " + cDim + strings.Repeat("─", right) + cReset)
        for i, s := range steps {
                fmt.Printf("  "+cYellow+"%d."+cReset+" %s\n", i+1, s)
        }
        fmt.Println()
}

func center(s string, width int) string {
        v := visibleLen(s)
        if v >= width {
                return s
        }
        left := (width - v) / 2
        right := width - v - left
        return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

func visibleLen(s string) int {
        n, in := 0, false
        for _, r := range s {
                if r == 0x1b {
                        in = true
                        continue
                }
                if in {
                        if r == 'm' {
                                in = false
                        }
                        continue
                }
                n += runeCellWidth(r)
        }
        return n
}

// runeCellWidth approximates how many terminal cells a rune occupies.
// Emoji and CJK characters are treated as width 2; everything else as 1.
// This keeps box-drawing alignment correct on most modern terminals.
func runeCellWidth(r rune) int {
        switch {
        case r == 0:
                return 0
        case r < 0x20: // control
                return 0
        case r >= 0x1F000: // most emoji blocks (Misc Symbols & Pictographs and beyond)
                return 2
        case r >= 0x2600 && r <= 0x27BF: // Misc Symbols + Dingbats (☀️ ✨ ✅ etc.)
                return 2
        case r >= 0x3000 && r <= 0x9FFF: // CJK + punctuation
                return 2
        case r >= 0xFE0F && r <= 0xFE0F: // variation selector — invisible
                return 0
        case r >= 0x200D && r <= 0x200D: // ZWJ — invisible
                return 0
        }
        return 1
}
