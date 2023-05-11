package main

import (
	"fmt"
	"os"
	"stscli/cards"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	baseStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))
)

type item struct {
	title, desc string
	color       CardColors
}
type CardColors string

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

const (
	Colorless CardColors = "Colorless"
	Watcher   CardColors = "Purple"
	Ironclad  CardColors = "Red"
	Silent    CardColors = "Green"
	Defect    CardColors = "Blue"
)

type model struct {
	table table.Model
	list  list.Model
	class CardColors
	view  string
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc":
			m.view = ""
			return m, cmd
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.class = CardColors(i.color)
				m.view = "table"
			}
		}
	}
	m.table, cmd = m.table.Update(msg)
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m model) View() string {
	if m.view != "" {
		return baseStyle.Render(m.table.View()) + "\n"
	}
	return "\n" + m.list.View()
}

func main() {
	cardData := cards.GetData()
	columns := []table.Column{
		{Title: "Name", Width: 10},
		{Title: "Color", Width: 10},
		{Title: "Rarity", Width: 10},
		{Title: "Type", Width: 10},
		{Title: "Cost", Width: 10},
		{Title: "Text", Width: 50},
	}

	rows := []table.Row{}
	for _, card := range cardData {
		rows = append(rows, table.Row{card["Name"], card["Color"], card["Rarity"], card["Type"], card["Cost"], card["Text"]})
	}

	class := Watcher
	items := []list.Item{
		item{title: "Colorless", desc: "Colorless card", color: Colorless},
		item{title: "Ironclad", desc: "Ironclad", color: Ironclad},
		item{title: "Silent", desc: "Silent", color: Silent},
		item{title: "Defect", desc: "Defect", color: Defect},
		item{title: "Watcher", desc: "Watcher", color: Watcher},
	}

	l := list.New(items, list.NewDefaultDelegate(), 20, 20)
	l.Title = "Which class cards do you want to search?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := model{table: t, list: l, class: class}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
