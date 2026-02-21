// Package chat provides the chat room UI model for the chat client.
package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sleklere/realtime-chat/cmd/client/internal/api"
	"github.com/sleklere/realtime-chat/cmd/client/internal/ws"
)

// LeaveRoomMsg signals that the user wants to leave the current room.
type LeaveRoomMsg struct{}

type historyLoadedMsg struct {
	messages []api.MessageResponse
}

type wsConnectedMsg struct {
	client *ws.Client
}

// Model is the Bubble Tea model for the chat room screen.
type Model struct {
	apiClient *api.Client
	wsClient  *ws.Client
	program   *tea.Program
	logger    *slog.Logger

	room     api.RoomResponse
	userID   int64
	username string
	wsURL    string
	token    string

	viewport viewport.Model
	input    textinput.Model
	messages []chatMessage
	err      string
	width    int
	height   int
}

type chatMessage struct {
	senderID       int64
	senderUsername string
	content        string
	timestamp      string
}

// New creates a new chat Model for the given room.
func New(
	apiClient *api.Client,
	program *tea.Program,
	logger *slog.Logger,
	room api.RoomResponse,
	userID int64,
	username, wsURL, token string,
	width, height int,
) Model {
	vp := viewport.New(width, height-5)

	input := textinput.New()
	input.Placeholder = "type a message..."
	input.Focus()
	input.CharLimit = 500
	input.Width = width - 4

	return Model{
		apiClient: apiClient,
		program:   program,
		logger:    logger,
		room:      room,
		userID:    userID,
		username:  username,
		wsURL:     wsURL,
		token:     token,
		viewport:  vp,
		input:     input,
		width:     width,
		height:    height,
	}
}

// Init initializes the chat model.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.loadHistory(),
		m.connectWS(),
		textinput.Blink,
	)
}

// Update handles messages for the chat model.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.cleanup()
			return m, func() tea.Msg { return LeaveRoomMsg{} }
		case "enter":
			return m.sendMessage()
		}

	case historyLoadedMsg:
		for _, m2 := range msg.messages {
			m.messages = append(m.messages, chatMessage{
				senderID:       m2.SenderID,
				senderUsername: fmt.Sprintf("user_%d", m2.SenderID),
				content:        m2.Body,
				timestamp:      m2.CreatedAt.Format("15:04"),
			})
		}
		m.updateViewport()
		return m, nil

	case wsConnectedMsg:
		m.wsClient = msg.client
		m.logger.Info("ws connected for chat", "room_id", m.room.ID)
		return m, nil

	case ws.IncomingMsg:
		return m.handleWSMessage(msg)

	case ws.ErrorMsg:
		m.err = msg.Err.Error()
		m.logger.Error("ws error in chat", "error", msg.Err)
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 5
		m.input.Width = msg.Width - 4
		m.updateViewport()
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the chat model.
func (m Model) View() string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color("241")).
		Width(m.width)

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(lipgloss.Color("241")).
		Width(m.width)

	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true)

	var b strings.Builder

	header := fmt.Sprintf(" #%s — %s", m.room.Slug, m.room.Name)
	b.WriteString(headerStyle.Render(header))
	b.WriteString("\n")
	b.WriteString(m.viewport.View())
	b.WriteString("\n")
	b.WriteString(inputStyle.Render(" " + m.input.View()))
	b.WriteString("\n")

	if m.err != "" {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
		b.WriteString(errorStyle.Render(m.err))
		b.WriteString(" ")
	}

	b.WriteString(helpStyle.Render("esc: leave room • enter: send"))

	return b.String()
}

func (m Model) sendMessage() (Model, tea.Cmd) {
	content := strings.TrimSpace(m.input.Value())
	if content == "" || m.wsClient == nil {
		return m, nil
	}

	m.input.SetValue("")

	if err := m.wsClient.SendRoomMessage(m.room.ID, content); err != nil {
		m.err = err.Error()
	}

	return m, nil
}

func (m Model) handleWSMessage(msg ws.IncomingMsg) (Model, tea.Cmd) {
	switch msg.Message.Type {
	case ws.TypeRoomMessage:
		var payload ws.RoomMessagePayload
		if err := json.Unmarshal(msg.Message.Payload, &payload); err != nil {
			m.logger.Error("failed to unmarshal room message", "error", err)
			return m, nil
		}

		m.messages = append(m.messages, chatMessage{
			senderID:       payload.SenderID,
			senderUsername: payload.SenderUsername,
			content:        payload.Content,
			timestamp:      msg.Message.Timestamp.Format("15:04"),
		})
		m.updateViewport()

	case ws.TypeError:
		var payload ws.ErrorPayload
		if err := json.Unmarshal(msg.Message.Payload, &payload); err == nil {
			m.err = payload.Message
		}
	}

	return m, nil
}

func (m *Model) updateViewport() {
	var lines []string
	ownStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	timeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	nameStyle := lipgloss.NewStyle().Bold(true)

	for _, msg := range m.messages {
		ts := timeStyle.Render(fmt.Sprintf("[%s]", msg.timestamp))
		var name string
		if msg.senderID == m.userID {
			name = ownStyle.Render("you")
		} else {
			name = nameStyle.Render(msg.senderUsername)
		}
		lines = append(lines, fmt.Sprintf("%s %s: %s", ts, name, msg.content))
	}

	m.viewport.SetContent(strings.Join(lines, "\n"))
	m.viewport.GotoBottom()
}

func (m *Model) cleanup() {
	if m.wsClient != nil {
		m.wsClient.Close()
		m.wsClient = nil
	}
}

func (m Model) loadHistory() tea.Cmd {
	return func() tea.Msg {
		messages, err := m.apiClient.GetMessages(m.room.ID, 50)
		if err != nil {
			return ws.ErrorMsg{Err: err}
		}
		return historyLoadedMsg{messages: messages}
	}
}

func (m Model) connectWS() tea.Cmd {
	return func() tea.Msg {
		client, err := ws.Connect(context.Background(), m.wsURL, m.token, m.program, m.logger)
		if err != nil {
			return ws.ErrorMsg{Err: err}
		}
		return wsConnectedMsg{client: client}
	}
}
