package conn

import (
        "context"
        "hisoka/src/handlers"
        "hisoka/src/helpers"
        "hisoka/src/settings"
        "os"
        "os/signal"
        "regexp"
        "strings"
        "syscall"
        "time"

        _ "hisoka/src/commands"

        _ "github.com/mattn/go-sqlite3"
        "github.com/mdp/qrterminal"
        "go.mau.fi/whatsmeow"
        "go.mau.fi/whatsmeow/proto/waCompanionReg"
        "go.mau.fi/whatsmeow/store"
        "go.mau.fi/whatsmeow/store/sqlstore"
        "go.mau.fi/whatsmeow/types"
        "google.golang.org/protobuf/proto"
)

type Template struct {
        Nama   string
        Status bool
}

func init() {
        store.DeviceProps.PlatformType = waCompanionReg.DeviceProps_EDGE.Enum()
        store.DeviceProps.Os = proto.String("Linux")
}

func StartClient() {
        helpers.Banner()

        ctx := context.Background()
        dbLog := helpers.PrettyWALogger("db")
        container, err := sqlstore.New(ctx, "sqlite3", "file:session.db?_foreign_keys=on", dbLog)
        if err != nil {
                panic(err)
        }
        handler := handlers.NewHandler(container)
        helpers.Step("Menyiapkan socket…")
        conn := handler.Client()
        conn.PrePairCallback = func(jid types.JID, platform, businessName string) bool {
                helpers.Connected()
                return true
        }

        if conn.Store.ID == nil {
                // No ID stored, new login
                pairingNumber := settings.PairingNumber

                if pairingNumber != "" {
                        pairingNumber = regexp.MustCompile(`\D+`).ReplaceAllString(pairingNumber, "")
                        handler.PairingNumber = pairingNumber

                        if err := conn.Connect(); err != nil {
                                panic(err)
                        }

                        var code string
                        for attempt := 1; ; attempt++ {
                                var err error
                                code, err = conn.PairPhone(context.Background(), pairingNumber, true, whatsmeow.PairClientChrome, "Edge (Linux)")
                                if err == nil {
                                        break
                                }
                                if strings.Contains(err.Error(), "rate-overlimit") || strings.Contains(err.Error(), "429") {
                                        wait := time.Duration(15*attempt) * time.Second
                                        if wait > 60*time.Second {
                                                wait = 60 * time.Second
                                        }
                                        helpers.SocketDown("WhatsApp rate-limit, tunggu " + wait.String() + " lalu coba lagi…")
                                        time.Sleep(wait)
                                        continue
                                }
                                panic(err)
                        }

                        handler.MarkPairIssued()
                        helpers.PairingPanel(pairingNumber, code)
                } else {
                        qrChan, _ := conn.GetQRChannel(context.Background())
                        if err := conn.Connect(); err != nil {
                                panic(err)
                        }

                        shown := false
                        for evt := range qrChan {
                                switch string(evt.Event) {
                                case "code":
                                        if !shown {
                                                helpers.QRPanel()
                                                shown = true
                                        }
                                        qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
                                }
                        }
                }
        } else {
                // Already logged in, just connect
                if err := conn.Connect(); err != nil {
                        panic(err)
                }
                helpers.Connected()
        }

        // Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
        c := make(chan os.Signal, 1)
        signal.Notify(c, os.Interrupt, syscall.SIGTERM)
        <-c

        conn.Disconnect()
}
