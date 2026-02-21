package ui

import "github.com/charmbracelet/lipgloss"

// Common style variables used across UI components.
var (
	// AccentColor is the primary accent color for the UI.
	AccentColor = lipgloss.Color("205")
	// SubtleColor is a muted color for secondary text.
	SubtleColor = lipgloss.Color("241")
	// ErrorColor is the color used for error messages.
	ErrorColor = lipgloss.Color("196")

	// TitleStyle is the style for title text.
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(AccentColor)

	// SubtitleStyle is the style for subtitle text.
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(SubtleColor)

	// ErrorStyle is the style for error text.
	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor)

	// HelpStyle is the style for help text.
	HelpStyle = lipgloss.NewStyle().
			Foreground(SubtleColor).
			Italic(true)
)
