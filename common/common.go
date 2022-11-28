package common

import (
    "crypto/md5"
    "encoding/binary"
)

func Hash(b []byte) uint64{
    sum := md5.Sum(b)
    return binary.BigEndian.Uint64([]byte(sum[:]))
}

type Seed struct {
    NodeId      int
    Bytes       []byte
    CovHash     uint64
    CovEdges    int
    Crash       bool
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
    Die      bool
}

type FTaskArgs struct {

}

type Coverage struct {
    NodeId   int
    Type     string
}

type UpdateReply struct {

}

type UpdateFTask struct {
    NodeId   int
    Ok       bool
    Id       uint64
    CovHash  uint64
    CovEdges int
    Crash    string
}
