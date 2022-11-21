package hopper

type FTask struct {
	Seed     Byte[]
    Id       uint64
    HashSeed uint64
	energy   int
	index    int 
}

type Crash struct {
    NodeId     int
    Type       string
}

type UpdateFTask struct {
    Id          uint64
    CovHash     uint64
    CrashReport Crash
}
