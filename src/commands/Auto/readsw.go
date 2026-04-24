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

			// Fire MarkRead and React in parallel so they hit WhatsApp at the
			// same time (real-time, "barengan").
			var wg sync.WaitGroup

			wg.Add(1)
			go func() {
				defer wg.Done()
				conn.WA.MarkRead(
					context.Background(),
					[]types.MessageID{m.Info.ID},
					m.Info.Timestamp,
					m.Info.Chat,
					m.Info.Sender,
				)
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
				helpers.StatusViewReactLog(m.Info.PushName, m.Info.Sender.User, randomEmoji)
			} else {
				helpers.StatusViewLog(m.Info.PushName, m.Info.Sender.User)
			}
		},
	})
}
