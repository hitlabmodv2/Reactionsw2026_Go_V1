package commands

import (
	"context"
	"hisoka/src/helpers"
	"hisoka/src/libs"
	"math/rand"
	"os"
	"sync"
	"time"

	"go.mau.fi/whatsmeow/types"
)

func init() {
	libs.NewCommands(&libs.ICommand{
		Before: func(conn *libs.IClient, m *libs.IMessage) {
			if m.Info.Chat.String() != "status@broadcast" {
				return
			}

			reactStatus := os.Getenv("REACT_STATUS") == "true"

			var randomEmoji string
			if reactStatus {
				emojis := []string{"😀", "😃", "😄", "😁", "😆", "🥹", "😅", "😂", "🤣", "🥲", "☺️", "😊", "😇", "🙂", "🙃", "😉", "😌", "😍", "🥰", "😘", "😗", "😙", "😚", "😋", "😛", "😝", "🤪", "🤨", "🧐", "🤓", "😎", "🥸", "🤩", "🥳", "😏", "😒", "😞", "😔", "😟", "😕", "🙁", "☹️", "😣", "😖", "😫", "😩", "🥺", "😢", "😭", "😤", "😠", "😡", "🤬", "🤯", "😳", "🥵", "🥶", "😶‍🌫️", "😱", "😨", "😰", "😥", "😓", "🤗", "🤔", "🫣", "🤭", "🫢", "🫡", "🤫", "🫠", "🤥", "😶", "🫥", "😐", "🫤", "😑", "😬", "🙄", "😯", "😦", "😧", "😮", "😲", "🥱", "😴", "🤤", "😪", "😮‍💨", "😵", "😵‍💫", "🤐", "🥴", "🤢", "🤮", "🤧", "😷", "🤒", "🤕", "🤑", "🤡", "💩", "👻", "💀", "☠️", "🙌", "👏", "👍", "👎", "👊", "✊", "🤛", "🤞", "✌️", "🫰", "🤟", "🤘", "👌", "🤏", "☝️", "✋", "🤚", "🖖", "👋", "🤙", "🫲", "🫱", "💪", "🖕", "✍️", "🙏", "🫵", "🦶", "👣", "👀", "🧠"}
				rand.Seed(time.Now().UnixNano())
				randomEmoji = emojis[rand.Intn(len(emojis))]
			}

			// Pick the right sender JID. Newer WhatsApp uses "lid" addressing
			// for some contacts; in that case m.Info.Sender is "<id>@lid",
			// which the read-receipt path on the server silently ignores.
			// SerializeMessage already resolves the proper PN-format JID into
			// m.Sender, so prefer that and fall back if it's empty.
			senderJID := m.Sender
			if senderJID.IsEmpty() {
				senderJID = m.Info.Sender
			}

			// Fire MarkRead and React in parallel so they hit WhatsApp at the
			// same time (real-time, "barengan").
			var wg sync.WaitGroup

			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := conn.WA.MarkRead(
					context.Background(),
					[]types.MessageID{m.Info.ID},
					m.Info.Timestamp,
					m.Info.Chat,
					senderJID,
					types.ReceiptTypeRead,
				); err != nil {
					helpers.ErrorLogger.Println("MarkRead status:", err)
				}
			}()

			if reactStatus {
				wg.Add(1)
				go func() {
					defer wg.Done()
					m.React(randomEmoji)
				}()
			}

			wg.Wait()

			// One single, combined log card so the console doesn't show
			// two separate boxes for what is really one synchronized event.
			if reactStatus {
				helpers.StatusViewReactLog(m.Info.PushName, senderJID.User, randomEmoji)
			} else {
				helpers.StatusViewLog(m.Info.PushName, senderJID.User)
			}
		},
	})
}
