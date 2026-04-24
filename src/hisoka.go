package conn

import (
        "context"
        "hisoka/src/handlers"
        "hisoka/src/helpers"
        "os"
        "os/signal"
        "regexp"
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
                pairingNumber := os.Getenv("PAIRING_NUMBER")

                if pairingNumber != "" {
                        pairingNumber = regexp.MustCompile(`\D+`).ReplaceAllString(pairingNumber, "")

                        if err := conn.Connect(); err != nil {
                                panic(err)
                        }

                        pairCtx := context.Background()
                        code, err := conn.PairPhone(pairCtx, pairingNumber, true, whatsmeow.PairClientChrome, "Edge (Linux)")
                        if err != nil {
                                panic(err)
                        }

                        helpers.PairingPanel(pairingNumber, code)

                        // Auto-refresh pairing code while user has not entered it yet.
                        go func() {
                                for {
                                        time.Sleep(65 * time.Second)
                                        if conn.Store.ID != nil {
                                                return
                                        }
                                        for !conn.IsConnected() && conn.Store.ID == nil {
                                                time.Sleep(2 * time.Second)
                                        }
                                        if conn.Store.ID != nil {
                                                return
                                        }
                                        newCode, err := conn.PairPhone(pairCtx, pairingNumber, true, whatsmeow.PairClientChrome, "Edge (Linux)")
                                        if err != nil {
                                                helpers.SocketDown("Gagal minta pairing code baru: " + err.Error())
                                                continue
                                        }
                                        helpers.PairingRefresh(newCode)
                                }
                        }()
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
