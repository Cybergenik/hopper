package master

import (
	"container/heap"
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	c "github.com/Cybergenik/hopper/common"
	//Profiling
	//_ "net/http/pprof"
)

type Hopper struct {
	// Campaign ctx
	ctx context.Context
	// Fuzz Campaign Corpus
	corpus [][]byte
	// Havoc level to use in mutator
	havoc uint64
	// seeds and Cov mutex
	mu sync.Mutex
	// PQ mutex
	pqMu sync.Mutex
	// Mutation function
	mutf func([]byte, uint64) []byte
	// PQ of seeds
	pq PriorityQueue
	// seed map, used as temporary rotating buffer while seeds are being fuzzed.
	// Seeds exist ephemerally
	seeds map[c.FTaskID][]byte
	// Seed BloomFilter, used as a set for deduping seeds
	seedBF *BloomFilter
	// Coverage BloomFilter, used as set for deduping same coverage seeds
	coverageBF *BloomFilter
	// Coverage per number of nodes
	crashes map[string][]uint64
	// Max Coverage in terms of edges
	maxCov uint64
	// Port to host RPC
	port int
	// Queue Channel to add new seeds based on energy
	qChan chan c.FTaskID
	// Node IDs
	nodes map[uint64]bool
	//Stats
	its    uint64
	crashN uint64
	seedsN uint64
	paths  uint64
}

const (
	EXP = `
   _____ __        __      
  / ___// /_____  / /______
  \__ \/ __/ __ \/ __/ ___/
 ___/ / /_/ /_/ / /_(__  ) 
/____/\__/\__,_/\__/____/  

Havoc:          %v
Seeds:          %v
Fuzz Instances: %v
Max Edges:      %v
Crashes:        %v
UniqueCrashes:  %v
UniquePaths:    %v
Nodes:          %v

%s
`
)

func (h *Hopper) Report() string {
	h.mu.Lock()
	defer h.mu.Unlock()
	crashes := "Crashes:\n"
	for cType, nodes := range h.crashes {
		crashes += cType + ": "
		for _, node := range nodes {
			crashes += fmt.Sprintf("Node%d ", node)
		}
		crashes += "\n"
	}
	return fmt.Sprintf(
		EXP,
		h.havoc,
		h.seedsN,
		h.its,
		h.maxCov,
		h.crashN,
		len(h.crashes),
		h.paths,
		len(h.nodes),
		crashes,
	)
}

func (h *Hopper) Stats() c.Stats {
	h.mu.Lock()
	defer h.mu.Unlock()
	//fmt.Printf("FTask Q size: %v, Energy PQ size: %v ", len(h.seeds), h.pq.Len())
	return c.Stats{
		Its:             h.its,
		Port:            h.port,
		Havoc:           h.havoc,
		CrashN:          h.crashN,
		SeedsN:          h.seedsN,
		ReadyToFuzzQ:    len(h.qChan),
		EnergyMutationQ: h.pq.Len(),
		MaxCov:          h.maxCov,
		UniqueCrashes:   len(h.crashes),
		UniquePaths:     h.paths,
		Nodes:           len(h.nodes),
	}
}

func (h *Hopper) killed() bool {
	select {
	case <-h.ctx.Done():
		return true
	default:
		return false
	}
}

func (h *Hopper) GetFTask(args *c.FTaskArgs, task *c.FTask) error {
	seedHash := <-h.qChan
	task.Id = seedHash
	h.mu.Lock()
	task.Seed = h.seeds[seedHash]
	h.mu.Unlock()
	task.Die = h.killed()
	return nil
}

func (h *Hopper) UpdateFTask(update *c.UpdateFTask, reply *c.UpdateReply) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.nodes[update.NodeId] = true
	// Check if seed in rotating Task buffer, has it already been processed
	if _, ok := h.seeds[update.Id]; !ok {
		return nil
	}
	h.its++
	// Dump Failed seeds
	if !update.Ok {
		delete(h.seeds, update.Id)
		return nil
	}
	// Track Crashes
	if update.Crash {
		h.crashN++
	}
	// Dedup based on similar Coverage hash
	if !h.coverageBF.ContainsHash(update.CovHash) || h.pq.Len() < 500 {
		h.paths++
		h.coverageBF.AddHash(update.CovHash)
		// Energy Mutations
		s := c.SeedInfo{
			NodeId:   update.NodeId,
			Id:       update.Id,
			CovHash:  update.CovHash,
			CovEdges: update.CovEdges,
			Bytes:    h.seeds[update.Id],
			Crash:    update.Crash,
		}
		go h.energyMutate(s, h.maxCov)

		// Update Max Edge coverage post mutation
		if update.CovEdges > h.maxCov {
			h.maxCov = s.CovEdges
		}
		// Found Unique crash, tell node to Log
		if update.Crash {
			reply.Log = true
			h.crashes[update.CrashMsg] = append(h.crashes[update.CrashMsg], update.NodeId)
		}
	}
	//Free mutated seed
	delete(h.seeds, update.Id)
	return nil
}

