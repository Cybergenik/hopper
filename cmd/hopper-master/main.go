package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"plugin"
	"time"

	m "github.com/Cybergenik/hopper/master"
	tui "github.com/Cybergenik/hopper/tui"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: Update Help
func printHelp() {
	fmt.Printf("Hopper Master: go run master.go -I input/ -H=2 -P=6969\n")
}

func initTUI(master *m.Hopper) {
	tui_model := tui.InitModel(master)
	p := tea.NewProgram(tui_model)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}

func readCorpus(path string) [][]byte {
	corpus := [][]byte{}
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("Cannot open input corpus dir: %v", err)
	}
	for _, file := range files {
		seed, err := os.ReadFile(path + "/" + file.Name())
		if err != nil {
			log.Fatalf("Cannot open input in corpus: %v %v", file.Name(), err)
		}
		corpus = append(corpus, seed)
	}
	return corpus
}

func main() {
	help := flag.Bool("help", false, "help menu")
	input := flag.String("I", "", "path to input corpus, directory containing files each being a seed")
	havoc := flag.Uint64("H", 1, "Havoc level to use in mutator, defaults to 1")
	port := flag.Int("P", 6969, "Port to use, defaults to :6969")
	noTui := flag.Bool("no-tui", false, "Don't Generate TUI")
	//TODO: impl thread mode, shouldn't be too hard
	//thread_mode := flag.Bool("T", false, "Port to use, defaults to :6969")
	flag.Parse()
	if *help {
		printHelp()
		os.Exit(0)
	}
	if *input == "" {
		fmt.Fprintln(os.Stderr, "Hopper Master: Please provide a directory with files as input seeds")
	}
	//Parse corpus seeds
	corpus := readCorpus(*input)

	if *noTui {
		h := m.InitHopper(*havoc, *port, m.Mutator, corpus)
		for {
			s := h.Stats()
			fmt.Printf("Its: %d\n", s.Its)
			time.Sleep(10 * time.Second)
		}
	}
	//Init TUI loop
	initTUI(m.InitHopper(*havoc, *port, m.Mutator, corpus))
}

func loadMutEngine(filename string) func([]byte, uint64) []byte {
	p, err := plugin.Open(filename)
	if err != nil {
		log.Fatalf("cannot load plugin %v", filename)
	}

	xmutf, err := p.Lookup("Mutator")
	if err != nil {
		log.Fatalf("cannot find Mutator in %v", filename)
	}

	return xmutf.(func([]byte, uint64) []byte)
}
