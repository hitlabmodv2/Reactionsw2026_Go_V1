package handlers

import (
        "context"
        "hisoka/src/helpers"
        "hisoka/src/libs"
        "os"
        "regexp"
        "strings"
        "sync/atomic"
        "time"

        "go.mau.fi/whatsmeow"
        "go.mau.fi/whatsmeow/store"
        "go.mau.fi/whatsmeow/store/sqlstore"
        "go.mau.fi/whatsmeow/types"
        "go.mau.fi/whatsmeow/types/events"
)

type IHandler struct {
        Container     *store.Device
        PairingNumber string
        pairingActive int32
        lastPairAt    time.Time
        pairBackoff   time.Duration
}

func NewHandler(container *sqlstore.Container) *IHandler {
        ctx := context.Background()
        deviceStore, err := container.GetFirstDevice(ctx)
        if err != nil {
                panic(err)
        }
        return &IHandler{
                Container: deviceStore,
        }
}

func (h *IHandler) Client() *whatsmeow.Client {
        clientLog := helpers.PrettyWALogger("client")
        conn := whatsmeow.NewClient(h.Container, clientLog)
        conn.AddEventHandler(h.RegisterHandler(conn))
        return conn
}

func (h *IHandler) RegisterHandler(conn *whatsmeow.Client) func(evt interface{}) {
        return func(evt interface{}) {
                sock := libs.SerializeClient(conn)
                switch v := evt.(type) {
                case *events.Message:
                        m := libs.SerializeMessage(v, sock)

                        // skip deleted message
                        if m.Message.GetProtocolMessage() != nil && m.Message.GetProtocolMessage().GetType() == 0 {
                                return
                        }

                        // log only when the message triggers a bot command, not every
                        // random message from groups/channels.
                        if m.Body != "" && libs.HasCommand(m.Command) {
                                helpers.MessageLog(v.Info.PushName, m.Info.Sender.User, m.Command, m.Body, string(m.Info.Type))
                        }

                        // Get command
                        go ExecuteCommand(sock, m)
                        return
                case *events.Connected, *events.PushNameSetting:
                        if len(conn.Store.PushName) == 0 {
                                return
                        }
                        conn.SendPresence(context.Background(), types.PresenceAvailable)
                case *events.Disconnected:
                        if conn.Store.ID == nil && h.PairingNumber != "" {
                                h.refreshPairing(conn)
                        } else if conn.Store.ID != nil {
                                helpers.SocketDown("Terputus dari WhatsApp — menyambung ulang…")
                        }
                case *events.StreamReplaced:
                        helpers.SocketDown("Sesi diambil alih oleh device lain.")
                case *events.LoggedOut:
                        helpers.SocketDown("Sesi logout — hapus session.db lalu jalankan ulang untuk pairing baru.")
                case *events.ClientOutdated:
                        helpers.SocketDown("Versi client outdated — update whatsmeow lalu rebuild.")
                }
        }
}

func (h *IHandler) MarkPairIssued() {
        h.lastPairAt = time.Now()
        h.pairBackoff = 0
}

func (h *IHandler) refreshPairing(conn *whatsmeow.Client) {
        if !atomic.CompareAndSwapInt32(&h.pairingActive, 0, 1) {
                return
        }
        go func() {
                defer atomic.StoreInt32(&h.pairingActive, 0)

                // Honor minimum 65s gap between pair requests + any active backoff.
                minGap := 65 * time.Second
                if h.pairBackoff > minGap {
                        minGap = h.pairBackoff
                }
                if !h.lastPairAt.IsZero() {
                        elapsed := time.Since(h.lastPairAt)
                        if elapsed < minGap {
                                time.Sleep(minGap - elapsed)
                        }
                }
                if conn.Store.ID != nil {
                        return
                }

                // Wait for whatsmeow auto-reconnect to come back up.
                deadline := time.Now().Add(90 * time.Second)
                for !conn.IsConnected() {
                        if conn.Store.ID != nil {
                                return
                        }
                        if time.Now().After(deadline) {
                                helpers.SocketDown("Tidak bisa nyambung ulang ke WhatsApp, coba restart workflow.")
                                return
                        }
                        time.Sleep(2 * time.Second)
                }

                newCode, err := conn.PairPhone(context.Background(), h.PairingNumber, true, whatsmeow.PairClientChrome, "Edge (Linux)")
                if err != nil {
                        if strings.Contains(err.Error(), "rate-overlimit") || strings.Contains(err.Error(), "429") {
                                if h.pairBackoff == 0 {
                                        h.pairBackoff = 90 * time.Second
                                } else {
                                        h.pairBackoff *= 2
                                }
                                if h.pairBackoff > 5*time.Minute {
                                        h.pairBackoff = 5 * time.Minute
                                }
                                helpers.SocketDown("Rate-limit WhatsApp, jeda " + h.pairBackoff.String() + " sebelum kode baru dikeluarkan.")
                                return
                        }
                        helpers.SocketDown("Gagal minta pairing code baru: " + err.Error())
                        return
                }
                h.lastPairAt = time.Now()
                h.pairBackoff = 0
                helpers.PairingRefresh(newCode)
        }()
}

func ExecuteCommand(c *libs.IClient, m *libs.IMessage) {
        var prefix string
        pattern := regexp.MustCompile(os.Getenv("PREFIX"))
        for _, f := range pattern.FindAllString(m.Command, -1) {
                prefix = f
        }
        lists := libs.GetList()
        for _, cmd := range lists {
                if cmd.Before != nil {
                        cmd.Before(c, m)
                }
                re := regexp.MustCompile(`^` + cmd.Name + `$`)
                if valid := len(re.FindAllString(strings.ReplaceAll(m.Command, prefix, ""), -1)) > 0; valid {
                        if cmd.Execute != nil {
                                if os.Getenv("PUBLIC") == "false" && !m.IsOwner {
                                        return
                                }

                                var cmdWithPref bool
                                var cmdWithoutPref bool
                                if cmd.IsPrefix && (prefix != "" && strings.HasPrefix(m.Command, prefix)) {
                                        cmdWithPref = true
                                } else {
                                        cmdWithPref = false
                                }

                                if !cmd.IsPrefix {
                                        cmdWithoutPref = true
                                } else {
                                        cmdWithoutPref = false
                                }

                                if !cmdWithPref && !cmdWithoutPref {
                                        continue
                                }

                                if cmd.IsOwner && !m.IsOwner {
                                        continue
                                }

                                if cmd.IsQuery && m.Text == "" {
                                        m.Reply("Query Required")
                                        continue
                                }

                                if cmd.IsGroup && !m.Info.IsGroup {
                                        m.Reply("Commands only work in Group Chat")
                                        continue
                                }

                                if cmd.IsPrivate && m.Info.IsGroup {
                                        m.Reply("Commands only work in Private Chat")
                                        continue
                                }

                                if cmd.IsMedia && m.IsMedia == "" {
                                        m.Reply("Reply to Media Message, or send Media with Command")
                                        continue
                                }

                                if cmd.IsWait {
                                        m.React("⏳")
                                }

                                ok := cmd.Execute(c, m)

                                if cmd.IsWait && !ok {
                                        m.React("❌")
                                }

                                if cmd.IsWait && ok {
                                        c.WA.MarkRead(context.Background(), []types.MessageID{m.Info.ID}, time.Now(), m.Info.Chat, m.Info.Sender)
                                        m.React("")
                                        continue
                                }
                        }
                }
        }
}
