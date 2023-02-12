package common

import (
    "crypto/md5"
    "encoding/binary"
)

type HashID uint64

func Hash(b []byte) HashID {
    sum := md5.Sum(b)
    return HashID(binary.BigEndian.Uint64([]byte(sum[:])))
}

type Seed struct {
    NodeId      int
    Bytes       *[]byte
    CovHash     HashID
    CovEdges    int
    Crash       bool
}

type Stats struct {
    Its           int
    Port          int
    Havoc         int
    CrashN        int
    SeedsN        int
    MaxSeed       Seed
    UniqueCrashes int
    UniquePaths   int
    Nodes         int
}

type FTask struct {
    Id       HashID
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
    Log      bool
}

type UpdateFTask struct {
    NodeId   int
    Ok       bool
    Id       HashID
    CovHash  HashID
    CovEdges int
    Crash    string
}