func (h *Hopper) mutGenerator() {
	for !h.killed() {
		availableCap := cap(h.qChan) - len(h.qChan)
		if h.pq.Len() > 0 && availableCap >= (cap(h.qChan)/2) {
			//Baseline .01% of available queue capacity
			baseline := float64(availableCap) * .001

			h.pqMu.Lock()
			energyItem := heap.Pop(&h.pq).(*PQItem)
			h.pqMu.Unlock()
			mutN := int(math.Max(1, energyItem.Energy*baseline))
			//fmt.Printf("baseline: %.2f * energy: %.2f = %d", baseline, energyItem.Energy, mutN)
			for i := 0; i < mutN; i++ {
				for ok := h.addSeed(h.mutf(energyItem.Seed, h.havoc)); !ok; {
					ok = h.addSeed(h.mutf(energyItem.Seed, h.havoc))
				}
			}
			// Avoid mem leak
			energyItem = nil
		} else if len(h.qChan) == 0 {
			// Add gen 1 mutated corpus if exploration dies
			for _, seed := range h.corpus {
				for ok := h.addSeed(h.mutf(seed, h.havoc)); !ok; {
					ok = h.addSeed(h.mutf(seed, h.havoc))
				}
			}
		}
	}
}

func (h *Hopper) energyMutate(seed c.SeedInfo, maxEdges uint64) {
	// Energy Range: (0, 1]
	energy := math.Min(1, float64(seed.CovEdges)/float64(maxEdges))
	if seed.Crash {
		energy *= 2
	}
	h.pqMu.Lock()
	heap.Push(
		&h.pq,
		&PQItem{
			Id:       seed.Id,
			Seed:     seed.Bytes,
			Energy:   energy,
			priority: energy,
		},
	)
	h.pqMu.Unlock()
}

// addSeed is by design blocking, we want to block the production of new seeds
// until there is enough space in the Queue
func (h *Hopper) addSeed(seed []byte) bool {
	if h.seedBF.Contains(seed) {
		return false
	}
	h.seedBF.Add(seed)
	seedHash := c.Hash(seed)
	h.mu.Lock()
	h.seedsN++
	h.seeds[seedHash] = seed
	h.mu.Unlock()
	h.qChan <- seedHash
	return true
}

func (h *Hopper) rpcServer() {
	rpc.Register(h)
	rpc.HandleHTTP()
	config := &net.ListenConfig{
		KeepAlive: 0,
	}
	l, e := config.Listen(nil, "tcp", fmt.Sprintf(":%d", h.port))
	//l, e := net.Listen("tcp", ":"+strconv.Itoa(h.port))
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

func (h *Hopper) logger() {
	logInt, ok := os.LookupEnv("HOPPER_LOG_INTERVAL")
	interval := 30
	if ok {
		err := error(nil)
		interval, err = strconv.Atoi(logInt)
		if err != nil {
			log.Fatalf("Invalid HOPPER_LOG_INTERVAL: %v", interval)
		}
	}
	n := 0
	for !h.killed() {
		time.Sleep(time.Minute * time.Duration(interval))
		out_dir, ok := os.LookupEnv("HOPPER_OUT")
		var out string
		if ok {
			out = path.Join(out_dir, fmt.Sprintf("hopper.report.%d", n))
		} else {
			out = fmt.Sprintf("hopper.report.%d", n)
		}

		os.WriteFile(out, []byte(h.Report()), 0666)
		n++
	}
}

func InitHopper(ctx context.Context, havocN uint64, port int, mutf func([]byte, uint64) []byte, corpus [][]byte) *Hopper {
	h := Hopper{
		ctx:        ctx,
		corpus:     corpus,
		havoc:      havocN,
		mutf:       mutf,
		pq:         make(PriorityQueue, 0),
		seeds:      make(map[c.FTaskID][]byte),
		seedBF:     NewWithEstimates(100_000_000, .001),
		coverageBF: NewWithEstimates(100_000_000, .001),
		crashes:    make(map[string][]uint64),
		maxCov:     0,
		port:       port,
		nodes:      make(map[uint64]bool),
		//TODO: consider using circular buffer: container/ring
		qChan:  make(chan c.FTaskID, 10_000),
		its:    0,
		crashN: 0,
		seedsN: 0,
	}

	// Add initial Corpus
	for _, seed := range corpus {
		h.addSeed(seed)
	}

	// Init PQ of energy mutation seeds
	heap.Init(&h.pq)

	// Spawn Energy Mutation Generator
	go h.mutGenerator()

	// Logger
	_, ok := os.LookupEnv("HOPPER_LOG")
	if ok {
		go h.logger()
	}

	// Spawn RPC server
	h.rpcServer()

	//Profiler
	//go func() {
	//    log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()

	return &h
}
