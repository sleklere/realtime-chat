// Package rooms provides the room list UI model for the chat client.
package rooms

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sleklere/realtime-chat/cmd/client/internal/api"
)

// RoomSelectedMsg signals that a room has been selected and joined.
type RoomSelectedMsg struct {
	Room api.RoomResponse
}

// RoomErrorMsg signals an error in room operations.
type RoomErrorMsg struct {
	Err error
}

type roomsLoadedMsg struct {
	rooms []api.RoomResponse
}

type roomCreatedMsg struct {
	room api.RoomResponse
}

type roomJoinedMsg struct {
	room api.RoomResponse
}

type roomItem struct {
	room api.RoomResponse
}

func (i roomItem) FilterValue() string { return i.room.Name }

type roomItemDelegate struct{}

func (d roomItemDelegate) Height() int                             { return 1 }
func (d roomItemDelegate) Spacing() int                            { return 0 }
func (d roomItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d roomItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(roomItem)
	if !ok {
		return
	}

	name := i.room.Name
	slug := i.room.Slug

	var str string
	if index == m.Index() {
		selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
		str = selectedStyle.Render(fmt.Sprintf("> %s ", name)) +
			lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(fmt.Sprintf("#%s", slug))
	} else {
		str = fmt.Sprintf("  %s ", name) +
			lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(fmt.Sprintf("#%s", slug))
	}

	_, _ = fmt.Fprint(w, str)
}

// Model is the Bubble Tea model for the room list screen.
type Model struct {
	apiClient   *api.Client
	list        list.Model
	creating    bool
	createInput textinput.Model
	err         string
	width       int
	height      int
}

// New creates a new rooms Model with the given API client and dimensions.
func New(apiClient *api.Client, width, height int) Model {
	l := list.New([]list.Item{}, roomItemDelegate{}, width, height-4)
	l.Title = "Rooms"
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.Styles.Title = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))

	input := textinput.New()
	input.Placeholder = "room name"
	input.CharLimit = 50

	return Model{
		apiClient:   apiClient,
		list:        l,
		createInput: input,
		width:       width,
		height:      height,
	}
}

// Init initializes the rooms model.
func (m Model) Init() tea.Cmd {
	return m.fetchRooms()
}

// Update handles messages for the rooms model.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.creating {
			return m.updateCreating(msg)
		}

		switch msg.String() {
		case "n":
			m.creating = true
			m.createInput.SetValue("")
			m.createInput.Focus()
			return m, textinput.Blink
		case "r":
			return m, m.fetchRooms()
		case "enter":
			if item, ok := m.list.SelectedItem().(roomItem); ok {
				return m, m.joinAndSelect(item.room)
			}
		}

	case roomsLoadedMsg:
		items := make([]list.Item, len(msg.rooms))
		for i, r := range msg.rooms {
			items[i] = roomItem{room: r}
		}
		m.list.SetItems(items)
		return m, nil

	case roomCreatedMsg:
		m.creating = false
		return m, m.fetchRooms()

	case roomJoinedMsg:
		return m, func() tea.Msg {
			return RoomSelectedMsg{Room: msg.room}
		}

	case RoomErrorMsg:
		m.err = msg.Err.Error()
		m.creating = false
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-4)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the rooms model.
func (m Model) View() string {
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true)

	var b strings.Builder

	b.WriteString(m.list.View())
	b.WriteString("\n")

	if m.creating {
		b.WriteString("New room: ")
		b.WriteString(m.createInput.View())
		b.WriteString("\n")
	}

	if m.err != "" {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
		b.WriteString(errorStyle.Render(m.err))
		b.WriteString("\n")
	}

	b.WriteString(helpStyle.Render("enter: join • n: new room • r: refresh • esc: quit"))

	return b.String()
}

func (m Model) updateCreating(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		name := strings.TrimSpace(m.createInput.Value())
		if name == "" {
			m.creating = false
			return m, nil
		}
		m.creating = false
		return m, m.createRoom(name)
	case "esc":
		m.creating = false
		return m, nil
	}

	var cmd tea.Cmd
	m.createInput, cmd = m.createInput.Update(msg)
	return m, cmd
}

func (m Model) fetchRooms() tea.Cmd {
	return func() tea.Msg {
		rooms, err := m.apiClient.ListRooms()
		if err != nil {
			return RoomErrorMsg{Err: err}
		}
		return roomsLoadedMsg{rooms: rooms}
	}
}

func (m Model) createRoom(name string) tea.Cmd {
	return func() tea.Msg {
		room, err := m.apiClient.CreateRoom(name)
		if err != nil {
			return RoomErrorMsg{Err: err}
		}
		return roomCreatedMsg{room: room}
	}
}

func (m Model) joinAndSelect(room api.RoomResponse) tea.Cmd {
	return func() tea.Msg {
		if err := m.apiClient.JoinRoom(room.ID); err != nil {
			return RoomErrorMsg{Err: err}
		}
		return roomJoinedMsg{room: room}
	}
}
