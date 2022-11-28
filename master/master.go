package master

import (
    "log"
    "net"
    "sync"
    "strconv"
    "net/rpc"
    "net/http"
    "sync/atomic"
    "hash/maphash"
    c "github.com/Cybergenik/hopper/common"
)

type Seed struct {
    nodeId      int
    bytes       []byte
    covHash     uint64
    covEdges    int
    crash       bool
}

type Hopper struct {
    // Havoc level to use in mutator
    havoc    int
    // seeds and Cov mutex
    mu       sync.Mutex
    // Mutation function
    mutf     func ([]byte, int) []byte
    // seed map, used as set for deduping seeds and keeping track of Crashes
    seeds    map[uint64]Seed
    // cov map, used as set for deduping same coverage seeds
    covHash  map[uint64]interface{}
    // Coverage per number of nodes
    crashes  map[string][]seedInfo
    // Max Coverage in terms of edges
    maxCov   Seed
    // Port to host RPC
    port     int
    // Seed for hashing coverage
    hashSeed maphash.Seed
    // Queue Channel to add new seeds based on energy
    qChan    chan c.FTask
    // Keeps track of whether Hopper has been killed
    dead     int32
    //Stats
    its      int
    crashN   int
    seedsN   int
}

func (h *Hopper) Kill() {
	atomic.StoreInt32(&h.dead, 1)
    h.mu.Lock()
    fmt.Println("%+v", h.crashes)
    h.mu.Unlock()
}

func (h *Hopper) Stats() c.Stats{
    h.mu.Lock()
    defer h.mu.Unlock()
    return c.Stats{
        Its:     h.its,
        Port:    h.port,
        Havoc:   h.havoc,
        CrashN:  h.crashN,
        SeedsN:  h.seedsN,
        MaxSeed: h.maxCov,
    }
}

func (h *Hopper) killed() bool {
	z := atomic.LoadInt32(&h.dead)
	return z == 1
}

func (h *Hopper) energyMutate(seed Seed){
    mutN := 10
    covDiff := seed.covEdges - h.maxCov.covEdges
    if covDiff >= 0 {
        mutN *= covDiff+1
    } else {
        mutN = mutN*(seed.covEdges/h.maxCov.covEdges)
    }
    if seed.crash {
        mutN += 10
    }
    for i:=0;i<mutN;i++{
        mutSeed := h.mutf(seed.bytes, h.havoc)
        h.addSeed(mutSeed)
    }
}

func (h *Hopper) GetFTask(args *interface{}, task *c.FTask) error {
    t := <-h.qChan 
    t.Die = h.killed()
    *task = t
    return nil
}

func (h *Hopper) UpdateFTask(update *c.UpdateFTask, reply *interface{}) error {
    h.mu.Lock()
    defer h.mu.Unlock()
    if val, ok:= h.seeds[update.Id]; ok && h.seeds[update.Id].nodeId == -1 {
        h.its++
        val.nodeId = update.NodeId
        val.covHash = update.CovHash
        val.covEdges = update.CovEdges
        h.seeds[update.Id] = val
        if val.CovEdges > h.maxCov.covEdges{
            h.maxCov = val
        }
        // Dedup based on similar Coverage hash
        if _, ok := h.covHash[update.CovHash]; !ok{
            h.covHash[update.CovHash] = nil
            if (update.Crash != "") {
                h.crashN++
                val := h.seeds[update.Id]
                val.crash = true
                h.seeds[update.Id] = val
                h.crashes[update.Crash] = append(h.crashes[update.Crash], h.seeds[update.Id])
            }
        }
        h.energyMutate(h.seeds[update.Id])
    }
    return nil
}

func (h *Hopper) addSeed(seed []byte){
    seedHash := c.Hash(seed, h.hashSeed)
    if _, ok := h.seeds[seedHash]; !ok {
        return
    }
    h.seedsN++
    h.seeds[seedHash] = seedInfo{
        nodeId:   -1,
        bytes:    seed,
        covHash:  0,
        covEdges: -1,
    }
    h.qChan<-c.FTask{
        Id:       seedHash,
        Seed:     seed,
        HashSeed: h.hashSeed,
    }
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

func InitHopper(havocN int, port int, mutf func([]byte, int) []byte, corpus [][]byte) Hopper{
    h := Hopper{
        havoc:    havocN,
        mutf:     mutf,
        seeds:    make(map[uint64]seedInfo),
        covHash:  make(map[uint64]interface{}),
        crashes:  make(map[string][]seedInfo),
        maxCov:   Seed{},
        port:     port,
        hashSeed: maphash.MakeSeed(),
        //TODO: consider using circular buffer: container/ring
        qChan:    make(chan c.FTask, 1000),
        dead:     0,
    }

	for _, seed := range corpus {
        h.addSeed(seed)
	}

    //Spawn RPC server
    h.rpcServer()

    return h
}

