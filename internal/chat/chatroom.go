package chat

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gliderlabs/ssh"
)

func NewClient(username string, pty ssh.Pty, send chan<- string, recv chan Msg, theme string) Client {
	ti := textarea.New()
	ti.Focus()
	ti.Prompt = "// "

	ti.SetHeight(3)
	ti.SetWidth(pty.Window.Width)

	ti.KeyMap.InsertNewline.SetEnabled(true)
	ti.KeyMap.InsertNewline.Unbind()
	ti.KeyMap.InsertNewline.SetKeys("ctrl+r")

	ti.FocusedStyle.CursorLine = lipgloss.NewStyle().Background(subtle)

	ti.ShowLineNumbers = false

	return Client{
		input:    ti,
		err:      nil,
		width:    pty.Window.Width,
		height:   pty.Window.Height,
		username: username,
		send:     send,
		recv:     recv,
		theme:    theme,
		Choice:   0,
	}
}

func (chatRoom *ChatRoom) withLock(tag string, f func()) {
	log.Printf("[🔒 %s] LOCK", tag)
	chatRoom.mutex.Lock()
	f()
	log.Printf("[🔒 %s] UNLOCK", tag)
	chatRoom.mutex.Unlock()
}

func (chatRoom *ChatRoom) Subscribe(username string) chan Msg {
	ch := make(chan Msg)
	go chatRoom.withLock("SUBSCRIBE", func() {
		v, ok := chatRoom.users[username]
		log.Printf("[🔔 %s] %t %v", username, ok, v)
		if ok {
			log.Println("[Subscribe] Already subscribed")
			return
		} else {
			chatRoom.users[username] = ch
		}
		chatRoom.Inbox <- username + " has joined\n"
		ch <- MsgChat{chatRoom.history()}
		chatRoom.Blast(MsgUserList{chatRoom.GetUsers()})
	})
	return ch
}

func (chatRoom *ChatRoom) GetUsers() []string {
	users := []string{}
	for k := range chatRoom.users {
		users = append(users, k)
	}
	return users
}

func (chatRoom *ChatRoom) Unsubscribe(username string) {
	chatRoom.withLock("UNSUBSCRIBE", func() {
		delete(chatRoom.users, username)
		chatRoom.Inbox <- username + " has left\n"
		chatRoom.Blast(MsgUserList{chatRoom.GetUsers()})
	})
}
func (chatRoom *ChatRoom) history() string {
	return strings.Join(chatRoom.lines, "\n")

}

func (chatRoom *ChatRoom) Blast(m Msg) {
	log.Printf("[BLAST]📤️ %d to send\n", len(chatRoom.users))
	for _, ch := range chatRoom.users {
		ch <- m
	}
}

func (chatRoom *ChatRoom) SendAll(m Msg) {
	chatRoom.withLock("SendAll:"+m.Tag(), func() {
		log.Printf("--- 📤️ %d to send\n", len(chatRoom.users))
		for _, ch := range chatRoom.users {
			ch <- m
		}
	})
}

func logTime(tag string, f func()) {
	now := time.Now()
	f()
	after := time.Now()
	log.Printf("[⏱ %s] took %s", tag, after.Sub(now))
}

func StartChatRoom() (context.Context, context.CancelFunc, *ChatRoom) {
	chatRoom := ChatRoom{
		lines: []string{},
		users: make(map[string]chan<- Msg),
		Inbox: make(chan string),
	}
	// Entry point for new messages from subscriptions.
	go func() {
		for msg := range chatRoom.Inbox {
			logTime("SendAll", func() {
				msg = CommandsOutput(msg)
				log.Printf("RECV: %s, %d chars long", msg, len(msg))
				chatRoom.lines = append(chatRoom.lines, msg)
				chat := MsgChat{chat: (strings.Join(chatRoom.lines, "\n"))}
				chatRoom.SendAll(chat)
			})
		}
	}()

	// Function to subscribe to the chat room.

	ctx, cancel := context.WithCancel(context.Background())
	return ctx, cancel, &chatRoom
}

/**
 * Private Functions
 */
