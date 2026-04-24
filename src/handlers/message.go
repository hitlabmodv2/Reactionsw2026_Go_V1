package handlers

import (
        "context"
        "hisoka/src/helpers"
        "hisoka/src/libs"
        "os"
        "regexp"
        "strings"
        "time"

        "go.mau.fi/whatsmeow"
        "go.mau.fi/whatsmeow/store"
        "go.mau.fi/whatsmeow/store/sqlstore"
        "go.mau.fi/whatsmeow/types"
        "go.mau.fi/whatsmeow/types/events"
)

type IHandler struct {
        Container *store.Device
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

                        // log
                        if m.Body != "" {
                                cmd := ""
                                if libs.HasCommand(m.Command) {
                                        cmd = m.Command
                                }
                                helpers.MessageLog(v.Info.PushName, m.Info.Sender.User, cmd, m.Body, string(m.Info.Type))
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
                        helpers.SocketDown("Terputus dari WhatsApp — menyambung ulang…")
                case *events.StreamReplaced:
                        helpers.SocketDown("Sesi diambil alih oleh device lain.")
                case *events.LoggedOut:
                        helpers.SocketDown("Sesi logout — hapus session.db lalu jalankan ulang untuk pairing baru.")
                case *events.ClientOutdated:
                        helpers.SocketDown("Versi client outdated — update whatsmeow lalu rebuild.")
                }
        }
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
