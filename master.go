package hopper

import (
    "fmt"
    "log"
    "net"
    "plugin"
    "strconv"
    "os"
    "net/rpc"
    "net/http"
    "hash/maphash"
)

type Hopper struct {
    // Generation N factor
    GenN     int
    // seed map, used as set for deduping seeds and keeping track of Crashes
    seeds    map[uint64]interface{}
    // cov map, used as set for deduping same coverage seeds
    covHash  map[uint64]interface{}
    // Coverage per number of nodes
    cov      []*Coverage
    // Port to host RPC
    port     string
    // PQ of Fuzz tasks (seeds)
    pq       []*FTask
    // Seed for hashing coverage
    hashSeed uint64
}


func (h *Hopper) GetFuzzTask(){
    
}

func (h *Hopper) rpcServer(){
    rpc.Register(h)
    rpc.HandleHTTP()
    l, e := net.Listen("tcp", ":"+strconv.Itoa(h.port))
    if e != nil {                              
        log.Fatal("listen error:", e)
    }                                   
    go http.Serve(l, nil)                         
}

func main() {
    fmt.Fprintf(os.Stderr, "Hopper: hopper -M xxx.so\n")
    GenN := flag.Int("N", 1000, "Number of mutations per generation")
    port := flag.Int("P", 6969, "Port to use, defaults to :6969")
    //thread_mode := flag.Bool("T", false, "Port to use, defaults to :6969")
    
    h := Hopper{
        GenN:       GenN,
        seeds:      make(map[string]interface{}),
        crashes:    make([]Crash, 0),
        port:       port,
        hashSeed:   maphash.NewSeed(),
    }

    h.rpcServer()
}

