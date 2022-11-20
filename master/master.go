package master

import (
    "fmt"
    "log"
    "net"
    "plugin"
    "os"
    "net/rpc"
    "net/http"
)

type Crash struct {

}

type Hopper struct {
    // Generation N factor
    GenN     int
    // seed map, used as set for deduping seeds
    seeds    map[string]interface{}
    // Coverage per number of nodes
    cov         []*Coverage
    // Port to host RPC
    port     string
}

func (h *Hopper) RPCServer(){
    rpc.Register(h)
    rpc.HandleHTTP()
    l, e := net.Listen("tcp", ":"+h.port)
    if e != nil {                              
        log.Fatal("listen error:", e)
    }                                   
    go http.Serve(l, nil)                         
}

func main() {
    fmt.Fprintf(os.Stderr, "hopper: hopper -M xxx.so\n")
    GenN := flag.Int("N", 1000, "Number of mutations per generation")
    port := flag.String("P", "6969", "Port to use, defaults to :6969")
    //thread_mode := flag.Bool("T", false, "Port to use, defaults to :6969")
    
    h := Hopper{
        GenN:       GenN,
        seeds:      make(map[string]interface{}),
        crashes:    make([]Crash, 0),
        port:       port,
    }

    h.RPCServer()
}

