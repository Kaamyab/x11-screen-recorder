package util

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choice int
}

var options = []string{"Capture Whole Screen", "Capture Specific Window"}

func RunTerminalUI() error {
	p := tea.NewProgram(initialModel())
	return p.Start()
}

func initialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1":
			m.choice = 1
			return m, captureWholeScreen
		case "2":
			m.choice = 2
			return m, captureWindow
		}
	}
	return m, nil
}

func (m model) View() string {
	return "Choose an option:\n" +
		"1. Capture Whole Screen\n" +
		"2. Capture Specific Window\n" +
		"\nPress q to quit."
}

func captureWholeScreen() tea.Msg {
	fmt.Println("Capturing the whole screen...")
	// Implement screen capturing logic here
	return nil
}

func captureWindow() tea.Msg {
	fmt.Println("Please click on a window to capture...")
	//windowID, err := capture.SelectWindow()
	//if err != nil {
	//	fmt.Printf("Error selecting window: %v\n", err)
	//	return nil
	//}
	//fmt.Printf("Selected window ID: %d\n", windowID)
	//// Implement window capturing logic here using the windowID
	return nil
}
