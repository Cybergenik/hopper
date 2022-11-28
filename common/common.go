package common

import "hash/maphash"

func Hash(b []byte, seed maphash.Seed) uint64{
    return maphash.Bytes(seed, b)
}

type Stats struct {
    Its     int
    Port    int
    Havoc   int
    CrashN  int
    SeedsN  int
    MaxSeed Seed
}

type FTask struct {
    Id       uint64
	Seed     []byte
    HashSeed maphash.Seed
    Die      bool
}

type Coverage struct {
    NodeId   int
    Type     string
}

type UpdateFTask struct {
    NodeId   int
    Id       uint64
    CovHash  uint64
    CovEdges int
    Crash    string
}
