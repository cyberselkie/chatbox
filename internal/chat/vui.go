package chat

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/gerow/go-color"
)

const useHighPerformanceRenderer = false

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
			m.send <- ((m.username) + ": ") + txt + "\n"
			m.input.Reset()
			m.viewport.GotoBottom()
		case tea.KeyCtrlT:
			m.Choice++
			if m.Choice > 3 {
				m.Choice = 0
			}
			m.themepark(m.Choice)
		}
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		inputHeight := lipgloss.Height(m.input.View())
		verticalMarginHeight := headerHeight + footerHeight + inputHeight + 1

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.Update(msg)
			m.viewport.SetContent(m.chat)
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
		if useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}
	cmds = append(cmds, tiCmd, vpCmd)
	return m, tea.Batch(cmds...)
}

func (m Client) View() string {
	s := fmt.Sprintf("%s\n%s\n%s\n%s", m.headerView(), m.chatboxView(), m.input.View(), m.footerView())
	return s
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
		Inherit(m.viewport.Style).
		Height(m.viewport.Height)
	return chatbox.Render(m.chatMD(m.viewport.View()))
}

// Style definitions.
var (

	// General colors

	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	plaintext = lipgloss.AdaptiveColor{Light: "#301934", Dark: "#E6E6FA"}

	subtly = lipgloss.NewStyle().Foreground(subtle)

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

	rand_color = lipgloss.NewStyle().Foreground(randomColor())
)

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

func randomColor() lipgloss.Color {
	hsl := color.HSL{H: rand.Float64(), S: 0.7, L: 0.7}
	return lipgloss.Color("#" + hsl.ToHTML())
}

// custom markdown themeing

func (m Client) chatMD(input string) string {
	output, err := glamour.Render(input, m.theme)
	if err != nil {
		output = input
	}
	return output
}

// cycle through themes
func (m *Client) themepark(int) { //pick string) {
	choices := [4]string{
		"dracula",
		"dark",
		"light",
		"notty"}
	pick := choices[m.Choice]
	// randomIndex := rand.Intn(len(choices))
	// pick := choices[randomIndex]
	m.theme = pick
	println(m.theme)
}
