package common

import (
    "crypto/md5"
    "encoding/binary"
)

type FTaskID uint64

type BFHash [4]uint64

func Hash(b []byte) FTaskID {
    sum := md5.Sum(b)
    return FTaskID(binary.BigEndian.Uint64([]byte(sum[:])))
}

// baseHashes returns the four hash values of data that are used to create k
// hashes
func BloomHash(data []byte) BFHash {
	var d Digest128 // murmur hashing
	hash1, hash2, hash3, hash4 := d.Sum256(data)
	return [4]uint64{
		hash1, hash2, hash3, hash4,
	}
}

type SeedInfo struct {
    NodeId      uint64
    Id          FTaskID
    Bytes       []byte
    CovHash     BFHash
    CovEdges    uint64
    Crash       bool
}

type Stats struct {
    Its           uint64
    Port          int
    Havoc         uint64
    CrashN        uint64
    SeedsN        uint64
    MaxCov        uint64
    UniqueCrashes int
    UniquePaths   uint64
    Nodes         int
}

type FTask struct {
    Id       FTaskID
    Seed     []byte
    Die      bool
}

type FTaskArgs struct {

}

type UpdateReply struct {
    Log      bool
}

type UpdateFTask struct {
    NodeId   uint64
    Ok       bool
    Id       FTaskID
    CovHash  BFHash
    CovEdges uint64
    Crash    string
}
