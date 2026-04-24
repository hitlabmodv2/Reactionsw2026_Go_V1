package commands

import (
        "context"
        "hisoka/src/helpers"
        "hisoka/src/libs"
        "math/rand"
        "os"
        "time"

        "go.mau.fi/whatsmeow/types"
)

func init() {
        libs.NewCommands(&libs.ICommand{
                Before: func(conn *libs.IClient, m *libs.IMessage) {
                        if m.Info.Chat.String() == "status@broadcast" {
                                conn.WA.MarkRead(context.Background(), []types.MessageID{m.Info.ID}, m.Info.Timestamp, m.Info.Chat, m.Info.Sender)
                                helpers.StatusViewLog(m.Info.PushName, m.Info.Sender.User)

                                if os.Getenv("REACT_STATUS") == "true" {
                                        emojis := []string{"😀", "😃", "😄", "😁", "😆", "🥹", "😅", "😂", "🤣", "🥲", "☺️", "😊", "😇", "🙂", "🙃", "😉", "😌", "😍", "🥰", "😘", "😗", "😙", "😚", "😋", "😛", "😝", "🤪", "🤨", "🧐", "🤓", "😎", "🥸", "🤩", "🥳", "😏", "😒", "😞", "😔", "😟", "😕", "🙁", "☹️", "😣", "😖", "😫", "😩", "🥺", "😢", "😭", "😤", "😠", "😡", "🤬", "🤯", "😳", "🥵", "🥶", "😶‍🌫️", "😱", "😨", "😰", "😥", "😓", "🤗", "🤔", "🫣", "🤭", "🫢", "🫡", "🤫", "🫠", "🤥", "😶", "🫥", "😐", "🫤", "😑", "😬", "🙄", "😯", "😦", "😧", "😮", "😲", "🥱", "😴", "🤤", "😪", "😮‍💨", "😵", "😵‍💫", "🤐", "🥴", "🤢", "🤮", "🤧", "😷", "🤒", "🤕", "🤑", "🤡", "💩", "👻", "💀", "☠️", "🙌", "👏", "👍", "👎", "👊", "✊", "🤛", "🤞", "✌️", "🫰", "🤟", "🤘", "👌", "🤏", "☝️", "✋", "🤚", "🖖", "👋", "🤙", "🫲", "🫱", "💪", "🖕", "✍️", "🙏", "🫵", "🦶", "👣", "👀", "🧠"}
                                        rand.Seed(time.Now().UnixNano())
                                        randomEmoji := emojis[rand.Intn(len(emojis))]

                                        m.React(randomEmoji)
                                        helpers.StatusReactLog(m.Info.PushName, m.Info.Sender.User, randomEmoji)
                                }
                        }
                },
        })
}
