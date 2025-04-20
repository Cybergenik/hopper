package main

import (
	"context"
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

func printHelp() {
	fmt.Printf(
		`NAME:
    Hopper: Master Node

SYNOPSIS: 
    hopper-master [OPTIONS]... -I <dir to corpus>

DESCRIPTION:
    Master node, used to orchestrate fuzzing campaigns

    -I (required)
        path to input corpus, directory containing files each being a seed
    -H 
        havoc level to use in mutator, defaults to 1 (recommended: increase havoc for larger seeds)
    -P 
        port to host Master on, defaults to 6969
    --no-tui 
        Don't Generate TUI, defaults to false
    --help 
        Prints this message

EXAMPLES:
    hopper-master -H=2 -P=6666 -I input/  
        runs master with havoc level 2, on port 6666, using corpus in "input/" directory, where each file is a seed

    hopper-master -H=5 -I test/in
        runs master with havoc level 5, on port 6969, using corpus in "test/in" directory, where each file is a seed
`)
}

func initTUI(ctx context.Context, master *m.Hopper) {
	tuiCtx, cancel := context.WithCancel(ctx)
	tui_model := tui.InitModel(tuiCtx, cancel, master)
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
	flag.Parse()
	if *help || *input == "" {
		printHelp()
		os.Exit(0)
	}
	//Parse corpus seeds
	corpus := readCorpus(*input)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hopper := m.InitHopper(ctx, *havoc, *port, m.Mutator, corpus)
	if *noTui {
		for {
			s := hopper.Stats()
			fmt.Printf("Stats: %+v\n", s)
			time.Sleep(1 * time.Second)
		}
	}
	//Init TUI loop
	initTUI(ctx, hopper)
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
