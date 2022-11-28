package tui

import (
	"log"
	"strings"
    "strconv"
	tea "github.com/charmbracelet/bubbletea"
    c "github.com/Cybergenik/hopper/common"
    "github.com/charmbracelet/lipgloss"
)

const (
    title = `
    __  __                           
   / / / /___  ____  ____  ___  _____
  / /_/ / __ \/ __ \/ __ \/ _ \/ ___/
 / __  / /_/ / /_/ / /_/ /  __/ /    
/_/ /_/\____/ .___/ .___/\___/_/     
           /_/   /_/                 
`
    hline = `-------------------------------------\n`
)

type StatsMsg struct {
    Stats   c.Stats
}

type Model struct {
    oldStats    c.Stats
    stats       c.Stats
}

//Style
const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	labelstyle = lipgloss.NewStyle().Foreground(hotPink)
	datastyle  = lipgloss.NewStyle().Foreground(darkGray)
)


type Stats struct {
    Its     int
    Port    int
    Havoc   int
    CrashN  int
    SeedsN  int
    MaxSeed Seed
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case StatsMsg:
        m.oldstats = m.stats
        m.stats = msg.Stats
	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	body := fmt.Sprintf(
        `
        %s
        %s
        %s

        %s      %s      %s
        %s      %s      %s

        %s      %s      %s
        %s      %s      %s
        `,
        Hopper,
        hline,
        datastyle.Width(30).Render("Master running on port: "+ stats.Port),
        // Fields
        labelstyle.Width(6).Render("Havoc:"),
        labelstyle.Width(6).Render("Its speed:"),
        labelstyle.Width(6).Render("Max Coverage:"),
        // Data
        datastyle.Width(6).Render(strconv.Itoa(stats.Havoc)),
        datastyle.Width(6).Render(strconv.Itoa(m.stats.Its - m.oldstats.Its)+"/s"),
        datastyle.Width(6).Render(strconv.Itoa(m.stats.MaxSeed.CovEdges)),
        // Fields
        labelstyle.Width(6).Render("# of Seeds:"),
        labelstyle.Width(6).Render("# of Crashes:"),
        labelstyle.Width(6).Render("Fuzz Instances:"),
        // Data
        datastyle.Width(6).Render(m.stats.SeedsN),
        datastyle.Width(6).Render(m.stats.CrashN),
        datastyle.Width(6).Render(m.stats.Its),
    )

	return body
}

func InitialModel() Model{
    return Model{
        oldstats: c.Stats{},
        stats: c.Stats{},
    }
}

