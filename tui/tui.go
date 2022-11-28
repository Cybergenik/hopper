package tui

import (
	"fmt"
	"strconv"
	"time"

	h "github.com/Cybergenik/hopper/master"
	c "github.com/Cybergenik/hopper/common"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	HOPPER = `
            __  __                           
           / / / /___  ____  ____  ___  _____
          / /_/ / __ \/ __ \/ __ \/ _ \/ ___/
         / __  / /_/ / /_/ / /_/ /  __/ /    
        /_/ /_/\____/ .___/ .___/\___/_/     
                   /_/   /_/                 
`
	hline = `---------------------------------------`
)

type TickMsg time.Time

type Model struct {
	oldStats c.Stats
	stats    c.Stats
    master   *h.Hopper
}

// Style
const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	labelstyle = lipgloss.NewStyle().Foreground(hotPink)
	datastyle  = lipgloss.NewStyle().Foreground(darkGray)
)

func tickStats() tea.Cmd {
	return tea.Every(
		time.Second,
		func(t time.Time) tea.Msg {
			return TickMsg(t)
		},
	)
}

func (m Model) Init() tea.Cmd {
    cmds := []tea.Cmd{tea.ClearScreen, tickStats()}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println("Killing Hopper")
			m.master.Kill()
			return m, tea.Quit
		}
	case TickMsg:
        m.oldStats = m.stats
        m.stats = m.master.Stats()
		return m, tickStats()
	}
	return m, nil
}

func (m Model) View() string {
	body := fmt.Sprintf(
		`
            %s
    %s
        %s

    %s      %s    %s
    %s      %s    %s

    %s      %s    %s
    %s      %s    %s

        %s
        `,
		labelstyle.Width(50).Render(HOPPER),
		datastyle.Width(50).Render(hline),
		datastyle.Width(30).Render("Master running on port: "+strconv.Itoa(m.stats.Port)),
		// Fields
		labelstyle.Width(6).Render("Havoc:"),
		labelstyle.Width(10).Render("Its speed:"),
		labelstyle.Width(6).Render("Edges:"),
		// Data
		datastyle.Width(6).Render(strconv.Itoa(m.stats.Havoc)),
		datastyle.Width(10).Render(strconv.Itoa(m.stats.Its-m.oldStats.Its)+"/s"),
		datastyle.Width(6).Render(strconv.Itoa(m.stats.MaxSeed.CovEdges)),
		// Fields
		labelstyle.Width(6).Render("Seeds:"),
		labelstyle.Width(10).Render("Crashes:"),
		labelstyle.Width(15).Render("Fuzz Instances:"),
		// Data
		datastyle.Width(6).Render(strconv.Itoa(m.stats.SeedsN)),
		datastyle.Width(10).Render(strconv.Itoa(m.stats.CrashN)),
		datastyle.Width(15).Render(strconv.Itoa(m.stats.Its)),
		// Quit
		datastyle.Width(30).Render("Press Esc or Ctrl+C to quit"),
	)

	return body
}

func InitModel(master *h.Hopper) Model {
	return Model{
        oldStats: c.Stats{},
        stats:    c.Stats{},
        master:   master,
	}
}
