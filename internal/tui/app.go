package tui

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Page identifies which TUI page is active.
type Page int

const (
	PageHome Page = iota
	PageLearn
	PageExplore
	PagePractice
	PageDashboard
)

// SubModel is implemented by each page model.
type SubModel interface {
	Init() tea.Cmd
	Update(tea.Msg) (SubModel, tea.Cmd)
	View() string
	SetSize(width, height int)
}

// SubModelFactory creates a SubModel. Used to lazily initialize pages.
type SubModelFactory func(width, height int) SubModel

// RootModel is the top-level TUI model that routes between pages.
type RootModel struct {
	currentPage  Page
	pages        map[Page]SubModel
	factories    map[Page]SubModelFactory
	windowWidth  int
	windowHeight int
}

func NewRootModel(factories map[Page]SubModelFactory) RootModel {
	return NewRootModelWithStart(factories, PageHome)
}

func NewRootModelWithStart(factories map[Page]SubModelFactory, startPage Page) RootModel {
	return RootModel{
		currentPage: startPage,
		pages:       make(map[Page]SubModel),
		factories:   factories,
	}
}

func (m RootModel) Init() tea.Cmd {
	return nil
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		for page, sub := range m.pages {
			sub.SetSize(msg.Width, msg.Height-2) // reserve for status bar
			m.pages[page] = sub
		}
		return m, nil

	case NavigateMsg:
		m.currentPage = msg.Page
		if _, ok := m.pages[msg.Page]; !ok {
			if factory, ok := m.factories[msg.Page]; ok {
				sub := factory(m.windowWidth, m.windowHeight-2)
				m.pages[msg.Page] = sub
				return m, sub.Init()
			}
		}
		return m, nil

	case BackMsg:
		if m.currentPage != PageHome {
			m.currentPage = PageHome
		}
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		// On home page, 'q' quits
		if m.currentPage == PageHome && msg.String() == "q" {
			return m, tea.Quit
		}
	}

	// Ensure current page is initialized
	if _, ok := m.pages[m.currentPage]; !ok {
		if factory, ok := m.factories[m.currentPage]; ok {
			sub := factory(m.windowWidth, m.windowHeight-2)
			m.pages[m.currentPage] = sub
			cmd := sub.Init()
			return m, cmd
		}
	}

	// Delegate to current page
	if sub, ok := m.pages[m.currentPage]; ok {
		newSub, cmd := sub.Update(msg)
		m.pages[m.currentPage] = newSub
		return m, cmd
	}

	return m, nil
}

func (m RootModel) View() string {
	if m.windowWidth == 0 {
		return "Loading..."
	}

	// Render current page
	var content string
	if sub, ok := m.pages[m.currentPage]; ok {
		content = sub.View()
	} else {
		content = "Initializing..."
	}

	return lipgloss.NewStyle().
		MaxWidth(m.windowWidth).
		MaxHeight(m.windowHeight).
		Render(content)
}
