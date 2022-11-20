package node

import (
    "fmt"
    "log"
)

type HopperNode struct {
    
}

func main() {
    fmt.Fprintf(os.Stderr, "Hopper Node: node -M xxx.so\n")
    mut_so := flag.String("M", "mut.so", "compiled object file containing mut engine")
    target := flag.String("t", "instrumented target binary")
    
    fmt.Println("Hello there!")
}

func loadMutEngine(filename string) func([]byte) []byte {
    p, err := plugin.Open(filename)
    if err != nil {
        log.Fatalf("cannot load plugin %v", filename)
    }

    xmutf, err := p.Lookup("Mutator")
    if err != nil {
        log.Fatalf("cannot find Mutator in %v", filename)
    }
    
    return xmutf.(func([]byte) []byte)
}
    
