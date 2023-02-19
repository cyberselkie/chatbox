package chat

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	errMsg error
)

func (m Client) Init() tea.Cmd {
	return tea.Batch(m.pollChat, textinput.Blink)
}

func (m Client) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
		cmds  []tea.Cmd
	)

	m.input, tiCmd = m.input.Update(msg)
	//m.viewport, vpCmd = m.viewport.Update(msg)
	switch msg := msg.(type) {
	case RecvMsg:
		return m.handleRecvMsg(msg)
	case tea.KeyMsg:
		switch msg.Type {
		// quit on ctrl+c
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyUp:
			m.viewport, vpCmd = m.viewport.Update(msg)
		case tea.KeyDown:
			m.viewport, vpCmd = m.viewport.Update(msg)
		case tea.KeyEnter:
			i := m.input
			txt := i.Value()
			m.send <- i.PromptStyle.Render(i.Prompt) + i.TextStyle.Render(txt)
			m.input.Reset()
			m.viewport.GotoBottom()
		}
	case errMsg:
		m.err = msg
		return m, nil
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		inputHeight := lipgloss.Height(m.input.View())
		verticalMarginHeight := headerHeight + footerHeight + inputHeight + 4
		m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - verticalMarginHeight
		m.viewport.Update(msg)
		m.viewport.SetContent(m.chat)
	}
	cmds = append(cmds, tiCmd, vpCmd)
	return m, tea.Batch(cmds...)
}
func (m Client) headerView() string {
	title := header.Render("epic_chat")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Client) footerView() string {
	info := header.Render("online_users")
	userlist := " // "
	userlist += roll_style.Render(strings.Join(m.users, " / "))
	userlist += " // "
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)-lipgloss.Width(userlist)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, userlist, info)
}

func (m Client) chatboxView() string {
	//style the chatbox
	var chatbox = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(pink).
		Padding(0, 1).
		Width(m.viewport.Width - 2).
		Inherit(m.viewport.Style)
	return chatbox.Render(m.viewport.View())
}

// Style definitions.
var (

	// General colors

	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	plaintext = lipgloss.AdaptiveColor{Light: "#301934", Dark: "#E6E6FA"}

	// pastel palette
	pink = lipgloss.Color("#ffafcc")
	blue = lipgloss.Color("#a2d2ff")

	// chat header
	header = lipgloss.NewStyle().
		Bold(true).
		Foreground(pink).
		Padding(1)

	// standout colors in commands
	roll_style = lipgloss.NewStyle().
			Foreground(blue)
)

func (m Client) View() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s", m.headerView(), m.chatboxView(), m.input.View(), m.footerView())
}

/**
 * Private Functions
 */

func (m Client) pollChat() tea.Msg {
	chat := <-m.recv
	m.viewport.SetContent(m.chat)
	return RecvMsg{Msg: chat}
}

func (m *Client) handleRecvMsg(msg RecvMsg) (tea.Model, tea.Cmd) {
	switch msg := msg.Msg.(type) {
	case MsgChat:
		m.chat = msg.chat
	case MsgUserList:
		m.users = msg.users
	}
	m.viewport.SetContent(m.chat)
	m.viewport.GotoBottom()
	m.polls = m.polls + 1
	return m, m.pollChat
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
