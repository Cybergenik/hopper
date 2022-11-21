package main

import (
    "fmt"
    "log"
    "crypto/md5"
    "net/rpc"
)

type HopperNode struct {
    master      string
}

func (n *HopperNode) hash(cov string) uint64{
    return md5.Sum([]byte(cov))
}

func main() {
    fmt.Fprintf(os.Stderr, "Hopper Node: node -M xxx.so\n")
    id     := flag.String("I", "Node ID, usually just a unique int")
    mut_so := flag.String("m", "mut.so", "compiled object file containing mut engine")
    target := flag.String("t", "instrumented target binary")
    master := flag.String("M", "localhost", "instrumented target binary")
    port   := flag.int("P", 6969, "instrumented target binary")
    
    n := HopperNode{
        master: master+":"+strconv.Itoa(port),
    }
    fmt.Println("Hello there!")
}

func (n *HopperNode) call(rpcname string, args interface{}, reply interface{}) bool {
    c, err := rpc.DialHTTP("tcp", n.master)
    if err != nil {                 
        log.Fatal("dialing:", err)
    }
    defer c.Close()       
                                                            
    err = c.Call(rpcname, args, reply)
    if err != nil {                    
        fmt.Println(err)
        return false
    }                                                          

    return true           
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
    
