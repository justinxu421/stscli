package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	//	"github.com/charmbracelet/bubbles/textarea"
	//	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type status int

const divisor = 4

const (
	todo status = iota
	inProgress
	done
)

type Task struct {
	status      status
	title       string
	description string
}

func (t *Task) Next() {
	if t.status == done {
		t.status = todo
	} else {
		t.status++
	}
}

// implement the list.Item interface
func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

/* MAIN MODEL */

type Model struct {
	lists   []list.Model
	focused status
	loaded  bool
	err     error
}

func New() *Model {
	return &Model{}
}

// TODO: call this on tea.WindowSizeMsg
func (m *Model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height/divisor)
	m.lists = []list.Model{defaultList, defaultList, defaultList}

	// Init todos
	m.lists[todo].Title = "To Do"
	m.lists[todo].SetItems([]list.Item{
		Task{status: todo, title: "buy milk", description: "strawberry milk"},
		Task{status: todo, title: "eat sushi", description: "negitoro roll, miso soup, rice"},
		Task{status: todo, title: "fold laundry", description: "or wear wrinkly t-shirts"},
	})

	// Init in progress
	m.lists[inProgress].Title = "In Progress"
	m.lists[inProgress].SetItems([]list.Item{
		Task{status: todo, title: "write code", description: "don't worry , it's Go"},
	})

	// Init done
	m.lists[done].Title = "Done"
	m.lists[done].SetItems([]list.Item{
		Task{status: todo, title: "stay cool", description: "as a cucumber"},
	})
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
            m.initLists(msg.Width, msg.Height)
			m.loaded = true
		}
		m.initLists(msg.Width, msg.Height)
	}
	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.loaded {
		return lipgloss.JoinHorizontal(lipgloss.Left, m.lists[todo].View(), m.lists[inProgress].View(), m.lists[done].View())
	} else {
		return "loading..."
	}
}

func main() {
	m := New()
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
