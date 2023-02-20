package chat

import (
	"sync"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
)

type RecvMsg struct {
	Msg
}

type ChatRoom struct {
	lines []string
	mutex sync.Mutex
	users map[string]chan<- Msg
	Inbox chan string
}

type Client struct {
	input    textarea.Model
	recv     chan Msg
	send     chan<- string
	username string
	users    []string
	polls    int
	chat     string
	width    int
	height   int
	viewport viewport.Model
	err      error
	theme    string
	Choice   int
	ready    bool
	//Chosen   bool
	//Loaded   bool
}

type Msg interface {
	Tag() string
}

type MsgChat struct {
	chat string
}

func (m MsgChat) Tag() string {
	return "CHAT"
}

type MsgUserList struct {
	users []string
}

func (m MsgUserList) Tag() string {
	return "USERLIST"
}
