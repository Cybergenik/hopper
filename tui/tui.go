package tui

import (
    "os"
	"fmt"
	"time"
    "path"

	h "github.com/Cybergenik/hopper/master"
	c "github.com/Cybergenik/hopper/common"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	HOPPER = `
      :::    :::  ::::::::  :::::::::  :::::::::  :::::::::: ::::::::: 
     :+:    :+: :+:    :+: :+:    :+: :+:    :+: :+:        :+:    :+: 
    +:+    +:+ +:+    +:+ +:+    +:+ +:+    +:+ +:+        +:+    +:+  
   +#++:++#++ +#+    +:+ +#++:++#+  +#++:++#+  +#++:++#   +#++:++#:    
  +#+    +#+ +#+    +#+ +#+        +#+        +#+        +#+    +#+    
 #+#    #+# #+#    #+# #+#        #+#        #+#        #+#    #+#     
###    ###  ########  ###        ###        ########## ###    ###      
`
	hline = `==================================================`
)

type TickMsg time.Time

type Model struct {
	oldStats c.Stats
	stats    c.Stats
    master   *h.Hopper
}

// Style
const (
	magneta  = lipgloss.Color("#FF00FF")
	darkGray = lipgloss.Color("#C0C0C0")
	snow	 = lipgloss.Color("#FFFAFA")
)

var (
	labelstyle = lipgloss.NewStyle().Foreground(magneta)
	datastyle  = lipgloss.NewStyle().Foreground(snow)
	substyle   =  lipgloss.NewStyle().Foreground(darkGray)
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
		case tea.KeySpace:
            out_dir, ok := os.LookupEnv("HOPPER_OUT")
            var out string
            if ok {
                out = path.Join(out_dir, "hopper.report.Space")
            } else {
                out = "hopper.report.Space"
            }
            os.WriteFile(out, []byte(m.master.Report()), 0666)
			return m, nil
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
    	
    	      %s    %s    %s
    	      %s    %s    %s
    	
    	      %s    %s    %s
    	      %s    %s    %s
    	
    	      %s    %s    %s
    	      %s    %s    %s
    
		    %s
		    %s
        `,
		labelstyle.Width(70).Render(HOPPER),
		datastyle.Width(60).Render(hline),
		datastyle.Width(30).Render(fmt.Sprintf("Master running on port: %d", m.stats.Port)),
		// Fields
		labelstyle.Width(6).Render("Havoc:"),
		labelstyle.Width(15).Render("Its/s:"),
		labelstyle.Width(6).Render("Edges:"),
		// Data
		datastyle.Width(6).Render(fmt.Sprintf("%d", m.stats.Havoc)),
		datastyle.Width(15).Render(fmt.Sprintf("%d/s", m.stats.Its-m.oldStats.Its)),
		datastyle.Width(6).Render(fmt.Sprintf("%d", m.stats.MaxCov)),
		// Fields
		labelstyle.Width(6).Render("Seeds:"),
		labelstyle.Width(15).Render("Crashes:"),
		labelstyle.Width(15).Render("Fuzz Instances:"),
		// Data
		datastyle.Width(6).Render(fmt.Sprintf("%d", m.stats.SeedsN)),
		datastyle.Width(15).Render(fmt.Sprintf("%d", m.stats.CrashN)),
		datastyle.Width(15).Render(fmt.Sprintf("%d", m.stats.Its)),
		// Fields
		labelstyle.Width(6).Render("Nodes:"),
		labelstyle.Width(15).Render("Unique Crashes:"),
		labelstyle.Width(13).Render("Unique Paths:"),
		// Data
		datastyle.Width(6).Render(fmt.Sprintf("%d", m.stats.Nodes)),
		datastyle.Width(15).Render(fmt.Sprintf("%d", m.stats.UniqueCrashes)),
		datastyle.Width(13).Render(fmt.Sprintf("%d", m.stats.UniquePaths)),
		// Quit
		substyle.Width(30).Render("Press Space to Generate Report"),
		substyle.Width(30).Render("Press Esc or Ctrl+C to quit"),
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
