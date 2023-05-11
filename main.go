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
	title, desc, color string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

const (
	Colorless string = "Colorless"
	Watcher   string = "Purple"
	Ironclad  string = "Red"
	Silent    string = "Green"
	Defect    string = "Blue"
)

type model struct {
	table table.Model
	list  list.Model
	class string
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
				m.class = i.color
				cardData := cards.GetData()
				rows := []table.Row{}
				for _, card := range cardData {
					if card["Color"] == i.color {
						rows = append(rows, table.Row{card["Name"], card["Color"], card["Rarity"], card["Type"], card["Cost"], card["Text"]})
					}
				}
				m.table.SetRows(rows)
				m.view = "table"
			}
		}
	}
	m.table, cmd = m.table.Update(msg)
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m model) View() string {
	if m.view == "table" {
		return baseStyle.Render(m.table.View()) + "\n"
	}
	return "\n" + m.list.View()
}

func main() {
	columns := []table.Column{
		{Title: "Name", Width: 10},
		{Title: "Color", Width: 10},
		{Title: "Rarity", Width: 10},
		{Title: "Type", Width: 10},
		{Title: "Cost", Width: 10},
		{Title: "Text", Width: 50},
	}

	view := "list"
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
		table.WithFocused(true),
		table.WithHeight(20),
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

	m := model{table: t, list: l, view: view}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
