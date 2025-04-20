package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	n "github.com/Cybergenik/hopper/node"
)

func printHelp() {
	fmt.Printf(
		`NAME:
    Hopper: Fuzzing Node

SYNOPSIS: 
    hopper-node [OPTIONS]... [--stdin|--args] -I <id> -T <path to target>

DESCRIPTION:
    Fuzzing Node, runs PUT and reports to Hopper Master

    -I 
        node ID, must be a unique unsigned integer
    -T
        path to instrumented target binary
    -M
        IP/address of Master, defaults to localhost
    -P
        port of Master, defaults to 6969
    --raw 
        feed raw seed bytes directly into the PUT, defaults to false. Hopper will put bytes in a file and feed that file to target
    --args 
        args to use against target, ex: --depth=1 @@
    --env 
        env variables for target seperated by a ';' ex: ENV1=foo;ENV2=bar;
    --stdin 
        feed seed through stdin, instead of as argument
    --help 
        prints this message

EXAMPLES:
    hopper-node -I 1 -T target --args "--depth=2 @@"
        runs fuzzing node with id 1, target binary "target", with PUT arg "--depth=2", and it will replace "@@" with a file name of the seed

    hopper-node -I 22 -T /home/user/trash/emacs --stdin
        runs fuzzing node with id 22, target binary "/home/user/trash/emacs", and will feed the seed directly through stdin
`)
}

func main() {
	help := flag.Bool("help", false, "help menu")
	id := flag.Uint64("I", 0, "Node ID, must be a unique unsigned integer")
	target := flag.String("T", "", "path to instrumented target binary")
	args := flag.String("args", "", "args to use against target, ex: --depth=1 @@")
	raw := flag.Bool("raw", false, "should input be fed as pure string (default: input as a file arg)")
	env := flag.String("env", "", "env variables for target seperated by a `;`, ex: ARG1=foo;ARG2=bar;")
	stdin := flag.Bool("stdin", false, "seed should be fed as stdin or as an argument")
	master := flag.String("M", "localhost", "IP/address of Master")
	port := flag.Int("P", 6969, "Port of Master")

	flag.Parse()
	if *help {
		printHelp()
		os.Exit(0)
	}
	var err string
	if *id == 0 {
		err = "Hopper Node: Provide a unique Node Id greater than 0: -I\n"
	} else if *target == "" {
		err = "Hopper Node: Provide a command to run the target: -C\n"
	} else if !strings.Contains(*args, "@@") && !*stdin {
		err = "Hopper Node: Must provide an @@ input in args if not using stdin: ex --stdin or --args @@\n"
	}
	if err != "" {
		err += "--help | -h : to show help menu"
		fmt.Println(err)
		os.Exit(1)
	}
	_, Err := exec.LookPath(n.SANCOV)
	if Err != nil {
		log.Fatalf("Hopper Node: Node requires clang-tools utils: sanvoc")
	}
	masterNode := fmt.Sprintf("%s:%d", *master, *port)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	n.Node(ctx, *id, *target, *args, *raw, *env, *stdin, masterNode)
}
