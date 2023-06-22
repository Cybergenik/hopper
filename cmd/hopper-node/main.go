package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	n "github.com/Cybergenik/hopper/node"
)

func main() {
	id := flag.Uint64("I", 0, "Node ID, usually just a unique int")
	target := flag.String("T", "", "instrumented target binary")
	args := flag.String("args", "", "args to use against target, ex: --depth=1 @@")
	raw := flag.Bool("raw", false, "should input be fed as pure string (default: input as a file arg)")
	env := flag.String("env", "", "env variables for target seperated by a `;`, ex: ARG1=foo;ARG2=bar;")
	stdin := flag.Bool("stdin", false, "seed should be fed as stdin or as an argument")
	master := flag.String("M", "localhost", "IP/address of Master")
	port := flag.Int("P", 6969, "Port of Master")

	flag.Parse()
	err := ""
	if *id == 0 {
		err += "Hopper Node: Provide a unique Node Id greater than 0: -I \n"
	}
	if *target == "" {
		err += "Hopper Node: Provide a command to run the target: -C\n"
	}
	if !strings.Contains(*args, "@@") && !*stdin {
		err += "Hopper Node: Must provide an @@ input in args if not using stdin: ex --stdin or --args @@\n"
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
	n.Node(*id, *target, *args, *raw, *env, *stdin, masterNode)
}
